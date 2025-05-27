package server

import (
	"crypto/rand"
	"fmt"
	"net/rpc"
	"os"
	"os/exec"
	"time"

	"codeflow.dananglin.me.uk/apollo/enbas/internal/config"
	"codeflow.dananglin.me.uk/apollo/enbas/internal/sessions"
	"codeflow.dananglin.me.uk/apollo/enbas/internal/utilities"
)

type Session struct {
	client    *rpc.Client
	sessionID string
}

func (s *Session) Client() *rpc.Client {
	return s.client
}

// StartSession starts a new client session with the server.
// It will attempt to make a connection to the server, generate a new session ID and register
// said session ID to the server. Upon successful completion the caller receives the Session
// which it can use to end the session once it has completed it's operation.
// If the server is not running then an attempt is made to run and connect to a temporary server.
func StartSession(cfg config.Server, configPath string) (*Session, error) {
	socketPath, err := utilities.AbsolutePath(cfg.SocketPath)
	if err != nil {
		return nil, fmt.Errorf(
			"unable to calculate the absolute path of the socket file: %w",
			err,
		)
	}

	exists, err := utilities.FileExists(socketPath)
	if err != nil {
		return nil, fmt.Errorf(
			"received an unexpected error after checking for the socket file: %w",
			err,
		)
	}

	// Attempt the server connection if the socket file is present.
	if exists {
		session := Session{
			client:    nil,
			sessionID: "",
		}

		if err := startSession(&session, socketPath); err != nil {
			return nil, fmt.Errorf("error starting a new session with the existing server: %w", err)
		}

		return &session, nil
	}

	// The socket file is not present so we will
	// attempt to start a new server process.
	server := exec.Command(os.Args[0], "--config", configPath, "start", "server")

	if err := server.Start(); err != nil {
		return nil, fmt.Errorf("error starting the server: %w", err)
	}

	// Attempt to start a new session with the server.
	session := Session{
		client:    nil,
		sessionID: "",
	}

	for range 3 {
		time.Sleep(100 * time.Millisecond)

		err = startSession(&session, socketPath)
		if err == nil {
			break
		}
	}

	if err != nil {
		return nil, fmt.Errorf("error connecting to the server: %w", err)
	}

	return &session, nil
}

// EndSession ends an existing session with the server.
// It will attempt to remove the session ID from the session store
// and close the connection to the server.
func EndSession(session *Session) error {
	defer session.client.Close()

	var initialised bool
	if err := session.client.Call(
		"SessionStore.IsInitialised",
		sessions.NoRPCArgs{},
		&initialised,
	); err != nil {
		return fmt.Errorf("error checking if the session store is initialised: %w", err)
	}

	if !initialised {
		return nil
	}

	if err := session.client.Call(
		"SessionStore.Remove",
		session.sessionID,
		nil,
	); err != nil {
		return fmt.Errorf("error removing the session ID from the session store: %w", err)
	}

	return nil
}

// startSession adds a new client and session ID to a Session value. If a client is
// already created then it only creates a new session ID if it has not been created
// already.
func startSession(session *Session, socketPath string) error {
	if session.client == nil {
		client, err := rpc.Dial("unix", socketPath)
		if err != nil {
			return fmt.Errorf(
				"error connecting to the RPC server: %w",
				err,
			)
		}

		session.client = client
	}

	if session.sessionID == "" {
		sessionID, err := newSessionID(session.client)
		if err != nil {
			return fmt.Errorf(
				"error creating the session ID: %w",
				err,
			)
		}

		session.sessionID = sessionID
	}

	return nil
}

// newSessionID attempts to create a new session ID and add it to the session store.
func newSessionID(client *rpc.Client) (string, error) {
	var initialised bool
	if err := client.Call(
		"SessionStore.IsInitialised",
		sessions.NoRPCArgs{},
		&initialised,
	); err != nil {
		return "", fmt.Errorf("error checking if the session store is initialised: %w", err)
	}

	if !initialised {
		return "", nil
	}

	sessionID := make([]byte, 16)

	if _, err := rand.Read(sessionID); err != nil {
		return "", fmt.Errorf("error creating the session ID: %w", err)
	}

	if err := client.Call(
		"SessionStore.Add",
		string(sessionID),
		nil,
	); err != nil {
		return "", fmt.Errorf("error adding the session ID to the session store: %w", err)
	}

	return string(sessionID), nil
}
