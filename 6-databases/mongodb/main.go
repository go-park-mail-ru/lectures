package main

import (
	"context"
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	// https://www.mongodb.com/blog/post/quick-start-golang--mongodb--modeling-documents-with-go-data-structures
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Item struct {
	Id          primitive.ObjectID `json:"id" bson:"_id"`
	Title       string             `json:"title" bson:"title"`
	Description string             `json:"description" bson:"description"`
	Updated     string             `json:"updated" bson:"updated"`
}

type Handler struct {
	Sess  *mongo.Client
	Items *mongo.Collection
	Tmpl  *template.Template
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {

	items := []*Item{}

	// bson.M{} - это типа условия для поиска
	c, err := h.Items.Find(r.Context(), bson.M{})
	__err_panic(err)
	err = c.All(r.Context(), &items)
	__err_panic(err)

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

	newItem := bson.M{
		"_id":         primitive.NewObjectID(),
		"title":       r.FormValue("title"),
		"description": r.FormValue("description"),
		"some_filed":  123,
	}
	_, err := h.Items.InsertOne(r.Context(), newItem)
	__err_panic(err)

	fmt.Println("Insert - LastInsertId:", newItem["_id"])

	http.Redirect(w, r, "/", http.StatusFound)
}

func (h *Handler) Edit(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	if !primitive.IsValidObjectID(vars["id"]) {
		http.Error(w, "bad id", 500)
		return
	}
	id, err := primitive.ObjectIDFromHex(vars["id"])
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	post := &Item{}
	err = h.Items.FindOne(r.Context(), bson.M{"_id": id}).Decode(post)
	__err_panic(err)

	err = h.Tmpl.ExecuteTemplate(w, "edit.html", post)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	if !primitive.IsValidObjectID(vars["id"]) {
		http.Error(w, "bad id", 500)
		return
	}
	id, err := primitive.ObjectIDFromHex(vars["id"])
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	post := &Item{}
	err = h.Items.FindOne(r.Context(), bson.M{"_id": id}).Decode(&post)
	__err_panic(err)

	post.Title = r.FormValue("title")
	post.Description = r.FormValue("description")
	post.Updated = "rvasily"

	res, err := h.Items.UpdateOne(
		r.Context(),
		bson.M{"_id": id},
		// про другие операторы помимо $set можно почитать тут
		// https://www.mongodb.com/docs/manual/reference/operator/update/
		bson.M{"$set": bson.M{
			"title":       r.FormValue("title"),
			"description": r.FormValue("description"),
			"updated":     "rvasily",
			"newField":    123,
		}},
	)
	affected := 1
	if err != nil {
		__err_panic(err)
	} else if res.ModifiedCount == 0 {
		affected = 0
	}

	fmt.Println("Update - RowsAffected", affected)

	http.Redirect(w, r, "/", http.StatusFound)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	if !primitive.IsValidObjectID(vars["id"]) {
		http.Error(w, "bad id", 500)
		return
	}
	id, err := primitive.ObjectIDFromHex(vars["id"])
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	res, err := h.Items.DeleteOne(r.Context(), bson.M{"_id": id})
	affected := 1
	if err != nil {
		__err_panic(err)
	} else if res.DeletedCount == 0 {
		affected = 0
	}

	w.Header().Set("Content-type", "application/json")
	resp := `{"affected": ` + strconv.Itoa(int(affected)) + `}`
	w.Write([]byte(resp))
}

func main() {
	ctx := context.Background()
	sess, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost"))
	__err_panic(err)

	collection := sess.Database("golang").Collection("items")

	// если коллекции не будет, то она создасться автоматически
	// для монги нет такого красивого дампа SQL, так что я вставляю демо-запись если коллекция пуста
	if n, _ := collection.CountDocuments(ctx, bson.M{}); n == 0 {
		collection.InsertOne(ctx, &Item{
			Id:          primitive.NewObjectID(),
			Title:       "mongodb",
			Description: "Рассказать про монгу",
			Updated:     "",
		})
		collection.InsertOne(ctx, &Item{
			Id:          primitive.NewObjectID(),
			Title:       "redis",
			Description: "Рассказать про redis",
			Updated:     "rvasily",
		})
	}

	handlers := &Handler{
		Items: collection,
		Tmpl:  template.Must(template.ParseGlob("./templates/*")),
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

	fmt.Println("starting server at :8088")
	err = http.ListenAndServe(":8088", r)
	__err_panic(err)
}

// не используйте такой код в прошакшене
// ошибка должна всегда явно обрабатываться
func __err_panic(err error) {
	if err != nil {
		panic(err)
	}
}
