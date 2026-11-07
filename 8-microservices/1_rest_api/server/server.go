package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

var sessionManager *SessionManager

// curl -v -X POST http://localhost:8081/session -d "login=bob" -d "ua=Mozilla/5.0"
func Create(w http.ResponseWriter, r *http.Request) {
	returnInternalError := func(format string, args ...any) {
		fmt.Printf(format, args...)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("failed to create session: internal error\n"))
	}

	if err := r.ParseForm(); err != nil {
		returnInternalError("failed to create session: failed to parse form: %v\n", err)
		return
	}

	login := r.Form.Get("login")
	if login == "" {
		fmt.Println("cannot create session: login is empty")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("form value 'login' is empty\n"))
		return
	}

	session := Session{
		Login:     login,
		UserAgent: r.Form.Get("ua"),
	}

	var sessionID SessionID
	if err := sessionManager.Create(&session, &sessionID); err != nil {
		returnInternalError("failed to create session: session manager error: %v\n", err)
		return
	}

	respData, err := json.Marshal(sessionID)
	if err != nil {
		returnInternalError("failed to create session: marshal response error: %v\n", err)
		return
	}

	w.Write(respData)
}

// curl -v http://localhost:8081/session/100500qq
func Check(w http.ResponseWriter, r *http.Request) {
	sessionID := SessionID{
		ID: r.PathValue("session_id"),
	}
	session := Session{}

	if err := sessionManager.Check(&sessionID, &session); err != nil {
		fmt.Println("failed to check session:", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("failed to check session\n"))
		return
	}

	if session.Login == "" {
		fmt.Println("failed to check session: session not found")
		w.Write([]byte("session not found\n"))
		return
	}

	respData, err := json.Marshal(session)
	if err != nil {
		fmt.Println("failed to check session: marshal response error", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("failed to check session\n"))
		return
	}

	w.Write(respData)
}

// curl -v -X DELETE http://localhost:8081/session/100500qq
func Delete(w http.ResponseWriter, r *http.Request) {
	sessionID := SessionID{
		ID: r.PathValue("session_id"),
	}

	var status int
	if err := sessionManager.Delete(&sessionID, &status); err != nil {
		fmt.Println("failed to delete session:", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("failed to delete session\n"))
		return
	}

	w.Write([]byte("delete session with id=" + sessionID.ID + "\n"))
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
	sessionManager = NewSessionManager()

	mux := http.NewServeMux()
	mux.HandleFunc("POST /session", Create)
	mux.HandleFunc("GET /session/{session_id}", Check)
	mux.HandleFunc("DELETE /session/{session_id}", Delete)
	muxWithMiddleware := accessLogMiddleware(mux)

	fmt.Println("starting server at :8081")
	log.Fatal(http.ListenAndServe(":8081", muxWithMiddleware))
}
