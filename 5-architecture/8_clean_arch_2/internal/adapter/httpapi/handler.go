package httpapi

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"lectures/cleanarchreview/internal/entity"
	"lectures/cleanarchreview/internal/usecase/article"
	"log/slog"
	"net/http"
	"strconv"
	"time"
)

// Входной порт объявлен у потребителя. Service реализует его неявно.
type Usecases interface {
	Create(context.Context, article.CreateInput) (article.View, error)
	GetByID(context.Context, int64) (article.View, error)
	Update(context.Context, article.UpdateInput) (article.View, error)
	Fetch(context.Context, int64, int) (article.Page, error)
	Delete(context.Context, int64) error
}

type Handler struct {
	usecases Usecases
	logger   *slog.Logger
	timeout  time.Duration
}

func New(u Usecases, logger *slog.Logger, timeout time.Duration) *Handler {
	return &Handler{usecases: u, logger: logger, timeout: timeout}
}

func (h *Handler) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /articles", h.create)
	mux.HandleFunc("GET /articles", h.fetch)
	mux.HandleFunc("GET /articles/{id}", h.get)
	mux.HandleFunc("PUT /articles/{id}", h.update)
	mux.HandleFunc("DELETE /articles/{id}", h.delete)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), h.timeout)
		defer cancel()
		mux.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (h *Handler) create(w http.ResponseWriter, r *http.Request) {
	var in createRequest
	if err := decode(w, r, &in); err != nil {
		h.fail(w, r, err)
		return
	}
	v, err := h.usecases.Create(r.Context(), article.CreateInput{AuthorID: in.AuthorID, Title: in.Title, Content: in.Content})
	if err != nil {
		h.fail(w, r, err)
		return
	}
	w.Header().Set("Location", "/articles/"+strconv.FormatInt(v.Article.ID(), 10))
	h.respond(w, http.StatusCreated, toResponse(v))
}

func (h *Handler) get(w http.ResponseWriter, r *http.Request) {
	id, err := pathID(r)
	if err != nil {
		h.fail(w, r, err)
		return
	}
	v, err := h.usecases.GetByID(r.Context(), id)
	if err != nil {
		h.fail(w, r, err)
		return
	}
	h.respond(w, http.StatusOK, toResponse(v))
}

func (h *Handler) update(w http.ResponseWriter, r *http.Request) {
	id, err := pathID(r)
	if err != nil {
		h.fail(w, r, err)
		return
	}
	var in updateRequest
	if err := decode(w, r, &in); err != nil {
		h.fail(w, r, err)
		return
	}
	v, err := h.usecases.Update(r.Context(), article.UpdateInput{ID: id, Title: in.Title, Content: in.Content})
	if err != nil {
		h.fail(w, r, err)
		return
	}
	h.respond(w, http.StatusOK, toResponse(v))
}

func (h *Handler) fetch(w http.ResponseWriter, r *http.Request) {
	after, err := queryInt(r, "after_id")
	if err != nil {
		h.fail(w, r, err)
		return
	}
	limit, err := queryInt(r, "num")
	if err != nil || limit < 0 || limit > 100 {
		h.fail(w, r, article.ErrInvalidInput)
		return
	}
	page, err := h.usecases.Fetch(r.Context(), after, int(limit))
	if err != nil {
		h.fail(w, r, err)
		return
	}
	items := make([]articleResponse, 0, len(page.Items))
	for _, v := range page.Items {
		items = append(items, toResponse(v))
	}
	if page.NextID != 0 {
		w.Header().Set("X-Next-ID", strconv.FormatInt(page.NextID, 10))
	}
	h.respond(w, http.StatusOK, items)
}

func (h *Handler) delete(w http.ResponseWriter, r *http.Request) {
	id, err := pathID(r)
	if err != nil {
		h.fail(w, r, err)
		return
	}
	if err := h.usecases.Delete(r.Context(), id); err != nil {
		h.fail(w, r, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func pathID(r *http.Request) (int64, error) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil || id <= 0 {
		return 0, article.ErrInvalidInput
	}
	return id, nil
}
func queryInt(r *http.Request, name string) (int64, error) {
	v := r.URL.Query().Get(name)
	if v == "" {
		return 0, nil
	}
	n, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		return 0, article.ErrInvalidInput
	}
	return n, nil
}
func decode(w http.ResponseWriter, r *http.Request, out any) error {
	dec := json.NewDecoder(http.MaxBytesReader(w, r.Body, 1<<20))
	dec.DisallowUnknownFields()
	if err := dec.Decode(out); err != nil {
		return article.ErrInvalidInput
	}
	// В теле допустим ровно один JSON-документ.
	if err := dec.Decode(new(any)); err != io.EOF {
		return article.ErrInvalidInput
	}
	return nil
}
func (h *Handler) respond(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(body); err != nil {
		h.logger.Error("encode response", "error", err)
	}
}
func (h *Handler) fail(w http.ResponseWriter, r *http.Request, err error) {
	status, message := http.StatusInternalServerError, "internal error"
	switch {
	case errors.Is(err, article.ErrInvalidInput), errors.Is(err, entity.ErrInvalidArticle), errors.Is(err, entity.ErrInvalidAuthor):
		status, message = http.StatusBadRequest, "invalid input"
	case errors.Is(err, article.ErrNotFound):
		status, message = http.StatusNotFound, "not found"
	case errors.Is(err, article.ErrConflict):
		status, message = http.StatusConflict, "article title already exists"
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
