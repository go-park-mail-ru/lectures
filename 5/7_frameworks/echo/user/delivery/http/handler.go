package http

import (
	"echo_example/middleware"
	"echo_example/user"
	"errors"

	"github.com/labstack/echo"
)

type userHandler struct {
	usecase user.Usecase
}

func NewUserHandler(e *echo.Echo, us user.Usecase) {
	handler := userHandler{usecase: us}

	e.GET("/user", handler.GetUser, middleware.PanicMiddleware)
	e.POST("/user", handler.CreateUser)

	e.GET("/users", handler.GetAllUsers)
}

func (h *userHandler) GetUser(ctx echo.Context) error {
	ctx.Error(errors.New("some error"))
}

func (h *userHandler) GetAllUsers(ctx echo.Context) error {
	return nil
}

func (h *userHandler) CreateUser(ctx echo.Context) error {
	return nil
}
