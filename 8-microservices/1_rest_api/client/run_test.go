package main

import (
	"testing"
)

func TestRun(t *testing.T) {

	sessionManager = NewSessionManager()

	// создаем сессию
	sessId, err := sessionManager.Create(
		&Session{
			Login:     "anton",
			UserAgent: "safari",
		})
	t.Log("sessId", sessId, err)

	// проверяем сессию
	sess := sessionManager.Check(
		&SessionID{
			ID: sessId.ID,
		})
	t.Log("sess", sess)

	// удаляем сессию
	sessionManager.Delete(
		&SessionID{
			ID: sessId.ID,
		})

	// проверяем еще раз
	sess = sessionManager.Check(
		&SessionID{
			ID: sessId.ID,
		})
	t.Log("sess", sess)
	t.Fail()
}
