package items

import (
	"database/sql"
)

type RepoPgx struct {
	DB *sql.DB
}

func NewPgxRepository(db *sql.DB) *RepoPgx {
	return &RepoPgx{DB: db}
}

func (repo *RepoPgx) GetAll() ([]*Item, error) {
	items := []*Item{}
	rows, err := repo.DB.Query("SELECT id, title, updated FROM items")
	if err != nil {
		return nil, err
	}
	defer rows.Close() // надо закрывать соединение, иначе будет течь
	for rows.Next() {
		post := &Item{}
		err = rows.Scan(&post.ID, &post.Title, &post.Updated)
		if err != nil {
			return nil, err
		}
		items = append(items, post)
	}
	/*
		func tx(db *sq.DB, fb func(tx *sql.Tx) error) error {
			tx := repo.DB.Begin()
			err := fb(tx)
			if err != nil {
				tx.Rollback()
				return err
			}
			tx.Commit()
			return nil
		}

		tx(repo.DB, func(tx *sql.Tx) error {
			tx.Query("select")

			if err != nil {
				return err
			}

			tx.Exec("update")

			return nil
		})
	*/

	return items, nil
}

func (repo *RepoPgx) GetByID(id int64) (*Item, error) {
	post := &Item{}
	// QueryRow сам закрывает коннект
	err := repo.DB.
		QueryRow(`SELECT id, title, updated, description FROM items WHERE id = $1`, id).
		Scan(&post.ID, &post.Title, &post.Updated, &post.Description)
	if err != nil {
		return nil, err
	}
	return post, nil
}

func (repo *RepoPgx) Add(elem *Item) (int64, error) {
	var lastInsertId int64
	err := repo.DB.QueryRow(
		`INSERT INTO items ("title", "description") VALUES ($1, $2) RETURNING id`,
		elem.Title,
		elem.Description,
	).Scan(&lastInsertId)
	if err != nil {
		return 0, err
	}
	return lastInsertId, nil
}

func (repo *RepoPgx) Update(elem *Item) (int64, error) {
	result, err := repo.DB.Exec(
		`UPDATE items SET "title" = $1`+
			`,"description" = $2`+
			`,"updated" = $3`+
			`WHERE id = $4`,
		elem.Title,
		elem.Description,
		"rvasily",
		elem.ID,
	)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func (repo *RepoPgx) Delete(id int64) (int64, error) {
	result, err := repo.DB.Exec(
		"DELETE FROM items WHERE id = $1",
		id,
	)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}
