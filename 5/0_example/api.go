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

var (
	userDB = map[string]User{}
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

func main() {
	http.HandleFunc("/user", createUserHandler)
	http.HandleFunc("/users", allUsersHandler)

	log.Println("listening on 127.0.0.1:8000")
	log.Fatal(http.ListenAndServe("127.0.0.1:8000", nil))
}
