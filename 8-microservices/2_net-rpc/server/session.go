package main

import (
	"fmt"
	"math/rand"
	"sync"
)

type Session struct {
	Login     string
	UserAgent string
}

type SessionID struct {
	ID string
}

const sessionKeyLen = 10

type SessionManager struct {
	sessions map[SessionID]*Session
	mu       sync.RWMutex
}

func NewSessionManager() *SessionManager {
	return &SessionManager{
		mu:       sync.RWMutex{},
		sessions: map[SessionID]*Session{},
	}
}

func (sm *SessionManager) Create(in *Session, out *SessionID) error {
	fmt.Printf("call Create, input=%#v\n", in)

	id := &SessionID{generateRandomString(sessionKeyLen)}
	sm.mu.Lock()
	sm.sessions[*id] = in
	sm.mu.Unlock()
	*out = *id

	return nil
}

func (sm *SessionManager) Check(in *SessionID, out *Session) error {
	fmt.Printf("call Check, input=%#v\n", in)

	sm.mu.RLock()
	defer sm.mu.RUnlock()
	if sess, ok := sm.sessions[*in]; ok {
		*out = *sess
	}

	return nil
}

func (sm *SessionManager) Delete(in *SessionID, out *int) error {
	fmt.Printf("call Delete, input=%#v\n", in)

	sm.mu.Lock()
	defer sm.mu.Unlock()
	delete(sm.sessions, *in)
	*out = 1

	return nil
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func generateRandomString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
