package main

import (
	"context"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/go-park-mail-ru/lectures/6-databases/crudapp_mongo/pkg/handlers"
	"github.com/go-park-mail-ru/lectures/6-databases/crudapp_mongo/pkg/items"
	"github.com/go-park-mail-ru/lectures/6-databases/crudapp_mongo/pkg/middleware"
	"github.com/go-park-mail-ru/lectures/6-databases/crudapp_mongo/pkg/session"
	"github.com/go-park-mail-ru/lectures/6-databases/crudapp_mongo/pkg/user"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.uber.org/zap"
)

func getMongo(cfg string) *mongo.Client {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg))
	if err != nil {
		log.Fatalln("cant connect to mongo", err)
	}
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatalln("cant ping mongo", err)
	}
	return client
}

func main() {

	// -----

	templates := template.Must(template.ParseGlob("./templates/*"))

	sm := session.NewSessionsMem()
	zapLogger, _ := zap.NewProduction()
	defer zapLogger.Sync() // flushes buffer, if any
	logger := zapLogger.Sugar()

	// dbMSSQL := sql.Open(...)

	userRepo := user.NewUserRepo()

	mongoClient := getMongo("mongodb://localhost:27017")
	itemsRepo := items.NewMongoRepository(mongoClient.Database("coursera").Collection("items"))

	userHandler := &handlers.UserHandler{
		Tmpl:     templates,
		UserRepo: userRepo,
		Logger:   logger,
		Sessions: sm,
	}

	handlers := &handlers.ItemsHandler{
		Tmpl:      templates,
		Logger:    logger,
		ItemsRepo: itemsRepo,
	}

	r := mux.NewRouter()
	r.HandleFunc("/", userHandler.Index).Methods("GET")
	r.HandleFunc("/login", userHandler.Login).Methods("POST")
	r.HandleFunc("/logout", userHandler.Logout).Methods("POST")

	r.HandleFunc("/items", handlers.List).Methods("GET")
	r.HandleFunc("/items/new", handlers.AddForm).Methods("GET")
	r.HandleFunc("/items/new", handlers.Add).Methods("POST")
	r.HandleFunc("/items/{id}", handlers.Edit).Methods("GET")
	r.HandleFunc("/items/{id}", handlers.Update).Methods("POST")
	r.HandleFunc("/items/{id}", handlers.Delete).Methods("DELETE")

	mux := middleware.Auth(sm, r)
	mux = middleware.AccessLog(logger, mux)
	mux = middleware.Panic(mux)

	addr := ":8080"
	logger.Infow("starting server",
		"type", "START",
		"addr", addr,
	)
	http.ListenAndServe(addr, mux)
}
