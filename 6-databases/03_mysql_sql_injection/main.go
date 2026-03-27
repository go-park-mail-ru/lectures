package main

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

var loginFormTmpl = `
<html>
	<body>
	<form action="/login" method="post">
		Login: <input type="text" name="login">
		Password: <input type="password" name="password">
		<input type="submit" value="Login">
	</form>
	</body>
</html>
`

func main() {

	// Основные настройки подключения к базе
	dsn := "root:love@tcp(localhost:3306)/golang?"
	// Указываем кодировку
	dsn += "&charset=utf8"
	// Отказываемся от prepared statements
	// Параметры подставляются сразу
	dsn += "&interpolateParams=true"

	var err error
	// Создаём структуру базы
	// Но соединение происходит только при первом запросе
	db, err := sql.Open("mysql", dsn)
	PanicOnErr(err)

	err = db.Ping() // Тут будет первое подключение к базе
	PanicOnErr(err)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(loginFormTmpl))
	})

	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		var (
			id          int
			login, body string
		)

		inputLogin := r.FormValue("login")
		body += fmt.Sprintln("inputLogin:", inputLogin)

		// ПЛОХО! НЕ ДЕЛАЙТЕ ТАК!
		// Параметры не экранированы должным образом
		// Мы подставляем в запрос параметр как есть
		query := fmt.Sprintf("SELECT id, login FROM users WHERE login = '%s' LIMIT 1", inputLogin) // SELECT id, login FROM users WHERE login = '' OR '1'<>'2' LIMIT 1
		//query := fmt.Sprintf("SELECT id, login FROM users WHERE (login = '%s' AND login <> 'admin') LIMIT 1", inputLogin) // SELECT id, login FROM users WHERE (login = 'admin') #' AND login <> 'admin') LIMIT 1

		body += fmt.Sprintln("Sprint query:", query)

		row := db.QueryRow(query)
		err := row.Scan(&id, &login)

		if errors.Is(err, sql.ErrNoRows) {
			body += fmt.Sprintln("Sprint case: NOT FOUND")
		} else {
			PanicOnErr(err)
			body += fmt.Sprintln("Sprint case: id:", id, "login:", login)
		}

		// ПРАВИЛЬНО
		// Мы используем плейсхолдеры, там параметры будет экранирован должным образом
		row = db.QueryRow("SELECT id, login FROM users WHERE login = ? LIMIT 1", inputLogin)
		err = row.Scan(&id, &login)
		if errors.Is(err, sql.ErrNoRows) {
			body += fmt.Sprintln("Placeholders case: NOT FOUND")
		} else {
			PanicOnErr(err)
			body += fmt.Sprintln("Placeholders id:", id, "login:", login)
		}

		w.Write([]byte(body))
	})

	fmt.Println("starting server at :8080")
	PanicOnErr(http.ListenAndServe(":8080", nil))
}

// PanicOnErr panics on error
func PanicOnErr(err error) {
	if err != nil {
		panic(err)
	}
}
