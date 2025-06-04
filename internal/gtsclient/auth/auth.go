package auth

import (
	"fmt"
	"sync"

	"codeflow.dananglin.me.uk/apollo/enbas/internal/config"
)

type Auth struct {
	mu               sync.RWMutex
	currentAccountID string
	instanceURL      string
	token            string
}

func NewAuthFromConfig(path string) (*Auth, error) {
	creds, err := config.NewCredentialsConfigFromFile(path)
	if err != nil {
		return nil, fmt.Errorf(
			"error getting the credentials from the credentials file: %w",
			err,
		)
	}

	// If the current account is not set, return a config with zero-valued fields.
	if creds.CurrentAccount == "" {
		return NewAuthZero(), nil
	}

	authCfg, ok := creds.Credentials[creds.CurrentAccount]
	if !ok {
		return nil, NewNoAuthenticationDetailsError(creds.CurrentAccount)
	}

	return &Auth{
		mu:               sync.RWMutex{},
		currentAccountID: "",
		instanceURL:      authCfg.Instance,
		token:            authCfg.AccessToken,
	}, nil
}

func NewAuthZero() *Auth {
	return &Auth{
		mu:               sync.RWMutex{},
		currentAccountID: "",
		instanceURL:      "",
		token:            "",
	}
}

func (a *Auth) UpdateAuth(authCfg config.Credentials) {
	a.mu.Lock()
	defer a.mu.Unlock()

	a.instanceURL = authCfg.Instance
	a.token = authCfg.AccessToken
	a.currentAccountID = ""
}

func (a *Auth) GetInstanceURL() string {
	a.mu.RLock()
	defer a.mu.RUnlock()

	return a.instanceURL
}

func (a *Auth) GetToken() string {
	a.mu.RLock()
	defer a.mu.RUnlock()

	return a.token
}

func (a *Auth) UpdateCurrentAccountID(accountID string) {
	a.mu.Lock()
	defer a.mu.Unlock()

	if a.currentAccountID == accountID {
		return
	}

	a.currentAccountID = accountID
}

func (a *Auth) GetCurrentAccountID() string {
	a.mu.RLock()
	defer a.mu.RUnlock()

	return a.currentAccountID
}
