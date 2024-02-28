package main

import (
	"math/rand"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type User struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

var (
	letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
)

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

type MyHandler struct {
	sessions map[string]uint
	users    map[string]*User
}

func NewMyHandler() *MyHandler {
	return &MyHandler{
		sessions: make(map[string]uint, 10),
		users: map[string]*User{
			"rvasily": {1, "rvasily", "love"},
		},
	}
}

// http://127.0.0.1:8081/login?login=rvasily&password=love

func (api *MyHandler) Login(w http.ResponseWriter, r *http.Request) {

	user, ok := api.users[r.FormValue("login")]
	if !ok {
		http.Error(w, `no user`, 404)
		return
	}

	if user.Password != r.FormValue("password") {
		http.Error(w, `bad pass`, 400)
		return
	}

	SID := RandStringRunes(32)

	api.sessions[SID] = user.ID

	cookie := &http.Cookie{
		Name:    "session_id",
		Value:   SID,
		Path:    "/",
		Expires: time.Now().Add(10 * time.Hour),
	}
	http.SetCookie(w, cookie)
	w.Write([]byte(SID))

}

func (api *MyHandler) Logout(w http.ResponseWriter, r *http.Request) {

	session, err := r.Cookie("session_id")
	if err == http.ErrNoCookie {
		http.Error(w, `no sess`, 401)
		return
	}

	if _, ok := api.sessions[session.Value]; !ok {
		http.Error(w, `no sess`, 401)
		return
	}

	delete(api.sessions, session.Value)

	session.Expires = time.Now().AddDate(0, 0, -1)
	http.SetCookie(w, session)
}

func (api *MyHandler) Root(w http.ResponseWriter, r *http.Request) {
	authorized := false
	session, err := r.Cookie("session_id")
	if err == nil && session != nil {
		_, authorized = api.sessions[session.Value]
	}

	if authorized {
		w.Write([]byte("autrorized"))
	} else {
		w.Write([]byte("not autrorized"))
	}
}

func main() {
	r := mux.NewRouter()

	api := NewMyHandler()
	r.HandleFunc("/", api.Root)
	r.HandleFunc("/login", api.Login)
	r.HandleFunc("/logout", api.Logout)

	http.ListenAndServe(":8081", r)
}
