// XSS - это внедрение вредоносного кода там где мы не ожидаем
// например в комменте пишем JS, который будет выполняться для всех пользователей, читающих его
// опасность заключается в том, что злоумышленник может вызывать от имени юзера какие-то методы
// например отправка спама от его имени или кража сессии
// лечится правильным экранированием всех внешних входных данных по отношению к сайту (в первую очередь - пользовательского ввода)

// данный пример ИСКУССТВЕННЫЙ, чтобы показать как проявляется XSS
// используйте пакет html/template
// он автоматичски экранирует все входящие данные с учетом контекста
// подрбонее https://golang.org/pkg/html/template/

package main

import (
	"fmt"
	"html/template"
	"math/rand"
	"net/http"
	// "text/template"
	"time"
)

var sessions = map[string]string{}
var messages = []string{"Hello World"}

var loginFormTmplRaw = `<html><body>
	<form action="/login" method="post">
		Login: <input type="text" name="login" value="DefaultUserName">
		Password: <input type="password" name="password" value="anypass">
		<input type="submit" value="Login">
	</form>
</body></html>`

var messagesTmpl = `<html><body>

	&lt;script&gt;alert(document.cookie)&lt;/script&gt;

	<br />
	<br />

	<form action="/comment" method="post">
		<textarea name="comment"></textarea><br />
		<input type="submit" value="Comment">
	</form>

	<br />

    {{range .Messages}}
		<div style="border: 1px solid black; padding: 5px; margin: 5px;">
			<!-- text/template по-умолч ничего не экранируется, надо указать | html --> 
			<!-- html/template по-умолч будет экранировать --> 

			{{.}}
		</div>
    {{end}}
</body></html>`

func checkSession(r *http.Request) bool {
	// обработка сессии
	// не используйте этот подход в продакшене
	sessionID, err := r.Cookie("session_id")
	if err == http.ErrNoCookie {
		return false
	} else if err != nil {
		PanicOnErr(err)
	}
	_, ok := sessions[sessionID.Value]
	if !ok {
		return false
	}
	return true
}

func main() {

	tmpl := template.New("main")
	tmpl, _ = tmpl.Parse(messagesTmpl)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if !checkSession(r) {
			w.Write([]byte(loginFormTmplRaw))
			return
		}

		// для отключения экранирования в html/template
		newMsgs := []template.HTML{}
		for _, v := range messages {
			newMsgs = append(newMsgs, template.HTML(v))
		}

		//выводим комментарии + форму отправки
		tmpl.Execute(w, struct {
			Messages []template.HTML
		}{
			Messages: newMsgs,
		})
	})

	// добавление комментария
	// добавим комментарий c текстом
	/*
		<script>alert(document.cookie)</script>
	*/
	// это выведет на экран куки сайта. дальше с ними можно сделать всё что угодно - например отправить на внешний сервис, который с сессией этого юзера будет слать спам пока может
	http.HandleFunc("/comment", func(w http.ResponseWriter, r *http.Request) {
		if !checkSession(r) {
			w.Write([]byte(loginFormTmplRaw))
			return
		}
		r.ParseForm()
		commentText := r.FormValue("comment")
		messages = append(messages, commentText)
		http.Redirect(w, r, "/", http.StatusFound)
	})

	// сервисный метод для очистки комментариев
	http.HandleFunc("/clear_comments", func(w http.ResponseWriter, r *http.Request) {
		if !checkSession(r) {
			w.Write([]byte(loginFormTmplRaw))
			return
		}

		messages = []string{}
		http.Redirect(w, r, "/", http.StatusFound)
	})

	// создаём сессию
	// не используйте этот подход в продакшене
	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		inputLogin := r.Form["login"][0]
		expiration := time.Now().Add(365 * 24 * time.Hour)

		sessionID := RandStringRunes(32)
		sessions[sessionID] = inputLogin

		cookie := http.Cookie{Name: "session_id", Value: sessionID, Expires: expiration}
		http.SetCookie(w, &cookie)
		http.Redirect(w, r, "/", http.StatusFound)
	})

	fmt.Println("starting server at :8080")
	http.ListenAndServe(":8080", nil)
}

//PanicOnErr panics on error
func PanicOnErr(err error) {
	if err != nil {
		panic(err)
	}
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
