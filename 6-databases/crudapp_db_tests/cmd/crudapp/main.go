package main

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/go-park-mail-ru/lectures/6-databases/crudapp_db_tests/pkg/handlers"
	"github.com/go-park-mail-ru/lectures/6-databases/crudapp_db_tests/pkg/items"
	"github.com/go-park-mail-ru/lectures/6-databases/crudapp_db_tests/pkg/middleware"
	"github.com/go-park-mail-ru/lectures/6-databases/crudapp_db_tests/pkg/session"
	"github.com/go-park-mail-ru/lectures/6-databases/crudapp_db_tests/pkg/user"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

func getMysql() *sql.DB {
	dsn := "root:love@tcp(localhost:3306)/golang?&charset=utf8&interpolateParams=true"
	db, err := sql.Open("mysql", dsn)
	err = db.Ping() // вот тут будет первое подключение к базе
	if err != nil {
		log.Fatalln(err)
	}
	db.SetMaxOpenConns(10)
	return db
}

func getPostgres() *sql.DB {
	dsn := "user=postgres dbname=golang password=love host=127.0.0.1 port=5432 sslmode=disable"
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Fatalln("cant parse config", err)
	}
	err = db.Ping() // вот тут будет первое подключение к базе
	if err != nil {
		log.Fatalln(err)
	}
	db.SetMaxOpenConns(10)
	return db
}

// http://jmoiron.github.io/sqlx/
func getSqlx() *sqlx.DB {
	return sqlx.NewDb(getMysql(), "mysql")
}

// https://gorm.io/
func getGorm() *gorm.DB {
	db, err := gorm.Open(mysql.New(mysql.Config{
		Conn: getMysql(),
	}), &gorm.Config{})
	if err != nil {
		log.Fatalln(err)
	}
	return db
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
	// itemsRepo := items.NewMysqlRepository(getMysql())
	// itemsRepo := items.NewSqlxRepository(getSqlx())
	itemsRepo := items.NewGormRepository(getGorm())
	// itemsRepo := items.NewPgxRepository(getPostgres())

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
