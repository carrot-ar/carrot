package carrot

import (
	"errors"
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
	Token         SessionToken
	expireTime    time.Time
	mutex         *sync.Mutex
	primaryDevice bool
	T_L           *offset
	T_P           *offset
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
	GetPrimaryDeviceToken() (SessionToken, error)
	Range(func(key, value interface{}) bool)
	Delete(SessionToken) error
	Length() int
}

type DefaultSessionStore struct {
	sessionStore *sync.Map
	length       int
	lengthMutex  *sync.RWMutex
}

// Provides a pointer to a singleton of the SessionStore interface
// for use in Carrot modules
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

// Confirms whether a session is the primary device or not
func (s *Session) isPrimaryDevice() bool {
	return s.primaryDevice
}

// Creates a new session and adds it to the SessionStore. Returns the generated, UUID,
// a pointer to the Session object, and an error
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

// Determines whether a given session exists based on the SessionToken provided
func (s *DefaultSessionStore) Exists(token SessionToken) bool {
	_, ok := s.sessionStore.Load(token)
	if !ok {
		return false
	} else {
		return true
	}
}

// Retrieves a session based on its SessionToken
//
// Get() does not guarantee that a connection is open this must be checked with the expired()
// function once the session is retrieved
func (s *DefaultSessionStore) Get(token SessionToken) (*Session, error) {
	ctx, ok := s.sessionStore.Load(token)

	if !ok {
		return nil, fmt.Errorf("session: session does not exist")
	}

	return ctx.(*Session), nil
}

// Retrieves the SessionToken of the primary device and returns an error if no
// primary device exists
func (s *DefaultSessionStore) GetPrimaryDeviceToken() (SessionToken, error) {
	var token SessionToken
	var err error
	s.sessionStore.Range(func(t, session interface{}) bool {
		s := session.(*Session)
		if s.isPrimaryDevice() {
			token = t.(SessionToken)
		}
		return true
	})
	if token == "" {
		err = errors.New("Could not locate primary device in sessions list")
	}
	return token, err
}

// Deletes a session and updates the session count variable
func (s *DefaultSessionStore) Delete(token SessionToken) error {
	s.sessionStore.Delete(token)
	s.lengthMutex.Lock()
	s.length -= 1
	s.lengthMutex.Unlock()
	return nil
}

// Provided a function, Range() iterates over every session in a thread-safe manner
// and applies that function to the session. If the function provided as an argument
// returns true after acting on the session, iteration continues. Otherwise, the range()
// function breaks
func (s *DefaultSessionStore) Range(f func(key, value interface{}) bool) {
	s.sessionStore.Range(f)
}

// Returns the number of sessions currently existing including inactive sessions
func (s *DefaultSessionStore) Length() int {
	s.lengthMutex.RLock()
	length := s.length
	s.lengthMutex.RUnlock()
	return length
}
