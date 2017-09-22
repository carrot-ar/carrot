package buddy

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"sync"
)

// Potentially will need to be a sync Map
type Context map[string]interface{}

type SessionToken string

type SessionStore interface {
	Get(SessionToken) (Context, error)
	NewSession() (SessionToken, error)
	Delete(SessionToken) error
	Exists(SessionToken) bool
	Length() int
}

type DefaultSessionStore struct {
	sessionStore *sync.Map
	length       int
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
	s.length -= 1
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
	s.length += 1

	return SessionToken(token), nil
}

func (s *DefaultSessionStore) Length() int {
	return s.length
}

func NewDefaultSessionManager() SessionStore {
	return &DefaultSessionStore{
		sessionStore: &sync.Map{},
		length:       0,
	}
}
