package main

import (
	"context"
	"fmt"
	"math/rand"
	"sync"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/go-park-mail-ru/lectures/8-microservices/5_grpc/session"
)

const sessKeyLen = 10

type SessionManager struct {
	session.UnimplementedAuthCheckerServer

	sessions map[string]*session.Session
	mu       sync.RWMutex
}

func NewSessionManager() *SessionManager {
	return &SessionManager{
		sessions: map[string]*session.Session{},
		mu:       sync.RWMutex{},
	}
}

func (sm *SessionManager) Create(ctx context.Context, in *session.Session) (*session.SessionID, error) {
	fmt.Printf("call Create, login=%s, ua=%s\n", in.Login, in.Useragent)

	id := &session.SessionID{ID: generateRandomString(sessKeyLen)}
	sm.mu.Lock()
	sm.sessions[id.ID] = in
	sm.mu.Unlock()

	return id, nil
}

func (sm *SessionManager) Check(ctx context.Context, in *session.SessionID) (*session.Session, error) {
	fmt.Printf("call Check, id=%s\n", in.ID)

	sm.mu.RLock()
	defer sm.mu.RUnlock()
	if sess, ok := sm.sessions[in.ID]; ok {
		return sess, nil
	}

	return nil, status.Errorf(codes.NotFound, "session not found")
}

func (sm *SessionManager) Delete(ctx context.Context, in *session.SessionID) (*session.Nothing, error) {
	fmt.Printf("call Delete, id=%s\n", in.ID)

	sm.mu.Lock()
	defer sm.mu.Unlock()
	delete(sm.sessions, in.ID)

	return &session.Nothing{Dummy: true}, nil
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func generateRandomString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
