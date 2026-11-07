// Package httpapi — входной адаптер: HTTP и JSON остаются на этой границе.
package httpapi

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"lectures/cleanarch3/internal/domain"
	"lectures/cleanarch3/internal/usecase"
)

// Входной интерфейс объявлен у потребителя — HTTP-контроллера.
// *usecase.UserUseCase удовлетворяет ему без implements или регистрации.
type UserUseCases interface {
	Register(context.Context, usecase.CreateUserInput) (domain.User, error)
	GetByID(context.Context, int64) (domain.User, error)
	List(context.Context, int64, int) (usecase.UserPage, error)
	Update(context.Context, usecase.UpdateUserInput) (domain.User, error)
	Delete(context.Context, int64) error
}

type UserHandler struct {
	users   UserUseCases
	logger  *slog.Logger
	timeout time.Duration
}

func NewUserHandler(users UserUseCases, logger *slog.Logger, timeout time.Duration) *UserHandler {
	return &UserHandler{users: users, logger: logger, timeout: timeout}
}

func (h *UserHandler) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /users", h.register)
	mux.HandleFunc("GET /users", h.list)
	mux.HandleFunc("GET /users/{id}", h.get)
	mux.HandleFunc("PUT /users/{id}", h.update)
	mux.HandleFunc("DELETE /users/{id}", h.delete)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), h.timeout)
		defer cancel()
		mux.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (h *UserHandler) register(w http.ResponseWriter, r *http.Request) {
	var in userRequest
	if err := decode(w, r, &in); err != nil {
		h.fail(w, r, err)
		return
	}
	u, err := h.users.Register(r.Context(), usecase.CreateUserInput{Name: in.Name, Email: in.Email})
	if err != nil {
		h.fail(w, r, err)
		return
	}
	w.Header().Set("Location", "/users/"+strconv.FormatInt(u.ID(), 10))
	h.respond(w, http.StatusCreated, toResponse(u))
}

func (h *UserHandler) get(w http.ResponseWriter, r *http.Request) {
	id, err := pathID(r)
	if err != nil {
		h.fail(w, r, err)
		return
	}
	u, err := h.users.GetByID(r.Context(), id)
	if err != nil {
		h.fail(w, r, err)
		return
	}
	h.respond(w, http.StatusOK, toResponse(u))
}

func (h *UserHandler) update(w http.ResponseWriter, r *http.Request) {
	id, err := pathID(r)
	if err != nil {
		h.fail(w, r, err)
		return
	}
	var in userRequest
	if err := decode(w, r, &in); err != nil {
		h.fail(w, r, err)
		return
	}
	u, err := h.users.Update(r.Context(), usecase.UpdateUserInput{ID: id, Name: in.Name, Email: in.Email})
	if err != nil {
		h.fail(w, r, err)
		return
	}
	h.respond(w, http.StatusOK, toResponse(u))
}

func (h *UserHandler) list(w http.ResponseWriter, r *http.Request) {
	afterID, err := queryInt(r, "after_id")
	if err != nil {
		h.fail(w, r, err)
		return
	}
	limit, err := queryInt(r, "num")
	if err != nil || limit < 0 || limit > 100 {
		h.fail(w, r, usecase.ErrInvalidInput)
		return
	}
	page, err := h.users.List(r.Context(), afterID, int(limit))
	if err != nil {
		h.fail(w, r, err)
		return
	}
	items := make([]userResponse, 0, len(page.Items))
	for _, u := range page.Items {
		items = append(items, toResponse(u))
	}
	if page.NextID != 0 {
		w.Header().Set("X-Next-ID", strconv.FormatInt(page.NextID, 10))
	}
	h.respond(w, http.StatusOK, items)
}

func (h *UserHandler) delete(w http.ResponseWriter, r *http.Request) {
	id, err := pathID(r)
	if err != nil {
		h.fail(w, r, err)
		return
	}
	if err := h.users.Delete(r.Context(), id); err != nil {
		h.fail(w, r, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func pathID(r *http.Request) (int64, error) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil || id <= 0 {
		return 0, usecase.ErrInvalidInput
	}
	return id, nil
}

func queryInt(r *http.Request, name string) (int64, error) {
	value := r.URL.Query().Get(name)
	if value == "" {
		return 0, nil
	}
	n, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return 0, usecase.ErrInvalidInput
	}
	return n, nil
}

func decode(w http.ResponseWriter, r *http.Request, out any) error {
	decoder := json.NewDecoder(http.MaxBytesReader(w, r.Body, 1<<20))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(out); err != nil {
		return usecase.ErrInvalidInput
	}
	// Принимаем ровно один JSON-документ, максимум 1 МБ.
	if err := decoder.Decode(new(any)); err != io.EOF {
		return usecase.ErrInvalidInput
	}
	return nil
}

func (h *UserHandler) respond(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(body); err != nil {
		h.logger.Error("encode response", "error", err)
	}
}

func (h *UserHandler) fail(w http.ResponseWriter, r *http.Request, err error) {
	status, message := http.StatusInternalServerError, "internal error"
	switch {
	case errors.Is(err, usecase.ErrInvalidInput), errors.Is(err, domain.ErrInvalidUser):
		status, message = http.StatusBadRequest, "invalid input"
	case errors.Is(err, usecase.ErrNotFound):
		status, message = http.StatusNotFound, "user not found"
	case errors.Is(err, usecase.ErrEmailTaken):
		status, message = http.StatusConflict, "email already taken"
	case errors.Is(err, context.DeadlineExceeded):
		status, message = http.StatusGatewayTimeout, "request timeout"
	}
	if status >= 500 {
		h.logger.ErrorContext(r.Context(), "request failed", "error", err)
	}
	h.respond(w, status, struct {
		Message string `json:"message"`
	}{message})
}
