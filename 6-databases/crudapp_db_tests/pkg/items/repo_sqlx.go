package items

import (
	// "database/sql"
	"github.com/jmoiron/sqlx"
)

type RepoSqlx struct {
	DB *sqlx.DB
}

func NewSqlxRepository(db *sqlx.DB) *RepoSqlx {
	return &RepoSqlx{DB: db}
}

func (repo *RepoSqlx) GetAll() ([]*Item, error) {
	items := make([]*Item, 0, 10)
	err := repo.DB.Select(&items, "SELECT id, title, updated FROM items")
	if err != nil {
		return nil, err
	}
	return items, nil
}

func (repo *RepoSqlx) GetAll_0() ([]*Item, error) {
	items := make([]*Item, 0, 10)
	rows, err := repo.DB.Queryx("SELECT id, title, updated FROM items")
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		item := &Item{}
		// MapScan, SliceScan
		err := rows.StructScan(&item)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}

func (repo *RepoSqlx) GetByID(id int64) (*Item, error) {
	post := &Item{}
	err := repo.DB.Get(post, `SELECT id, title, updated, description FROM items WHERE id = ?`, id)
	if err != nil {
		return nil, err
	}
	return post, nil
}

func (repo *RepoSqlx) Add(elem *Item) (int64, error) {
	result, err := repo.DB.NamedExec(
		`INSERT INTO person (first_name,last_name,email) VALUES (:title, :description)`,
		map[string]interface{}{
			"title":       elem.Title,
			"description": elem.Description,
		})
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (repo *RepoSqlx) Update(elem *Item) (int64, error) {
	result, err := repo.DB.Exec(
		"UPDATE items SET"+
			"`title` = ?"+
			",`description` = ?"+
			",`updated` = ?"+
			"WHERE id = ?",
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

func (repo *RepoSqlx) Delete(id int64) (int64, error) {
	result, err := repo.DB.Exec(
		"DELETE FROM items WHERE id = ?",
		id,
	)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}
