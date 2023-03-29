package main

import (
	"log"
	"net/http"
	"net/url"

	"github.com/casbin/casbin"
)

// https://casbin.org/en/editor
func main() {
	ur := &UsersRepo{
		db: map[string]*User{
			"admin@mail.ru": {1, "admin", "555"},
			"user@mail.ru":  {2, "member", "123"},
		},
	}
	e, err := casbin.NewEnforcerSafe("basic_model.conf", "basic_policy.csv")
	if err != nil {
		panic(err)
	}

	r := http.NewServeMux()
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("root page")) })
	r.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("you just logged int")) })
	r.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("you just logged out")) })

	r.HandleFunc("/member/profile", func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("here is your profile")) })
	r.HandleFunc("/member/profiles", func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("here is your profile")) })

	r.HandleFunc("/admin/user", func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("user edditing")) })
	r.HandleFunc("/admin/users", func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("users list")) })

	mux := PermissionsMiddleware(e, ur, r)

	http.ListenAndServe(":8080", mux)
}

func PermissionsMiddleware(e *casbin.Enforcer, ur *UsersRepo, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		v, _ := url.ParseQuery(r.URL.RawQuery)
		email := v.Get("email")
		// check auth and cookie/token ...

		role := "anonymous"
		user := ur.GetUser(email)
		if user != nil {
			role = user.Role
		}

		res, _ := e.EnforceSafe(role, r.URL.Path, r.Method)
		log.Printf("path=%s role=%s access=%v", r.URL.Path, role, res)
		if res {
			next.ServeHTTP(w, r)
		} else {
			http.Error(w, "forbidden", http.StatusForbidden)
			return
		}
	})

}

type User struct {
	UserID   int
	Role     string
	Password string
}
type UsersRepo struct {
	db map[string]*User
}

func (ur *UsersRepo) GetUser(email string) *User {
	return ur.db[email]
}
