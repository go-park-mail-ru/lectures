package main

import (
	"html/template"
	"net/http"

	"crudapp/internal/crudapp/middleware"
	items_delivery "crudapp/internal/pkg/items/delivery"
	items_repo "crudapp/internal/pkg/items/repository"
	"crudapp/internal/pkg/session"
	user_delivery "crudapp/internal/pkg/user/delivery"
	user_repo "crudapp/internal/pkg/user/repository"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

func main() {
	templates := template.Must(template.ParseGlob("./templates/*"))

	sm := session.NewSessionsMem()
	zapLogger, _ := zap.NewProduction()
	defer zapLogger.Sync() // flushes buffer, if any
	logger := zapLogger.Sugar()

	userRepo := user_repo.NewUserRepo()
	itemsRepo := items_repo.NewRepo()

	userHandler := &user_delivery.UserHandler{
		Tmpl:     templates,
		UserRepo: userRepo,
		Logger:   logger,
		Sessions: sm,
	}

	handlers := &items_delivery.ItemsHandler{
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
