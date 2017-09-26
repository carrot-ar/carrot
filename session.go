package buddy

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"sync"
)

const (
	nilSessionToken = ""
)

// Potentially will need to be a sync Map
type Context struct {
	Client *Client

	// bad name, still not sure of the use cases yet
	//itemMap map[string]interface{}
}

type SessionToken string

type SessionStore interface {
	Get(SessionToken) (*Context, error)
	SetClient(SessionToken, *Client) error
	NewSession() (SessionToken, error)
	NewSessionClient(client *Client) (SessionToken, error)
	Range(func(key, value interface{}) bool)
	Delete(SessionToken) error
	Exists(SessionToken) bool
	Length() int
}

type DefaultSessionStore struct {
	sessionStore *sync.Map
	length       int
}

func (s *DefaultSessionStore) Get(token SessionToken) (*Context, error) {
	ctx, ok := s.sessionStore.Load(token)
	if !ok {
		return nil, fmt.Errorf("session does not exist")
	}

	return ctx.(*Context), nil
}

func (s *DefaultSessionStore) SetClient(token SessionToken, client *Client) error {
	ctx, err := s.Get(token)
	if err != nil {
		return err
	}

	ctx.Client = client

	return nil
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

func (s *DefaultSessionStore) Range(f func(key, value interface{}) bool) {
	s.sessionStore.Range(f)
}

func (s *DefaultSessionStore) NewSession() (SessionToken, error) {
	// Initialize context fields
	ctx := Context{}
	b := make([]byte, 64)
	_, err := rand.Read(b)
	if err != nil {
		return nilSessionToken, err
	}

	token := SessionToken(base64.URLEncoding.EncodeToString(b))

	s.sessionStore.Store(token, &ctx)
	s.length += 1

	return SessionToken(token), nil
}

func (s *DefaultSessionStore) NewSessionClient(client *Client) (SessionToken, error) {
	token, err := s.NewSession()
	if err != nil {
		return nilSessionToken, err
	}

	c, ok := s.sessionStore.Load(token)
	if ok != true {
		return nilSessionToken, err
	}

	ctx := c.(Context)
	ctx.Client = client

	return token, nil
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
