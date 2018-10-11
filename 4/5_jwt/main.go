package main

import (
	"fmt"
	"net/http"

	"github.com/dgrijalva/jwt-go"
)

var SECRET = []byte("myawesomesecret")

func main() {
	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"username": r.FormValue("username"),
		})

		str, err := token.SignedString(SECRET)
		if err != nil {
			w.Write([]byte("=(" + err.Error()))
			return
		}

		cookie := &http.Cookie{
			Name:  "session_id",
			Value: str,
		}

		http.SetCookie(w, cookie)
		w.Write([]byte(str))
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_id")
		if err != nil {
			w.Write([]byte("=("))
			return
		}
		token, err := jwt.Parse(cookie.Value, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}

			return SECRET, nil
		})

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			w.Write([]byte("hello" + claims["username"].(string)))
			return
		}
		w.Write([]byte("not authorized"))
		fmt.Println(err)

	})

	http.ListenAndServe(":9999", nil)
}
