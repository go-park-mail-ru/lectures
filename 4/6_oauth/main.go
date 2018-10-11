package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/vk"
)

const (
	APP_ID     = "6703763"
	APP_KEY    = "NJaZc5qmq5jnVu9eo7SO"
	APP_SECRET = "bc547f20bc547f20bc547f202cbc3235b3bbc54bc547f20e7ff0660a087e7cb533a64fe"
	API_URL    = "https://api.vk.com/method/users.get?fields=email,photo_50&access_token=%s&v=5.52"
)

type Response struct {
	Response []struct {
		FirstName string `json:"first_name"`
		Photo     string `json:"photo_50"`
	}
}

// https://oauth.vk.com/authorize?client_id=6703763&redirect_uri=http://127.0.0.1:8080/

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		code := r.FormValue("code")
		conf := oauth2.Config{
			ClientID:     APP_ID,
			ClientSecret: APP_KEY,
			RedirectURL:  "http://127.0.0.1:8080/",
			Endpoint:     vk.Endpoint,
		}

		token, err := conf.Exchange(ctx, code)
		if err != nil {
			log.Println("cannot exchange", err)
			w.Write([]byte("=("))
			return
		}

		client := conf.Client(ctx, token)
		resp, err := client.Get(fmt.Sprintf(API_URL, token.AccessToken))
		if err != nil {
			log.Println("cannot request data", err)
			w.Write([]byte("=("))
			return
		}

		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Println("cannot read buffer", err)
			w.Write([]byte("=("))
			return
		}

		data := &Response{}
		json.Unmarshal(body, data)

		w.Write([]byte(`
		<div>
			<img src="` + data.Response[0].Photo + `"/>
			` + data.Response[0].FirstName + `
		</div>
		`))
	})

	http.ListenAndServe(":8080", nil)
}
