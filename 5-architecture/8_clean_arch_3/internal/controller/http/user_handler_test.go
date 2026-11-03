package httpapi_test

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	httpapi "lectures/cleanarch3/internal/controller/http"
	"lectures/cleanarch3/internal/domain"
	"lectures/cleanarch3/internal/repository/memory"
	"lectures/cleanarch3/internal/usecase"
)

func TestUserHTTPFlow(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	uc := usecase.NewUserUseCase(memory.NewUserRepository())
	h := httpapi.NewUserHandler(uc, logger, time.Second).Routes()
	steps := []struct {
		method, path, body string
		status             int
		response, nextID   string
	}{
		{"GET", "/users", "", 200, `[]`, ""},
		{"POST", "/users", `{"name":" Анна ","email":"Anna@Example.com"}`, 201, `{"id":1,"name":"Анна","email":"anna@example.com"}`, ""},
		{"POST", "/users", `{"name":"Другая Анна","email":"ANNA@example.com"}`, 409, `{"message":"email already taken"}`, ""},
		{"POST", "/users", `{"name":"Борис","email":"boris@example.com"}`, 201, `{"id":2,"name":"Борис","email":"boris@example.com"}`, ""},
		{"GET", "/users?num=1", "", 200, `[{"id":1,"name":"Анна","email":"anna@example.com"}]`, "1"},
		{"GET", "/users?num=1&after_id=1", "", 200, `[{"id":2,"name":"Борис","email":"boris@example.com"}]`, ""},
		{"PUT", "/users/1", `{"name":"Анна","email":"boris@example.com"}`, 409, `{"message":"email already taken"}`, ""},
		{"GET", "/users/1", "", 200, `{"id":1,"name":"Анна","email":"anna@example.com"}`, ""},
		{"PUT", "/users/1", `{"name":"Анна Иванова","email":"new@example.com"}`, 200, `{"id":1,"name":"Анна Иванова","email":"new@example.com"}`, ""},
		{"DELETE", "/users/1", "", 204, "", ""},
		{"GET", "/users/1", "", 404, `{"message":"user not found"}`, ""},
		{"DELETE", "/users/1", "", 404, `{"message":"user not found"}`, ""},
		{"GET", "/users/abc", "", 400, `{"message":"invalid input"}`, ""},
		{"GET", "/users/0", "", 400, `{"message":"invalid input"}`, ""},
		{"GET", "/users?num=101", "", 400, `{"message":"invalid input"}`, ""},
		{"GET", "/users?after_id=-1", "", 400, `{"message":"invalid input"}`, ""},
		{"GET", "/users?num=abc", "", 400, `{"message":"invalid input"}`, ""},
		{"POST", "/users", `{"name":" ","email":"other@example.com"}`, 400, `{"message":"invalid input"}`, ""},
		{"POST", "/users", `{"name":"Анна","email":"invalid"}`, 400, `{"message":"invalid input"}`, ""},
		{"POST", "/users", `{"name":"Анна","email":"other@example.com","role":"admin"}`, 400, `{"message":"invalid input"}`, ""},
		{"POST", "/users", `{"name":"Анна","email":"other@example.com"} {}`, 400, `{"message":"invalid input"}`, ""},
		{"POST", "/users", `null`, 400, `{"message":"invalid input"}`, ""},
	}
	for _, step := range steps {
		w := httptest.NewRecorder()
		h.ServeHTTP(w, httptest.NewRequest(step.method, step.path, strings.NewReader(step.body)))
		if w.Code != step.status || strings.TrimSpace(w.Body.String()) != step.response {
			t.Fatalf("%s %s: status=%d, body=%s; want status=%d, body=%s",
				step.method, step.path, w.Code, w.Body.String(), step.status, step.response)
		}
		if got := w.Header().Get("X-Next-ID"); got != step.nextID {
			t.Fatalf("%s: next ID=%q, want %q", step.path, got, step.nextID)
		}
		if step.status == 201 && w.Header().Get("Location") == "" {
			t.Fatal("created response has no Location")
		}
	}
}

type failingUseCases struct {
	httpapi.UserUseCases
	err error
}

func (uc failingUseCases) GetByID(context.Context, int64) (domain.User, error) {
	return domain.User{}, uc.err
}

func TestInternalErrorDoesNotLeakStorageDetails(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	h := httpapi.NewUserHandler(failingUseCases{err: errors.New("private storage details")}, logger, time.Second).Routes()
	w := httptest.NewRecorder()
	h.ServeHTTP(w, httptest.NewRequest("GET", "/users/1", nil))
	if w.Code != 500 || strings.TrimSpace(w.Body.String()) != `{"message":"internal error"}` {
		t.Fatalf("unexpected internal error response: %d %s", w.Code, w.Body.String())
	}
}
