package main

import (
	"fmt"
	"net/http"

	"github.com/getsentry/sentry-go"
	sentryecho "github.com/getsentry/sentry-go/echo"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	_ = sentry.Init(sentry.ClientOptions{
		Dsn: "https://73eda8ec507b464caeb9d8fd331bc9ce@o427877.ingest.sentry.io/5372588",
		BeforeSend: func(event *sentry.Event, hint *sentry.EventHint) *sentry.Event {
			fmt.Println(event)
			return event
		},
		Debug:            true,
		AttachStacktrace: true,
	})

	app := echo.New()

	app.Use(middleware.Logger())
	app.Use(middleware.Recover())

	app.Use(sentryecho.New(sentryecho.Options{
		Repanic: true,
	}))

	app.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			if hub := sentryecho.GetHubFromContext(ctx); hub != nil {
				hub.Scope().SetTag("someRandomTag", "maybeYouNeedIt")
			}
			return next(ctx)
		}
	})

	app.GET("/", func(ctx echo.Context) error {
		if hub := sentryecho.GetHubFromContext(ctx); hub != nil {
			hub.WithScope(func(scope *sentry.Scope) {
				scope.SetExtra("unwantedQuery", "someQueryDataMaybe")
				hub.CaptureMessage("User provided unwanted query string, but we recovered just fine")
			})
		}
		return ctx.String(http.StatusOK, "Hello, World!")
	})

	app.GET("/foo", func(ctx echo.Context) error {
		panic("y tho")
	})

	app.Logger.Fatal(app.Start(":8080"))
}
