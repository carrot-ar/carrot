package buddy

import (
	"crypto/rand"
	"fmt"
	"sync"
	"encoding/base64"
)

// Potentially will need to be a sync Map
type Context map[string]interface{}

type SessionToken string

type SessionStore interface {
	Get(SessionToken) (Context, error)
	NewSession() (SessionToken, error)
	Delete(SessionToken) error
	Exists(SessionToken) bool
}

type DefaultSessionStore struct {
	sessionStore *sync.Map
}

func (s *DefaultSessionStore) Get(token SessionToken) (Context, error) {
	ctx, ok := s.sessionStore.Load(token)
	if !ok {
		return nil, fmt.Errorf("session does not exist")
	}

	return ctx.(Context), nil
}

func (s *DefaultSessionStore) Delete(token SessionToken) error {
	s.sessionStore.Delete(token)
	return nil
}

func (s *DefaultSessionStore) Exists(token SessionToken) bool {
	_, ok := s.sessionStore.Load(token)
	if !ok {
		return false
	} else {
		return true
	}
}

func (s *DefaultSessionStore) NewSession() (SessionToken, error) {
	// Initialize context fields
	ctx := Context{}
	b := make([]byte, 64)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	token := SessionToken(base64.URLEncoding.EncodeToString(b))

	s.sessionStore.Store(token, ctx)

	return SessionToken(token), nil
}

func NewDefaultSessionManager() SessionStore {
	return &DefaultSessionStore{
		sessionStore: &sync.Map{},
	}
}
