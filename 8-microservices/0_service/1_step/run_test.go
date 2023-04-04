package main

import (
	"testing"
)

func TestExample(t *testing.T) {
	// создаем сессию
	sessId, err := AuthCreateSession(
		&Session{
			Login:     "anton",
			Useragent: "safari",
		})
	t.Log("sessId", sessId, err)

	// проеряем сессию
	sess := AuthCheckSession(
		&SessionID{
			ID: sessId.ID,
		})
	t.Log("sess", sess)

	// удаляем сессию
	AuthSessionDelete(
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
