package session

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
)

type Session struct {
	ID     string
	UserID uint32
}

func NewSession(userID uint32) *Session {
	// лучше генерировать из заданного алфавита, но так писать меньше и для учебного примера ОК
	randID := make([]byte, 16)
	rand.Read(randID)

	return &Session{
		ID:     fmt.Sprintf("%x", randID),
		UserID: userID,
	}
}

var (
	ErrNoAuth = errors.New("No session found")
)

type sessKey string

var SessionKey sessKey = "sessionKey"

func SessionFromContext(ctx context.Context) (*Session, error) {
	sess, ok := ctx.Value(SessionKey).(*Session)
	if !ok || sess == nil {
		return nil, ErrNoAuth
	}
	return sess, nil
}
