package main

import (
	"fmt"
	"log"
	"net/rpc"
)

type Session struct {
	Login     string
	Useragent string
}

type SessionID struct {
	ID string
}

type SessionManagerI interface {
	Create(*Session) (*SessionID, error)
	Check(*SessionID) *Session
	Delete(*SessionID)
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
	fmt.Println("call Create", in)
	err := sm.client.Call("SessionManager.Create", in, id)
	if err != nil {
		fmt.Println("SessionManager.Create error:", err)
		return nil, nil
	}
	return id, nil
}

func (sm *SessionManager) Check(in *SessionID) *Session {
	sess := new(Session)
	fmt.Println("call Check", in)
	err := sm.client.Call("SessionManager.Check", in, sess)
	if err != nil {
		fmt.Println("SessionManager.Check error:", err)
		return nil
	}
	return sess
}

func (sm *SessionManager) Delete(in *SessionID) {
	var reply int
	fmt.Println("call Delete", in)
	err := sm.client.Call("SessionManager.Delete", in, &reply)
	if err != nil {
		fmt.Println("SessionManager.Delete error:", err)
	}
}
