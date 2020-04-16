package middleware

import (
	"fmt"
	"math/rand"

	"github.com/labstack/echo"
)

func RequestIDMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		requestID := fmt.Sprintf("%016x", rand.Int())[:10]

		ctx.Logger().SetPrefix(fmt.Sprintf("%s rid=%s", ctx.Logger().Prefix(), requestID))
		return next(ctx)
	}
}
