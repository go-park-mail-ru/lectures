package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/go-park-mail-ru/lectures/8-microservices/5_grpc/session"
)

var (
	loginFormTmpl = []byte(`
<html>
	<body>
	<form action="/login" method="post">
		Login: <input type="text" name="login">
		Password: <input type="password" name="password">
		<input type="submit" value="Login">
	</form>
	</body>
</html>
`)

	sessManager session.AuthCheckerClient
)

func checkSession(r *http.Request) (*session.Session, error) {
	cookieSessionID, err := r.Cookie("session_id")
	if err == http.ErrNoCookie {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	sess, err := sessManager.Check(
		r.Context(),
		&session.SessionID{
			ID: cookieSessionID.Value,
		})
	if err != nil {
		return nil, err
	}

	return sess, nil
}

func innerPage(w http.ResponseWriter, r *http.Request) {
	sess, err := checkSession(r)
	if err != nil {
		w.Write(loginFormTmpl)
		return
	}
	if sess == nil {
		w.Write(loginFormTmpl)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintln(w, "Welcome, "+sess.Login+" <br />")
	fmt.Fprintln(w, "Session ua: "+sess.Useragent+" <br />")
	fmt.Fprintln(w, `<a href="/logout">logout</a>`)
}

func loginPage(w http.ResponseWriter, r *http.Request) {
	inputLogin := r.FormValue("login")
	expiration := time.Now().Add(365 * 24 * time.Hour)

	sess, err := sessManager.Create(
		r.Context(),
		&session.Session{
			Login:     inputLogin,
			Useragent: r.UserAgent(),
		})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	cookie := http.Cookie{
		Name:    "session_id",
		Value:   sess.ID,
		Expires: expiration,
	}
	http.SetCookie(w, &cookie)
	http.Redirect(w, r, "/", http.StatusFound)
}

func logoutPage(w http.ResponseWriter, r *http.Request) {
	cookieSessionID, err := r.Cookie("session_id")
	if err == http.ErrNoCookie {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, _ = sessManager.Delete(
		r.Context(),
		&session.SessionID{
			ID: cookieSessionID.Value,
		})

	cookieSessionID.Expires = time.Now().AddDate(0, 0, -1)
	http.SetCookie(w, cookieSessionID)

	http.Redirect(w, r, "/", http.StatusFound)
}

func accessLogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("accessLogMiddleware", r.URL.Path)
		start := time.Now()

		next.ServeHTTP(w, r)

		tm := time.Now().Format(time.RFC3339)
		fmt.Printf("%s [%s: %s] req_time=%s\n\n", tm, r.Method, r.URL.Path, time.Since(start))
	})
}

func main() {
	grpcClient, err := grpc.NewClient(
		"127.0.0.1:8081",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("cant connect to grpc")
	}
	defer grpcClient.Close()

	sessManager = session.NewAuthCheckerClient(grpcClient)

	mux := http.NewServeMux()
	mux.HandleFunc("/", innerPage)
	mux.HandleFunc("/login", loginPage)
	mux.HandleFunc("/logout", logoutPage)
	muxWithMiddleware := accessLogMiddleware(mux)

	fmt.Println("starting server at :8080")
	log.Fatal(http.ListenAndServe(":8080", muxWithMiddleware))
}
