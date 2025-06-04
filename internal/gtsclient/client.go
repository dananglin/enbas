package gtsclient

import (
	"fmt"
	"net/http"
	"time"

	"codeflow.dananglin.me.uk/apollo/enbas/internal/config"
	"codeflow.dananglin.me.uk/apollo/enbas/internal/gtsclient/auth"
	"codeflow.dananglin.me.uk/apollo/enbas/internal/info"
	"codeflow.dananglin.me.uk/apollo/enbas/internal/utilities"
)

const applicationJSON string = "application/json; charset=utf-8"

type (
	NoRPCArgs    struct{}
	NoRPCResults struct{}

	GTSClient struct {
		auth         *auth.Auth
		httpClient   http.Client
		timeout      time.Duration
		mediaTimeout time.Duration
		userAgent    string
	}
)

// NewGTSClient creates GTSClient value for connecting with the GoToSocial instance. If the credentials
// file is present then the authentication details is retrieved for the current account in use.
// If the file is not present then a zero-valued authentication value is used which must be updated later.
func NewGTSClient(cfg config.Config) (*GTSClient, error) {
	var newAuth *auth.Auth

	exists, err := utilities.FileExists(cfg.CredentialsFile)
	if err != nil {
		return nil, fmt.Errorf("error checking for the credentials file: %w", err)
	}

	if exists {
		newAuth, err = auth.NewAuthFromConfig(cfg.CredentialsFile)
		if err != nil {
			return nil, fmt.Errorf(
				"error getting the authentication details from the credentials file: %w",
				err,
			)
		}
	} else {
		newAuth = auth.NewAuthZero()
	}

	gtsClient := GTSClient{
		auth:         newAuth,
		httpClient:   http.Client{},
		timeout:      time.Duration(cfg.GTSClient.Timeout) * time.Second,
		mediaTimeout: time.Duration(cfg.GTSClient.MediaTimeout) * time.Second,
		userAgent:    info.ApplicationTitledName + "/" + info.BinaryVersion,
	}

	return &gtsClient, nil
}

// UpdateAuthentication updates the authentication details for the GTSClient.
func (g *GTSClient) UpdateAuthentication(authCfg config.Credentials, _ *NoRPCResults) error {
	g.auth.UpdateAuth(authCfg)

	return nil
}
