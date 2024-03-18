package main

import (
	"encoding/json"
	"fmt"
	"log"

	tarantool "github.com/tarantool/go-tarantool"
)

type Session struct {
	Login     string
	Useragent string
}

type SessionID struct {
	ID string
}

type SessionManager struct {
	tConn *tarantool.Connection
}

func NewSessionManager(conn *tarantool.Connection) *SessionManager {
	return &SessionManager{
		tConn: conn,
	}
}

func (sm *SessionManager) Create(in *Session) (*SessionID, error) {
	dataSerialized, err := json.Marshal(in)
	if err != nil {
		return nil, fmt.Errorf("data serialization failed: %s", err)
	}

	dataStr := string(dataSerialized)
	log.Printf("try to save data: %s", dataStr)

	resp, err := sm.tConn.Eval("return new_session(...)", []interface{}{dataStr})
	if err != nil {
		return nil, fmt.Errorf("error while calling function: %s", err)
	}

	data := resp.Data[0]
	if id, ok := data.(string); ok {
		return &SessionID{id}, nil
	}

	return nil, fmt.Errorf("cannot cast into int")
}

func (sm *SessionManager) Check(in *SessionID) (*Session, error) {
	resp, err := sm.tConn.Call("check_session", []interface{}{in.ID})
	if err != nil {
		fmt.Println("cannot check session", err)
		return nil, err
	}

	data := resp.Data[0]
	if data == nil {
		return &Session{}, nil
	}
	sessionDataSlice, ok := data.([]interface{})
	if !ok {
		return nil, fmt.Errorf("cannot cast data: %v", sessionDataSlice)
	}

	if sessionDataSlice[1] == nil {
		return nil, nil
	}

	sessionData, ok := sessionDataSlice[1].(string)
	if !ok {
		return nil, fmt.Errorf("cannot cast data: %v", sessionDataSlice[1])
	}

	sess := &Session{}
	err = json.Unmarshal([]byte(sessionData), sess)
	if err != nil {
		log.Printf("cant unpack session data(%s): %v\n", sessionData, err)
		return nil, nil
	}
	return sess, nil
}

func (sm *SessionManager) Delete(in *SessionID) {
	// mkey := "sessions:" + in.ID
	// _, err := redis.Int(sm.redisConn.Do("DEL", mkey))
	// if err != nil {
	// 	log.Println("redis error:", err)
	// }
}
