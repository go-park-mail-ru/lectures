package main

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

var privateAlbums = map[string]struct{}{
	"selfies": struct{}{},
	"nudes":   struct{}{},
}

var privateAlbumsACL = map[string]map[int]struct{}{
	"selfies": map[int]struct{}{
		1: struct{}{},
		2: struct{}{},
		3: struct{}{},
	},
	"nudes": map[int]struct{}{
		2: struct{}{},
	},
}

func checkACL(albumName string, uid int) bool {
	if _, ok := privateAlbums[albumName]; !ok {
		return true
	}
	allowedUsers, ok := privateAlbumsACL[albumName]
	if !ok {
		return true
	}
	_, ok = allowedUsers[uid]
	return ok
}

func getSession(r string) int {
	uid, _ := strconv.Atoi(r)
	return uid
}

func main() {

	http.HandleFunc("/auth", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("AUTH request", r)

		// fmt.Println("X-Original-URI", r.Header.Get("X-Original-URI"))
		req, _ := url.Parse(r.Header.Get("X-Original-URI"))
		// fmt.Println("X-User-ID", req, err)
		// fmt.Println("Req Data", req.Path, " -- ", req.Query().Get("user_id"))

		uid := getSession(req.Query().Get("user_id"))

		str := req.Path
		albumName := strings.ReplaceAll(str, "/albums/", "")
		fmt.Println("PARAMS", albumName, uid)

		if !checkACL(albumName, uid) {
			fmt.Println("ACL failed", albumName, uid)
			http.Error(w, "", 403)
		}

		w.Header().Set("WWW-Authenticate", req.Query().Get("user_id"))
		fmt.Println("ACL OK", albumName, uid)
		http.Error(w, "", 200)
	})

	http.HandleFunc("/albums/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("incoming request", r)
		fmt.Fprintln(w, "hi", r.Header.Get("WWW-Authenticate"))
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("ROOT incoming request", r)
		fmt.Fprintln(w, "hi", r.Header.Get("WWW-Authenticate"))
	})

	fmt.Println("start server at :8080")
	http.ListenAndServe(":8080", nil)
}
