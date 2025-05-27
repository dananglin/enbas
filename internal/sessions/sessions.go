package sessions

import (
	"sync"
)

type (
	NoRPCArgs    struct{}
	NoRPCResults struct{}

	SessionStore struct {
		mu          *sync.Mutex
		sessionIDs  map[string]struct{}
		initialised bool
	}
)

// NewSessionStore returns an initialised or a non-initialised SessionStore.
func NewSessionStore(initialise bool) *SessionStore {
	if !initialise {
		// return an initialised SessionStore typically
		// for a server process launched without an idle
		// timeout.
		return &SessionStore{
			mu:          nil,
			sessionIDs:  nil,
			initialised: false,
		}
	}

	// return an initialised SessionStore for a server
	// launched with an idle timeout.
	return &SessionStore{
		mu:          &sync.Mutex{},
		sessionIDs:  make(map[string]struct{}),
		initialised: true,
	}
}

// IsInitialised returns true if the SessionStore was initialised.
func (s *SessionStore) IsInitialised(_ NoRPCArgs, isInitialised *bool) error {
	*isInitialised = s.initialised

	return nil
}

// Add adds a session ID to the SessionStore at the request of the client session.
func (s *SessionStore) Add(sessionID string, _ *NoRPCResults) error {
	if !s.initialised {
		return NonInitialisedSessionStoreError{}
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.sessionIDs[sessionID]; exists {
		return NewSessionIDExistsError(sessionID)
	}

	s.sessionIDs[sessionID] = struct{}{}

	return nil
}

// Remove removes a session ID to the SessionStore at the request of the client session.
func (s *SessionStore) Remove(sessionID string, _ *NoRPCResults) error {
	if !s.initialised {
		return NonInitialisedSessionStoreError{}
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.sessionIDs, sessionID)

	return nil
}

// NumSessionIDs returns the number of session IDs stored in the SessionStore.
func (s *SessionStore) NumSessionIDs() int {
	return len(s.sessionIDs)
}

// IsEmpty returns true if the SessionStore has no session IDs.
func (s *SessionStore) IsEmpty() bool {
	return len(s.sessionIDs) == 0
}
