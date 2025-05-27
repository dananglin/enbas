package sessions

type SessionIDExistsError struct {
	sessionID string
}

func NewSessionIDExistsError(sessionID string) SessionIDExistsError {
	return SessionIDExistsError{sessionID}
}

func (e SessionIDExistsError) Error() string {
	return "the session ID (" +
		e.sessionID +
		") is already present in the SessionStore"
}

type NonInitialisedSessionStoreError struct{}

func (e NonInitialisedSessionStoreError) Error() string {
	return "the SessionStore is not initialised"
}
