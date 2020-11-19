package main

import (
	"github.com/pkg/errors"
	"log"
	"net/http"
	"time"

	"github.com/getsentry/sentry-go"
)

func firstError() error {
	return errors.New("error")
}

func returnError() error {
	return errors.Wrap(firstError(), "failed to do something")
}

func main() {
	err := sentry.Init(sentry.ClientOptions{
		Dsn:   "https://73eda8ec507b464caeb9d8fd331bc9ce@o427877.ingest.sentry.io/5372588",
		Debug: true,
	})
	if err != nil {
		log.Fatalf("sentry.Init: %s", err)
	}

	defer sentry.Flush(2 * time.Second)

	http.HandleFunc(
		"/500",
		func(writer http.ResponseWriter, request *http.Request) {
			sentry.CaptureException(returnError())
		},
	)

	http.ListenAndServe(":8080", nil)
}
