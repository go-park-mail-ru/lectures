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
	APP_ID     = "7065390"
	APP_KEY    = "cQZe3Vvo4mHotmetUdXK"
	APP_SECRET = "1bbf49951bbf49951bbf49953b1bd486bb11bbf1bbf4995468b3d76e2cb2114610654e0"
	API_URL    = "https://api.vk.com/method/users.get?fields=email,photo_50&access_token=%s&v=5.131"
)

type Response struct {
	Response []struct {
		FirstName string `json:"first_name"`
		Photo     string `json:"photo_50"`
	}
}

// https://oauth.vk.com/authorize?client_id=7065390&redirect_uri=http://localhost:8080/&response_type=code&scope=email

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		code := r.FormValue("code")
		conf := oauth2.Config{
			ClientID:     APP_ID,
			ClientSecret: APP_KEY,
			RedirectURL:  "http://localhost:8080/",
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
