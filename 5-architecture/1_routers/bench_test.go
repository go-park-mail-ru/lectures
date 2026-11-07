package routerbench

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/buaazp/fasthttprouter"
	"github.com/gorilla/mux"
	"github.com/julienschmidt/httprouter"
	"github.com/valyala/fasthttp"
)

func okHandler(w http.ResponseWriter, r *http.Request)                            {}
func okHandlerParams(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {}
func okHandlerFasthttp(ctx *fasthttp.RequestCtx)                                  {}

// go test -bench=. -benchmem
func BenchmarkFasthttprouter(b *testing.B) {
	r := fasthttprouter.New()
	r.GET("/users/:id", okHandlerFasthttp)

	ctx := &fasthttp.RequestCtx{}
	ctx.Request.SetRequestURI("/users/42")
	ctx.Request.Header.SetMethod(http.MethodGet)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r.Handler(ctx)
	}
}

func BenchmarkHttprouter(b *testing.B) {
	r := httprouter.New()
	r.GET("/users/:id", okHandlerParams)

	req := httptest.NewRequest(http.MethodGet, "/users/42", nil)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r.ServeHTTP(httptest.NewRecorder(), req)
	}
}

func BenchmarkStdServeMux(b *testing.B) {
	mux := http.NewServeMux()
	mux.HandleFunc("/users/", okHandler)

	req := httptest.NewRequest(http.MethodGet, "/users/42", nil)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		mux.ServeHTTP(httptest.NewRecorder(), req)
	}
}

func BenchmarkGorillaMux(b *testing.B) {
	r := mux.NewRouter()
	r.HandleFunc("/users/{id}", okHandler)

	req := httptest.NewRequest(http.MethodGet, "/users/42", nil)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r.ServeHTTP(httptest.NewRecorder(), req)
	}
}
