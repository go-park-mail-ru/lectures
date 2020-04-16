package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
)

type User struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type TODO struct {
	Text string `json:"text"`
}

var (
	userDB = map[string]User{
		"predefined": {
			Login:    "predefined",
			Password: "123",
		},
	}
	sessionDB = map[string]string{
		"predefined": "predefined",
	}
	todoDB = map[string][]TODO{
		"predefined": {{Text: "first thing"}, {Text: "second thing"}},
	}
)

func genRequestID() string {
	return fmt.Sprintf("%016x", rand.Int())[:10]
}

func createUserHandler(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("panic during request: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("internal error"))
		}
	}()

	requestID := genRequestID()
	log.Printf("incoming request [rid=%s] %s %q", requestID, r.Method, r.URL)

	switch r.Method {
	// curl 'http://127.0.0.1:8080/user' --cookie "session=78629a0f5f3f164f"
	case http.MethodGet:
		sessionCookie, err := r.Cookie("session")
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("unauthorized"))
			return
		}

		username := sessionDB[sessionCookie.Value]
		user, ok := userDB[username]
		if !ok {
			panic("user not exists")
		}

		res, _ := json.Marshal(user)

		w.WriteHeader(http.StatusOK)
		w.Write(res)

	// curl 'http://127.0.0.1:8080/user' -X POST -v -d '{"login": "user", "password": "123"}'
	case http.MethodPost:
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			panic("request body read failure")
		}

		var user User
		err = json.Unmarshal(body, &user)
		if err != nil {
			log.Printf("[rid=%s] unmarshal failure: %s", requestID, err)
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("bad body"))
			return
		}

		userDB[user.Login] = user

		sessionID := fmt.Sprintf("%016x", rand.Int())
		sessionDB[sessionID] = user.Login

		http.SetCookie(w, &http.Cookie{
			Name:  "session",
			Value: sessionID,
			Path:  "/",
		})
		w.WriteHeader(200)

	default:
		w.WriteHeader(http.StatusNotFound)
	}
}

func allUsersHandler(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("panic during request: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("internal error"))
		}
	}()

	requestID := genRequestID()
	log.Printf("incoming request [rid=%s] %s %q", requestID, r.Method, r.URL)

	switch r.Method {
	case http.MethodGet:
		response, err := json.Marshal(userDB)
		if err != nil {
			panic(err)
		}

		w.WriteHeader(200)
		w.Write(response)

	default:
		w.WriteHeader(http.StatusNotFound)
	}
}

func todoHandler(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("panic during request: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("internal error"))
		}
	}()

	requestID := genRequestID()
	log.Printf("incoming request [rid=%s] %s %q", requestID, r.Method, r.URL)

	sessionCookie, err := r.Cookie("session")
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("unauthorized"))
		return
	}

	username := sessionDB[sessionCookie.Value]

	switch r.Method {
	case http.MethodGet:
		todos, _ := json.Marshal(todoDB[username])
		w.WriteHeader(http.StatusOK)
		w.Write(todos)

	case http.MethodPost:
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			panic("request body read failure")
		}

		var todo TODO
		json.Unmarshal(body, &todo)
		todoDB[username] = append(todoDB[username], todo)

		w.WriteHeader(200)
		w.Write([]byte("ok"))

	default:
		w.WriteHeader(404)
		w.Write([]byte("not exists"))
	}
}

func main() {
	userMux := http.NewServeMux()

	userMux.HandleFunc("/user", createUserHandler)
	userMux.HandleFunc("/users", allUsersHandler)
	userMux.HandleFunc("/todo", todoHandler)

	log.Println("listening on 127.0.0.1:8080")
	log.Fatal(http.ListenAndServe("127.0.0.1:8080", userMux))
}
