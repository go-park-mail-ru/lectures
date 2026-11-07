package main

import (
	"context"
	"testing"
)

func TestRun(t *testing.T) {

	sessionManager = NewSessionManager()

	ctx := context.Background()

	// создаем сессию
	sessId, err := sessionManager.Create(
		ctx,
		&Session{
			Login:     "anton",
			UserAgent: "safari",
		})
	t.Log("sessId", sessId, err)

	// проверяем сессию
	sess := sessionManager.Check(
		ctx,
		&SessionID{
			ID: sessId.ID,
		})
	t.Log("sess", sess)

	// удаляем сессию
	sessionManager.Delete(
		ctx,
		&SessionID{
			ID: sessId.ID,
		})

	// проверяем еще раз
	sess = sessionManager.Check(
		ctx,
		&SessionID{
			ID: sessId.ID,
		})
	t.Log("sess", sess)
	t.Fail()
}
