package main

import (
	"testing"
)

func TestExample(t *testing.T) {
	// создаем сессию
	sessId, err := AuthCreateSession(
		&Session{
			Login:     "anton",
			UserAgent: "safari",
		})
	t.Log("sessId", sessId, err)

	// проверяем сессию
	sess := AuthCheckSession(
		&SessionID{
			ID: sessId.ID,
		})
	t.Log("sess", sess)

	// удаляем сессию
	AuthDeleteSession(
		&SessionID{
			ID: sessId.ID,
		})

	// проверяем еще раз
	sess = AuthCheckSession(
		&SessionID{
			ID: sessId.ID,
		})
	t.Log("sess", sess)
	t.Fail()
}
