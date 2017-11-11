package carrot

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"sync"
	"time"
)

const (
	defaultSessionClosedTimeoutDuration = 10 // seconds
)

var (
	oneSessionStore      sync.Once
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
	lengthMutex  *sync.RWMutex
}

func NewDefaultSessionManager() SessionStore {
	oneSessionStore.Do(func() {
		sessionStoreInstance = &DefaultSessionStore{
			sessionStore: &sync.Map{},
			length:       0,
			lengthMutex:  &sync.RWMutex{},
		}
	})

	return sessionStoreInstance
}

func (s *DefaultSessionStore) NewSession() (SessionToken, *Session, error) {

	uuid, err := generateUUID()
	if err != nil {
		return "", nil, err
	}

	token := SessionToken(uuid)

	// set to expiryTime not time.Now
	ctx := Session{
		Token:      token,
		expireTime: refreshExpiryTime(),
	}

	s.sessionStore.Store(token, &ctx)
	s.lengthMutex.Lock()
	s.length += 1
	s.lengthMutex.Unlock()

	log.WithFields(log.Fields{
		"session_token": token,
		"session_count": s.length,
	})

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
	s.lengthMutex.Lock()
	s.length -= 1
	s.lengthMutex.Unlock()
	return nil
}

func (s *DefaultSessionStore) Range(f func(key, value interface{}) bool) {
	s.sessionStore.Range(f)
}

func (s *DefaultSessionStore) Length() int {
	s.lengthMutex.RLock()
	length := s.length
	s.lengthMutex.RUnlock()
	return length
}
