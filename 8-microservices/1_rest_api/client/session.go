package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
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

const (
	httpRequestTimeout = 100 * time.Millisecond
	baseURL            = "http://localhost:8081/session"
)

type SessionManager struct {
	client *http.Client
}

func NewSessionManager() *SessionManager {
	client := &http.Client{
		Timeout: httpRequestTimeout,
	}

	return &SessionManager{
		client: client,
	}
}

func (sm *SessionManager) Create(in *Session) (*SessionID, error) {
	data := url.Values{}
	data.Add("login", in.Login)
	data.Add("ua", in.UserAgent)

	resp, err := sm.client.PostForm(baseURL, data)
	if err != nil {
		fmt.Println("SessionManager/create request error:", err)
		return nil, err
	}

	respData, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("SessionManager/create: failed to read response body: %v\n", err)
		return nil, err
	}

	if resp.StatusCode == http.StatusBadRequest {
		fmt.Printf("SessionManager/create: invalid request, resp=%s\n", respData)
		return nil, nil
	}

	id := new(SessionID)
	if err = json.Unmarshal(respData, id); err != nil {
		fmt.Printf("SessionManager/create: failed to unmarshal response: %v, resp=%s\n", err, respData)
		return nil, err
	}

	return id, nil
}

func (sm *SessionManager) Check(in *SessionID) *Session {
	resp, err := sm.client.Get(fmt.Sprintf("%s/%s", baseURL, in.ID))
	if err != nil {
		fmt.Println("SessionManager/check request error:", err)
		return nil
	}

	respData, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("SessionManager/check: failed to read response body:", err)
		return nil
	}

	session := new(Session)
	if err = json.Unmarshal(respData, session); err != nil {
		fmt.Printf("SessionManager/check: failed to unmarshal response: %v, resp=%s", err, respData)
		return nil
	}

	if session.Login == "" {
		return nil
	}

	return session
}

func (sm *SessionManager) Delete(in *SessionID) {
	parsedURL, err := url.Parse(fmt.Sprintf("%s/%s", baseURL, in.ID))
	if err != nil {
		fmt.Printf("SessionManager/delete: failed to parse request url: %v, url=%s", err, parsedURL)
		return
	}

	req := &http.Request{
		Method: http.MethodDelete,
		URL:    parsedURL,
	}

	if _, err := sm.client.Do(req); err != nil {
		fmt.Println("SessionManager/delete request error:", err)
	}
}
