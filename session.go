package carrot

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"sync"
	"time"
)

const (
	nilSessionToken                     = ""
	defaultSessionClosedTimeoutDuration = 10 // seconds
)

var (
	oneSessionStore sync.Once

	sessionStoreInstance SessionStore
)

// Potentially will need to be a sync Map
type Session struct {
	Token      SessionToken
	expireTime time.Time
	mutex      *sync.Mutex
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

type SessionToken string

type SessionStore interface {
	NewSession() (SessionToken, *Session, error)
	Exists(SessionToken) bool
	Get(SessionToken) (*Session, error)
	Range(func(key, value interface{}) bool)
	Delete(SessionToken) error
	Length() int
}

type DefaultSessionStore struct {
	sessionStore *sync.Map
	length       int
	mutex        *sync.Mutex
}

func NewDefaultSessionManager() SessionStore {
	oneSessionStore.Do(func() {
		sessionStoreInstance = &DefaultSessionStore{
			sessionStore: &sync.Map{},
			length:       0,
			mutex:        &sync.Mutex{},
		}
	})

	return sessionStoreInstance
}

func (s *DefaultSessionStore) NewSession() (SessionToken, *Session, error) {
	// Initialize context fields

	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return nilSessionToken, nil, err
	}

	token := SessionToken(base64.URLEncoding.EncodeToString(b))

	// set to expiryTime not time.Now
	ctx := Session{
		Token:      token,
		expireTime: refreshExpiryTime(),
	}

	s.sessionStore.Store(token, &ctx)
	s.mutex.Lock()
	s.length += 1
	s.mutex.Unlock()

	log.Printf("session: new session created %v, total: %v\n", token, s.length)

	return token, &ctx, nil
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
	//fmt.Printf("session: getting session %v\n", token)

	if !ok {
		return nil, fmt.Errorf("session: session does not exist")
	}

	return ctx.(*Session), nil
}

func (s *DefaultSessionStore) Delete(token SessionToken) error {
	s.sessionStore.Delete(token)
	s.mutex.Lock()
	s.length -= 1
	s.mutex.Unlock()
	return nil
}

func (s *DefaultSessionStore) Range(f func(key, value interface{}) bool) {
	s.sessionStore.Range(f)
}

func (s *DefaultSessionStore) Length() int {
	s.mutex.Lock()
	length := s.length
	s.mutex.Unlock()
	return length
}
