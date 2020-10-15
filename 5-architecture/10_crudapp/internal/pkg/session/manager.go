package session

import (
	"crudapp/internal/pkg/models"
	"net/http"
	"sync"
	"time"
)

type SessionsManager struct {
	data map[string]*models.Session
	mu   *sync.RWMutex
}

func NewSessionsMem() *SessionsManager {
	return &SessionsManager{
		data: make(map[string]*models.Session, 10),
		mu:   &sync.RWMutex{},
	}
}

func (sm *SessionsManager) Check(r *http.Request) (*models.Session, error) {
	sessionCookie, err := r.Cookie("session_id")
	if err == http.ErrNoCookie {
		return nil, models.ErrNoAuth
	}

	sm.mu.RLock()
	sess, ok := sm.data[sessionCookie.Value]
	sm.mu.RUnlock()

	if !ok {
		return nil, models.ErrNoAuth
	}

	return sess, nil
}

func (sm *SessionsManager) Create(w http.ResponseWriter, userID uint32) (*models.Session, error) {
	sess := models.NewSession(userID)

	sm.mu.Lock()
	sm.data[sess.ID] = sess
	sm.mu.Unlock()

	cookie := &http.Cookie{
		Name:    "session_id",
		Value:   sess.ID,
		Expires: time.Now().Add(90 * 24 * time.Hour),
		Path:    "/",
	}
	http.SetCookie(w, cookie)
	return sess, nil
}

func (sm *SessionsManager) DestroyCurrent(w http.ResponseWriter, r *http.Request) error {
	sess, err := models.SessionFromContext(r.Context())
	if err != nil {
		return err
	}

	sm.mu.Lock()
	delete(sm.data, sess.ID)
	sm.mu.Unlock()

	cookie := http.Cookie{
		Name:    "session_id",
		Expires: time.Now().AddDate(0, 0, -1),
		Path:    "/",
	}
	http.SetCookie(w, &cookie)
	return nil
}
