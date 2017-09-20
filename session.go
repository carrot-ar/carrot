package buddy

import (
	"crypto/sha256"
	"fmt"
	"hash"
	"sync"
)

// Potentially will need to be a sync Map
type Context map[string]interface{}

type SessionToken hash.Hash

type SessionStore interface {
	Get(SessionToken) (Context, error)
	NewSession() SessionToken
	Delete(SessionToken) error
	Exists(SessionToken) bool
}

type DefaultSessionStore struct {
	sessionStore sync.Map
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

func (s *DefaultSessionStore) NewSession() SessionToken {
	// Initialize context fields
	ctx := Context{}
	token := sha256.New()
	s.sessionStore.Store(token, ctx)

	return token
}

func NewDefaultSessionManager() SessionStore {
	return &DefaultSessionStore{}
}
