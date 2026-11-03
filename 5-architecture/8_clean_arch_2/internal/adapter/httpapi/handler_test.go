package httpapi_test

import (
	"io"
	"lectures/cleanarchreview/internal/adapter/httpapi"
	"lectures/cleanarchreview/internal/adapter/memory"
	"lectures/cleanarchreview/internal/entity"
	"lectures/cleanarchreview/internal/platform"
	"lectures/cleanarchreview/internal/usecase/article"
	"log/slog"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestArticleHTTPFlow(t *testing.T) {
	a, err := entity.NewAuthor(1, "Автор")
	if err != nil {
		t.Fatal(err)
	}
	s := article.New(memory.NewArticles(), memory.NewAuthors(a), platform.Clock{}, &platform.IDs{})
	h := httpapi.New(s, slog.New(slog.NewTextHandler(io.Discard, nil)), time.Second).Routes()
	steps := []struct {
		method, path, body string
		status             int
	}{
		{"POST", "/articles", `{"author_id":1,"title":"Статья","content":"Текст"}`, 201},
		{"POST", "/articles", `{"author_id":1,"title":"Статья","content":"Текст"}`, 409},
		{"GET", "/articles/1", "", 200},
		{"GET", "/articles?num=abc", "", 400},
		{"POST", "/articles", `{"author_id":1,"title":" ","content":"Текст"}`, 400},
		{"POST", "/articles", `{"author_id":99,"title":"Другая","content":"Текст"}`, 404},
		{"PUT", "/articles/1", `{"title":"Изменённая","content":"Новый текст"}`, 200},
		{"DELETE", "/articles/1", "", 204},
		{"GET", "/articles/1", "", 404},
	}
	for _, step := range steps {
		w := httptest.NewRecorder()
		h.ServeHTTP(w, httptest.NewRequest(step.method, step.path, strings.NewReader(step.body)))
		if w.Code != step.status {
			t.Fatalf("%s %s: got %d, want %d: %s", step.method, step.path, w.Code, step.status, w.Body.String())
		}
	}
}
