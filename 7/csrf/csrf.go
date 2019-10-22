// CSRF - это выполнение какие-то действий на сайте от имени другого пользователя

// данный пример ИСКУССТВЕННЫЙ, чтобы показать как проявляется CSRF
// используйте пакет html/template
// он автоматичски экранирует все входящие данные с учетом контекста
// подрбонее https://golang.org/pkg/html/template/

package main

import (
	"net/http"
	// "html/template"
	"fmt"
	"math/rand"
	"strconv"
	"text/template" // надо заменить text/template на html/template чтобы по-умоллчанию было правильное экранирование
	"time"
)

var sessions = map[string]string{}
var cnt = 1

type Msg struct {
	ID      int
	Message string
	Rating  int
}

var messages = map[int]*Msg{}

var loginFormTmplRaw = `<html><body>
	<form action="/login" method="post">
		Login: <input type="text" name="login" value="DefaultUserName">
		Password: <input type="password" name="password" value="anypass">
		<input type="submit" value="Login">
	</form>
</body></html>`

var messagesTmpl = `<html>
<head>
<script>
	function rateComment(id, vote) {
		var request = new XMLHttpRequest();
		request.open('POST', '/rate?id='+id+"&vote="+(vote > 0 ? "up" : "down"), true);

		request.onload = function() {
		    var resp = JSON.parse(request.responseText);
			console.log(resp, resp.id)
			console.log('#rating-'+resp.id)
			console.log(document.querySelector('#rating-'+resp.id))
			document.querySelector('#rating-'+resp.id).innerHTML = resp.rating;
		};
		request.send();
	}
</script>
</head>
<body>
	&lt;img src=&quot;/rate?id=1&amp;vote=up&quot;&gt;
	<br />
	<br />

	<form action="/comment" method="post">
		<textarea name="comment"></textarea><br />
		<input type="submit" value="Comment">
	</form>
	<br />
	
    {{range $idx, $var := .Messages}}
		<div style="border: 1px solid black; padding: 5px; margin: 5px;">
			<button onclick="rateComment({{$var.ID}}, 1)">&uarr;</button>
			<span id="rating-{{$var.ID}}">{{$var.Rating}}</span>
			<button onclick="rateComment({{$var.ID}}, -1)">&darr;</button>
			&nbsp;

			<!-- text/template по-умолчанию ничего не экранируется --> 
			<!-- html/template по-умолчанию будет экранировать --> 
			{{$var.Message}}
		</div>
    {{end}}
</body></html>`

func main() {

	tmpl := template.New("main")
	tmpl, _ = tmpl.Parse(messagesTmpl)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if !checkSession(r) {
			w.Write([]byte(loginFormTmplRaw))
			return
		}
		//выводим комментарии + форму отправки
		tmpl.Execute(w, struct {
			Messages map[int]*Msg
		}{
			Messages: messages,
		})
	})

	// добавление комментария
	// добавим комментарий c текстом
	/*
		<img src="/rate?id=1&vote=up">
	*/
	// это выведет на экран куки сайта. дальше с ними можно сделать всё что угодно - например отправить ан внешний сервис, который с сессией этого юзера будет слать спам пока может
	http.HandleFunc("/comment", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		commentText := r.FormValue("comment")
		id := cnt
		messages[id] = &Msg{
			ID:      id,
			Message: commentText,
			Rating:  0,
		}
		cnt++
		http.Redirect(w, r, "/", http.StatusFound)
	})

	// функция для изменения рейтинга
	// тут происхрдит CSRF т.к. <img который мы вставили в комменте выше вызывается каждым пользователем который его видит без его ведома
	http.HandleFunc("/rate", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		emptyResponse := []byte(`{"id":0, "rating":0}`)
		if !checkSession(r) || r.Method == http.MethodGet {
			w.Write([]byte(emptyResponse))
			return
		}

		id, _ := strconv.Atoi(r.URL.Query().Get("id"))
		vote := r.URL.Query().Get("vote")

		if msg, ok := messages[id]; ok {
			if vote == "up" {
				msg.Rating++
			} else if vote == "down" {
				msg.Rating--
			}
			w.Write([]byte(fmt.Sprintf(`{"id":%d, "rating":%d}`, msg.ID, msg.Rating)))
		} else {
			w.Write([]byte(emptyResponse))
		}
	})

	// сервисный метод для очистки комментариев
	http.HandleFunc("/clear_comments", func(w http.ResponseWriter, r *http.Request) {
		if !checkSession(r) {
			w.Write([]byte(loginFormTmplRaw))
			return
		}
		messages = map[int]*Msg{}
		http.Redirect(w, r, "/", http.StatusFound)
	})

	// создаём сессию
	// не используйте эитот подход в продакшене
	http.HandleFunc("/login", loginHandler)

	fmt.Println("starting server at :8080")
	http.ListenAndServe(":8080", nil)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	inputLogin := r.Form["login"][0]
	expiration := time.Now().Add(365 * 24 * time.Hour)

	sessionID := RandStringRunes(32)
	sessions[sessionID] = inputLogin

	cookie := http.Cookie{Name: "session_id", Value: sessionID, Expires: expiration}
	http.SetCookie(w, &cookie)
	http.Redirect(w, r, "/", http.StatusFound)
}

func checkSession(r *http.Request) bool {
	// обработка сессии
	// не используйте эитот подход в продакшене
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
