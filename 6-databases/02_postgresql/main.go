package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type Item struct {
	Id          int
	Title       string
	Description string
	Updated     sql.NullString
}

type Handler struct {
	DB   *sql.DB
	Tmpl *template.Template
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {

	items := []*Item{}
	// Не надо так: SELECT * FROM items
	rows, err := h.DB.QueryContext(r.Context(), "SELECT id, title, updated FROM items")
	panicOnErr(err)

	// Надо закрывать соединение, иначе будет течь
	defer rows.Close()

	for rows.Next() {
		post := &Item{}
		err = rows.Scan(&post.Id, &post.Title, &post.Updated)
		panicOnErr(err)
		items = append(items, post)
	}

	err = h.Tmpl.ExecuteTemplate(w, "index.html", struct {
		Items []*Item
	}{
		Items: items,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) AddForm(w http.ResponseWriter, r *http.Request) {
	err := h.Tmpl.ExecuteTemplate(w, "create.html", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) Add(w http.ResponseWriter, r *http.Request) {
	// В целях упрощения примера пропущена валидация
	var insertedID int
	err := h.DB.QueryRow(
		"INSERT INTO items (title, description) VALUES ($1, $2) RETURNING id",
		r.FormValue("title"),
		r.FormValue("description"),
	).Scan(&insertedID)
	panicOnErr(err)

	fmt.Println("Insert - InsertedId: ", insertedID)

	http.Redirect(w, r, "/", http.StatusFound)
}

func (h *Handler) Edit(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	panicOnErr(err)

	post := &Item{}
	row := h.DB.QueryRow("SELECT id, title, updated, description FROM items WHERE id = $1", id)

	// Scan сам закрывает коннект
	err = row.Scan(&post.Id, &post.Title, &post.Updated, &post.Description)
	panicOnErr(err)

	err = h.Tmpl.ExecuteTemplate(w, "edit.html", post)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	panicOnErr(err)

	// В целях упрощения примера пропущена валидация
	result, err := h.DB.Exec(
		"UPDATE items SET title = $1, description = $2, updated = $3 WHERE id = $4",
		r.FormValue("title"), r.FormValue("description"), "user", id,
	)
	panicOnErr(err)

	affected, err := result.RowsAffected()
	panicOnErr(err)

	fmt.Println("Update - RowsAffected", affected)

	http.Redirect(w, r, "/", http.StatusFound)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	panicOnErr(err)

	result, err := h.DB.Exec(
		"DELETE FROM items WHERE id = $1",
		id,
	)
	panicOnErr(err)

	affected, err := result.RowsAffected()
	panicOnErr(err)

	fmt.Println("Delete - RowsAffected", affected)

	w.Header().Set("Content-type", "application/json")
	resp := `{"affected": ` + strconv.Itoa(int(affected)) + `}`
	w.Write([]byte(resp))
}

func main() {

	// Основные настройки подключения к базе
	dsn := "postgres://postgres_user:postgres_password@localhost:5432/golang?"
	// Отключаем использование SSL
	dsn += "sslmode=disable"

	db, err := sql.Open("pgx", dsn)
	panicOnErr(err)

	db.SetMaxOpenConns(10)

	err = db.Ping() // Тут будет первое подключение к базе
	if err != nil {
		panic(err)
	}

	handlers := &Handler{
		DB:   db,
		Tmpl: template.Must(template.ParseGlob("templates/*")),
	}

	// В целях упрощения примера пропущена авторизация и csrf
	r := mux.NewRouter()
	r.HandleFunc("/", handlers.List).Methods("GET")
	r.HandleFunc("/items", handlers.List).Methods("GET")
	r.HandleFunc("/items/new", handlers.AddForm).Methods("GET")
	r.HandleFunc("/items/new", handlers.Add).Methods("POST")
	r.HandleFunc("/items/{id}", handlers.Edit).Methods("GET")
	r.HandleFunc("/items/{id}", handlers.Update).Methods("POST")
	r.HandleFunc("/items/{id}", handlers.Delete).Methods("DELETE")

	fmt.Println("starting server at :8080")
	fmt.Println(http.ListenAndServe(":8080", r))
}

// Не используйте такой код в продакшене
// Ошибка должна всегда явно обрабатываться
func panicOnErr(err error) {
	if err != nil {
		panic(err)
	}
}
