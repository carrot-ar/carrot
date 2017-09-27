package buddy

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"sync"
	"time"
)

const (
	nilSessionToken                     = ""
	defaultSessionClosedTimeoutDuration = 10 // seconds
)

// Potentially will need to be a sync Map
type Session struct {
	Token      SessionToken
	Client     *Client
	expireTime time.Time

	// bad name, still not sure of the use cases yet
	//itemMap map[string]interface{}
}

func refreshExpiryTime() time.Time {
	return time.Now().Add(time.Second * defaultSessionClosedTimeoutDuration)
}

func (c *Session) sessionDurationExpired() bool {
	if c.expireTime.Before(time.Now()) {
		return true
	}

	return false
}

func (c *Session) SessionExpired() bool {
	return !c.Client.open && c.sessionDurationExpired()
}

type SessionToken string

type SessionStore interface {
	NewSession() (SessionToken, error)
	Exists(SessionToken) bool
	Get(SessionToken) (*Session, error)
	GetByClient(client *Client) (*Session, error)
	SetClient(SessionToken, *Client) error
	Range(func(key, value interface{}) bool)
	Delete(SessionToken) error
	Length() int
}

type DefaultSessionStore struct {
	sessionStore *sync.Map
	length       int
}

func NewDefaultSessionManager() SessionStore {
	return &DefaultSessionStore{
		sessionStore: &sync.Map{},
		length:       0,
	}
}

func (s *DefaultSessionStore) NewSession() (SessionToken, error) {
	// Initialize context fields

	b := make([]byte, 64)
	_, err := rand.Read(b)
	if err != nil {
		return nilSessionToken, err
	}

	token := SessionToken(base64.URLEncoding.EncodeToString(b))

	// set to expiryTime not time.Now
	ctx := Session{
		Token:      token,
		expireTime: refreshExpiryTime(),
	}

	s.sessionStore.Store(token, &ctx)
	s.length += 1

	return token, nil
}

func (s *DefaultSessionStore) Exists(token SessionToken) bool {
	_, ok := s.sessionStore.Load(token)
	if !ok {
		return false
	} else {
		return true
	}
}

// Get does not guarantee that a connection is open for the given context
// this must be checked with the expired() function once the context is retrieved
func (s *DefaultSessionStore) Get(token SessionToken) (*Session, error) {
	ctx, ok := s.sessionStore.Load(token)
	if !ok {
		return nil, fmt.Errorf("session does not exist")
	}

	return ctx.(*Session), nil
}

func (s *DefaultSessionStore) GetByClient(client *Client) (*Session, error) {
	var session *Session

	s.sessionStore.Range(func(key, value interface{}) bool {
		s := value.(*Session)
		if s.Client == client {
			session = s
			return false
		}
		return true
	})

	if session == nil {
		return nil, fmt.Errorf("No session found for client %v", client)
	}

	return session, nil
}

func (s *DefaultSessionStore) SetClient(token SessionToken, client *Client) error {
	ctx, err := s.Get(token)
	if err != nil {
		return err
	}

	ctx.Client = client
	ctx.expireTime = refreshExpiryTime()

	return nil
}

func (s *DefaultSessionStore) Delete(token SessionToken) error {
	s.sessionStore.Delete(token)
	s.length -= 1
	return nil
}

func (s *DefaultSessionStore) Range(f func(key, value interface{}) bool) {
	s.sessionStore.Range(f)
}

func (s *DefaultSessionStore) Length() int {
	return s.length
}
