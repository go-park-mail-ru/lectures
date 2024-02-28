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
	APP_ID     = "51569264"
	APP_KEY    = "HGgEdwoK7Sf4DfwX57Hc"
	APP_SECRET = "ce7841cfce7841cfce7841cf4dcd6aa3bfcce78ce7841cfad8724af176c113b04175fbc"
	API_URL    = "https://api.vk.com/method/users.get?fields=email,photo_50&access_token=%s&v=5.131"
)

type Response struct {
	Response []struct {
		FirstName string `json:"first_name"`
		Photo     string `json:"photo_50"`
	}
}

var AUTH_URL = `https://oauth.vk.com/authorize?client_id=51569264&redirect_uri=http://localhost:8082/user/login_oauth&response_type=code&scope=email`

// https://oauth.vk.com/authorize?client_id=7065390&redirect_uri=http://localhost:8080/&response_type=code&scope=email

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		code := r.FormValue("code")
		// инициализируем конфиг для vk oauth2
		conf := oauth2.Config{
			ClientID:     APP_ID,
			ClientSecret: APP_KEY,
			RedirectURL:  "http://localhost:8082/user/login_oauth", // ссылка, на которую мы будем перенаправлены после авторизации
			Endpoint:     vk.Endpoint,
		}

		if code == "" {
			// если code пустой, юзер еще не авторизован, надо выдать ему url для авторизации
			w.Write([]byte(`
				<div>
					<a href="` + AUTH_URL + `"> auth</a>
				</div>
				`))
			return
		}

		// code есть, значит юзер прошел по url авторизации, надо выдать токен
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

	http.ListenAndServe(":8075", nil)
}
