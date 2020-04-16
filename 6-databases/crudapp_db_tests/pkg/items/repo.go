package items

import (
	"database/sql"
)

type ItemRepository struct {
	DB *sql.DB
}

func NewRepository(db *sql.DB) *ItemRepository {
	return &ItemRepository{DB: db}
}

func (repo *ItemRepository) GetAll() ([]*Item, error) {
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
	return items, nil
}

func (repo *ItemRepository) GetByID(id int64) (*Item, error) {
	post := &Item{}
	// QueryRow сам закрывает коннект
	err := repo.DB.
		QueryRow("SELECT id, title, updated, description FROM items WHERE id = ?", id).
		Scan(&post.ID, &post.Title, &post.Updated, &post.Description)
	if err != nil {
		return nil, err
	}
	return post, nil
}

func (repo *ItemRepository) Add(elem *Item) (int64, error) {
	result, err := repo.DB.Exec(
		"INSERT INTO items (`title`, `description`) VALUES (?, ?)",
		elem.Title,
		elem.Description,
	)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (repo *ItemRepository) Update(elem *Item) (int64, error) {
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

func (repo *ItemRepository) Delete(id int64) (int64, error) {
	result, err := repo.DB.Exec(
		"DELETE FROM items WHERE id = ?",
		id,
	)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}
