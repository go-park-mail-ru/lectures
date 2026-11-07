package main

import (
	"fmt"
	"log"
	"net/rpc"
)

type SessionManagerI interface {
	Create(*Session) (*SessionID, error)
	Check(*SessionID) *Session
	Delete(*SessionID)
}

type Session struct {
	Login     string
	UserAgent string
}

type SessionID struct {
	ID string
}

type SessionManager struct {
	client *rpc.Client
}

func NewSessionManager() *SessionManager {
	client, err := rpc.DialHTTP("tcp", "localhost:8081")
	if err != nil {
		log.Fatal("dialing:", err)
	}

	return &SessionManager{
		client: client,
	}
}

func (sm *SessionManager) Create(in *Session) (*SessionID, error) {
	id := new(SessionID)

	if err := sm.client.Call("SessionManager.Create", in, id); err != nil {
		fmt.Println("SessionManager.Create error:", err)
		return nil, nil
	}

	return id, nil
}

func (sm *SessionManager) Check(in *SessionID) *Session {
	sess := new(Session)

	if err := sm.client.Call("SessionManager.Check", in, sess); err != nil {
		fmt.Println("SessionManager.Check error:", err)
		return nil
	}

	return sess
}

func (sm *SessionManager) Delete(in *SessionID) {
	var reply int

	if err := sm.client.Call("SessionManager.Delete", in, &reply); err != nil {
		fmt.Println("SessionManager.Delete error:", err)
	}
}
