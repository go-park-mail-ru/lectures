package main

import (
	"testing"
)

func TestExample(t *testing.T) {
	sessManager = NewSessionManager()

	// создаем сессию
	sessId, err := sessManager.Create(
		&Session{
			Login:     "anton",
			UserAgent: "safari",
		})
	t.Log("sessId", sessId, err)

	// проверяем сессию
	sess := sessManager.Check(
		&SessionID{
			ID: sessId.ID,
		})
	t.Log("sess", sess)

	// удаляем сессию
	sessManager.Delete(
		&SessionID{
			ID: sessId.ID,
		})

	// проверяем еще раз
	sess = sessManager.Check(
		&SessionID{
			ID: sessId.ID,
		})
	t.Log("sess", sess)
	t.Fail()
}
