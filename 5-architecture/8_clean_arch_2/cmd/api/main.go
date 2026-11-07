// main — composition root: единственное место выбора и соединения реализаций.
package main

import (
	"errors"
	"lectures/cleanarchreview/internal/adapter/httpapi"
	"lectures/cleanarchreview/internal/adapter/memory"
	"lectures/cleanarchreview/internal/entity"
	"lectures/cleanarchreview/internal/platform"
	"lectures/cleanarchreview/internal/usecase/article"
	"log/slog"
	"net/http"
	"os"
	"time"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	author, err := entity.NewAuthor(1, "Учебный автор")
	if err != nil {
		logger.Error("seed author", "error", err)
		os.Exit(1)
	}
	articles := memory.NewArticles()
	authors := memory.NewAuthors(author)
	service := article.New(articles, authors, platform.Clock{}, &platform.IDs{})
	handler := httpapi.New(service, logger, 3*time.Second)
	server := &http.Server{Addr: ":9090", Handler: handler.Routes(), ReadHeaderTimeout: 5 * time.Second}
	logger.Info("listening", "address", server.Addr)
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		logger.Error("serve", "error", err)
		os.Exit(1)
	}
}
