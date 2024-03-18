package items

import (
	"database/sql"
)

type RepoMysql struct {
	DB *sql.DB
}

func NewMysqlRepository(db *sql.DB) *RepoMysql {
	return &RepoMysql{DB: db}
}

func (repo *RepoMysql) GetAll(limit int) ([]*Item, error) {
	items := make([]*Item, 0, limit)
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

/*

	dsn += "&interpolateParams=false" (или нет параметра)
	QueryRow("SELECT * FROM items WHERE id = ?", id).
	->
	smtp := db.PrepareStatement("SELECT * FROM items WHERE id = ?")
	row := smtp.Execute(smtp, 1)



	dsn += "&interpolateParams=true"
	smtp := db.QueryRaw("SELECT * FROM items WHERE id = 1")



	params := make([]string, 0, len(manyIds))
	values := make([]interface{}, 0, len(manyIds))
	for _, val := manyIds {
		params = append(params, "?")
		values = append(values, val)
	}

	q := fmt.Sprintf(`where id in(%s)`, string.Join(params, `,`))
	db.Query(q, values...)

*/

func (repo *RepoMysql) GetByID(id int64) (*Item, error) {
	post := &Item{}
	// QueryRow сам закрывает коннект
	err := repo.DB.
		QueryRow(`SELECT id, title, updated, description FROM items WHERE id = ?`, id).
		Scan(&post.ID, &post.Title, &post.Updated, &post.Description)
	// если запись не найден - вернется sql.ErrNoRows
	if err != nil {
		return nil, err
	}
	return post, nil
}

func (repo *RepoMysql) Add(elem *Item) (int64, error) {
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

func (repo *RepoMysql) Update(elem *Item) (int64, error) {
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

func (repo *RepoMysql) Delete(id int64) (int64, error) {
	result, err := repo.DB.Exec(
		"DELETE FROM items WHERE id = ?",
		id,
	)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}
