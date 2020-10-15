package middleware

import (
	"context"
	"fmt"
	"net/http"

	"crudapp/internal/pkg/models"
	"crudapp/internal/pkg/session"
)

var (
	noAuthUrls = map[string]struct{}{
		"/login": struct{}{},
	}
	noSessUrls = map[string]struct{}{
		"/": struct{}{},
	}
)

func Auth(sm *session.SessionsManager, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("auth middleware")
		if _, ok := noAuthUrls[r.URL.Path]; ok {
			next.ServeHTTP(w, r)
			return
		}
		sess, err := sm.Check(r)
		_, canbeWithouthSess := noSessUrls[r.URL.Path]
		if err != nil && !canbeWithouthSess {
			fmt.Println("no auth")
			http.Redirect(w, r, "/", 302)
			return
		}
		ctx := context.WithValue(r.Context(), models.SessionKey, sess)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
