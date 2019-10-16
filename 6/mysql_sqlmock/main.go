package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	_ "github.com/go-sql-driver/mysql"
)

type ItemCRUD interface {
	ListAll() ([]*Item, error)
	SelectByID(int64) (*Item, error)
	Create(*Item) (int64, error)
	Update(*Item) (int64, error)
	Delete(int64) (int64, error)
}

type Handler struct {
	Items ItemCRUD
	Tmpl  *template.Template
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	items, err := h.Items.ListAll()
	if err != nil {
		log.Println("Items.ListAll err:", err)
		http.Error(w, "db err", 500)
		return
	}
	err = h.Tmpl.ExecuteTemplate(w, "index.html", struct {
		Items []*Item
	}{
		Items: items,
	})
	if err != nil {
		log.Println("Tmpl.ExecuteTemplate err:", err)
		http.Error(w, "template expand err", 500)
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
	elem := &Item{
		Title:       r.FormValue("title"),
		Description: r.FormValue("description"),
	}
	lastID, err := h.Items.Create(elem)
	if err != nil {
		log.Println("Items.Create err:", err)
		http.Error(w, "db err", 500)
		return
	}

	fmt.Println("LastInsertId: ", lastID)

	http.Redirect(w, r, "/", http.StatusFound)
}

func (h *Handler) Edit(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil || id == 0 {
		log.Println("bad id:", id, err)
		http.Error(w, "bad request", 400)
		return
	}

	post, err := h.Items.SelectByID(int64(id))
	if err != nil {
		log.Println("Items.SelectByID err:", err)
		http.Error(w, "db err", 500)
		return
	}

	err = h.Tmpl.ExecuteTemplate(w, "edit.html", post)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	__err_panic(err)

	elem := &Item{
		ID:          int64(id),
		Title:       r.FormValue("title"),
		Description: r.FormValue("description"),
	}
	affected, err := h.Items.Update(elem)
	if err != nil {
		log.Println("Items.Update err:", err)
		http.Error(w, "db err", 500)
		return
	}

	fmt.Println("Update - RowsAffected", affected)

	http.Redirect(w, r, "/", http.StatusFound)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	__err_panic(err)

	affected, err := h.Items.Delete(int64(id))
	fmt.Println("Delete - RowsAffected", affected, err)

	w.Header().Set("Content-type", "application/json")
	resp := `{"affected": ` + strconv.Itoa(int(affected)) + `}`
	w.Write([]byte(resp))
}

func main() {

	// основные настройки к базе
	dsn := "root:1234@tcp(localhost:3306)/tech?"
	// указываем кодировку
	dsn += "&charset=utf8"
	// отказываемся от prapared statements
	// параметры подставляются сразу
	dsn += "&interpolateParams=true"

	db, err := sql.Open("mysql", dsn)

	db.SetMaxOpenConns(10)

	err = db.Ping() // вот тут будет первое подключение к базе
	if err != nil {
		panic(err)
	}

	itemsRepo := &ItemRepository{
		DB: db,
	}

	handlers := &Handler{
		Items: itemsRepo,
		Tmpl:  template.Must(template.ParseGlob("../crud_templates/*")),
	}

	// в целям упрощения примера пропущена авторизация и csrf
	r := mux.NewRouter()
	r.HandleFunc("/", handlers.List).Methods("GET")
	r.HandleFunc("/items", handlers.List).Methods("GET")
	r.HandleFunc("/items/new", handlers.AddForm).Methods("GET")
	r.HandleFunc("/items/new", handlers.Add).Methods("POST")
	r.HandleFunc("/items/{id}", handlers.Edit).Methods("GET")
	r.HandleFunc("/items/{id}", handlers.Update).Methods("POST")
	r.HandleFunc("/items/{id}", handlers.Delete).Methods("DELETE")

	fmt.Println("starting server at :8080")
	http.ListenAndServe(":8080", r)
}

// не используйте такой код в прошакшене
// ошибка должна всегда явно обрабатываться
func __err_panic(err error) {
	if err != nil {
		panic(err)
	}
}
