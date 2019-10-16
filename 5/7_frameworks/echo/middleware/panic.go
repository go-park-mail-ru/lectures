package middleware

import (
	"fmt"

	"github.com/labstack/echo"
)

func PanicMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		defer func() {
			if err := recover(); err != nil {
				ctx.Error(fmt.Errorf("%s", err))
			}
		}()

		return next(ctx)
	}
}
