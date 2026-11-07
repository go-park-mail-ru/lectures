// main собирает зависимости и запускает сервер — это внешний слой приложения.
package main

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	httpapi "lectures/cleanarch3/internal/controller/http"
	"lectures/cleanarch3/internal/repository/memory"
	"lectures/cleanarch3/internal/usecase"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	if err := run(logger); err != nil {
		logger.Error("application stopped", "error", err)
		os.Exit(1)
	}
}

func run(logger *slog.Logger) error {
	users := memory.NewUserRepository()

	// Dependency Injection: создаём объекты и передаём их через конструкторы.
	userUC := usecase.NewUserUseCase(users)
	handler := httpapi.NewUserHandler(userUC, logger, 3*time.Second)
	address := os.Getenv("HTTP_ADDR")
	if address == "" {
		address = ":9090"
	}
	server := &http.Server{
		Addr:              address,
		Handler:           handler.Routes(),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}
	logger.Info("listening", "address", server.Addr)
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("serve HTTP: %w", err)
	}
	return nil
}
