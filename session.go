package carrot

import (
	"crypto/rand"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"sync"
	"time"
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
	return time.Now().Add(time.Second * config.Session.DefaultSessionClosedTimeoutDuration)
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
	// Initialize context fields

	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return config.Session.NilSessionToken, nil, err
	}

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

// generate UUID fulfilling RFC 4122
func generateUUID() (string, error) {
	uuid := make([]byte, 16)
	n, err := io.ReadFull(rand.Reader, uuid)
	if n != len(uuid) || err != nil {
		return "", err
	}

	// variant bits
	uuid[8] = uuid[8]&^0xc0 | 0x80

	// version 4
	uuid[6] = uuid[6]&^0xf0 | 0x40
	return fmt.Sprintf("%x-%x-%x-%x-%x", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:]), nil
}
