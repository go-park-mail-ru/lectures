package main

import (
	"context"
	"fmt"

	jsonrpc "github.com/ybbus/jsonrpc/v3"
)

type SessionManagerI interface {
	Create(context.Context, *Session) (*SessionID, error)
	Check(context.Context, *SessionID) *Session
	Delete(context.Context, *SessionID)
}

type Session struct {
	Login     string
	UserAgent string
}

type SessionID struct {
	ID string
}

type SessionManager struct {
	client jsonrpc.RPCClient
}

func NewSessionManager() *SessionManager {
	client := jsonrpc.NewClient("http://localhost:8081/rpc")

	return &SessionManager{
		client: client,
	}
}

func (sm *SessionManager) Create(ctx context.Context, in *Session) (*SessionID, error) {
	id := new(SessionID)

	// Из-за особенностей библиотеки нужно завернуть аргументы в слайс:
	if err := sm.client.CallFor(ctx, id, "SessionManager.Create", []any{in}); err != nil {
		fmt.Println("SessionManager.Create error:", err)
		return nil, nil
	}

	return id, nil
}

func (sm *SessionManager) Check(ctx context.Context, in *SessionID) *Session {
	session := new(Session)

	if err := sm.client.CallFor(ctx, session, "SessionManager.Check", []any{in}); err != nil {
		fmt.Println("SessionManager.Check error:", err)
		return nil
	}

	return session
}

func (sm *SessionManager) Delete(ctx context.Context, in *SessionID) {
	var reply int

	if err := sm.client.CallFor(ctx, &reply, "SessionManager.Delete", []any{in}); err != nil {
		fmt.Println("SessionManager.Delete error:", err)
	}
}
