package main

import (
	"github.com/go-park-mail-ru/lectures/5-architecture/7_frameworks/echo/middleware"
	userhttp "github.com/go-park-mail-ru/lectures/5-architecture/7_frameworks/echo/user/delivery/http"
	"github.com/go-park-mail-ru/lectures/5-architecture/7_frameworks/echo/user/repository"
	"github.com/go-park-mail-ru/lectures/5-architecture/7_frameworks/echo/user/usecase"

	"github.com/labstack/echo"
	echomiddleware "github.com/labstack/echo/middleware"
)

const listenAddr = "127.0.0.1:8080"

func main() {
	e := echo.New()

	e.Use(middleware.RequestIDMiddleware)
	e.Use(echomiddleware.Logger())
	e.Use(middleware.PanicMiddleware)

	e.HTTPErrorHandler = middleware.ErrorHandler

	userhttp.NewUserHandler(e, usecase.NewUserUsecase(repository.NewUserMemoryRepository()))

	e.Logger.Warnf("start listening on %s", listenAddr)
	err := e.Start("127.0.0.1:8080")
	if err != nil {
		e.Logger.Errorf("server error: %s", err)
	}

	e.Logger.Warnf("shutdown")
}
