package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/rpc"
	"net/rpc/jsonrpc"
	"time"
)

type HttpConn struct {
	in  io.Reader
	out io.Writer
}

func (c *HttpConn) Read(p []byte) (n int, err error) {
	return c.in.Read(p)
}

func (c *HttpConn) Write(d []byte) (n int, err error) {
	return c.out.Write(d)
}

func (c *HttpConn) Close() error {
	return nil
}

/*
Структура запроса:
{
   "jsonrpc":"2.0", // Версия json-rpc
   "id":1,          // ID запроса, генерируемый клиентом. В ответе сервера будет этот же ID
   "method":"SessionManager.Create",
   "params":[       // Аргументы запроса
      {
         "login":"rvasily",
         "useragent":"chrome"
      }
   ]
}

curl -v -X POST -H "Content-Type: application/json" \
    -d '{"jsonrpc":"2.0", "id": 1, "method": "SessionManager.Create", "params": [{"login":"rvasily", "useragent": "chrome"}]}'\
	http://localhost:8081/rpc

curl -v -H "Content-Type: application/json" \
    -d '{"jsonrpc":"2.0", "id": 2, "method": "SessionManager.Check", "params": [{"id":"100500qq"}]}' \
	http://localhost:8081/rpc

*/

type Handler struct {
	rpcServer *rpc.Server
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")

	serverCodec := jsonrpc.NewServerCodec(&HttpConn{
		in:  r.Body,
		out: w,
	})

	if err := h.rpcServer.ServeRequest(serverCodec); err != nil {
		log.Printf("error while serving JSON request: %v\n", err)
		http.Error(w, `{"error":"cant serve request"}`, http.StatusInternalServerError)
	}
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
	sessionManager := NewSessionManager()

	server := rpc.NewServer()
	server.Register(sessionManager)

	sessionHandler := &Handler{
		rpcServer: server,
	}
	http.Handle("/rpc", accessLogMiddleware(sessionHandler))

	fmt.Println("starting server at :8081")
	log.Fatal(http.ListenAndServe(":8081", nil))
}
