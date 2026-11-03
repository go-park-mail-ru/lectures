package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"lectures/cleanarch3/internal/domain"
	"lectures/cleanarch3/internal/usecase"
)

// Модель строки БД остаётся внутри адаптера.
type userRow struct {
	id          int64
	name, email string
}

func (row userRow) toDomain() (domain.User, error) {
	return domain.RestoreUser(row.id, row.name, row.email)
}

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

var _ usecase.UserRepository = (*UserRepository)(nil)

func (r *UserRepository) Create(ctx context.Context, u domain.User) (domain.User, error) {
	if u.ID() != 0 {
		return domain.User{}, usecase.ErrInvalidInput
	}
	if _, err := domain.NewUser(u.Name(), u.Email()); err != nil {
		return domain.User{}, err
	}
	var row userRow
	err := r.db.QueryRowContext(ctx, `
		INSERT INTO users (name, email) VALUES ($1, $2)
		RETURNING id, name, email`, u.Name(), u.Email()).
		Scan(&row.id, &row.name, &row.email)
	if err != nil {
		return domain.User{}, fmt.Errorf("insert user: %w", translateError(err))
	}
	return row.toDomain()
}

func (r *UserRepository) GetByID(ctx context.Context, id int64) (domain.User, error) {
	var row userRow
	err := r.db.QueryRowContext(ctx, `SELECT id, name, email FROM users WHERE id = $1`, id).
		Scan(&row.id, &row.name, &row.email)
	if err != nil {
		return domain.User{}, fmt.Errorf("select user: %w", translateError(err))
	}
	return row.toDomain()
}

func (r *UserRepository) List(ctx context.Context, afterID int64, limit int) ([]domain.User, error) {
	if afterID < 0 || limit <= 0 {
		return nil, usecase.ErrInvalidInput
	}
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, name, email FROM users
		WHERE id > $1 ORDER BY id LIMIT $2`, afterID, limit)
	if err != nil {
		return nil, fmt.Errorf("list users: %w", err)
	}
	defer rows.Close()
	items := make([]domain.User, 0)
	for rows.Next() {
		var row userRow
		if err := rows.Scan(&row.id, &row.name, &row.email); err != nil {
			return nil, fmt.Errorf("scan user: %w", err)
		}
		u, err := row.toDomain()
		if err != nil {
			return nil, err
		}
		items = append(items, u)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("read users: %w", err)
	}
	return items, nil
}

func (r *UserRepository) Update(ctx context.Context, u domain.User) error {
	if _, err := domain.RestoreUser(u.ID(), u.Name(), u.Email()); err != nil {
		return err
	}
	result, err := r.db.ExecContext(ctx, `UPDATE users SET name = $1, email = $2 WHERE id = $3`,
		u.Name(), u.Email(), u.ID())
	if err != nil {
		return fmt.Errorf("update user: %w", translateError(err))
	}
	return requireAffectedRow(result)
}

func (r *UserRepository) Delete(ctx context.Context, id int64) error {
	result, err := r.db.ExecContext(ctx, `DELETE FROM users WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete user: %w", err)
	}
	return requireAffectedRow(result)
}

func requireAffectedRow(result sql.Result) error {
	count, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("count affected users: %w", err)
	}
	if count == 0 {
		return usecase.ErrNotFound
	}
	return nil
}

func translateError(err error) error {
	if errors.Is(err, sql.ErrNoRows) {
		return usecase.ErrNotFound
	}
	// Драйвер может сообщать код ошибки через этот маленький интерфейс.
	// В нашей схеме ID выдаёт БД, а уникальность проверяется для email.
	// При добавлении других UNIQUE нужно различать имена ограничений.
	var state interface{ SQLState() string }
	if errors.As(err, &state) && state.SQLState() == "23505" {
		return usecase.ErrEmailTaken
	}
	return err
}
