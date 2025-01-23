package gtsclient

import (
	"fmt"
	"net/http"
	"time"

	"codeflow.dananglin.me.uk/apollo/enbas/internal/config"
	"codeflow.dananglin.me.uk/apollo/enbas/internal/info"
	"codeflow.dananglin.me.uk/apollo/enbas/internal/utilities"
)

const (
	applicationJSON   string = "application/json; charset=utf-8"
	redirectURI       string = "urn:ietf:wg:oauth:2.0:oob"
	authCodeURLFormat string = "%s/oauth/authorize?client_id=%s&redirect_uri=%s&response_type=code"
)

type (
	NoRPCArgs    struct{}
	NoRPCResults struct{}

	GTSClient struct {
		authentication config.Credentials
		httpClient     http.Client
		timeout        time.Duration
		mediaTimeout   time.Duration
		userAgent      string
	}
)

// NewGTSClient creates GTSClient value for connecting with the GoToSocial instance. If the credentials
// file is present then the authentication details is retrieved for the current account in use.
// If the file is not present then a zero-valued authentication value is used which must be updated later.
func NewGTSClient(cfg *config.Config) (*GTSClient, error) {
	var auth config.Credentials

	exists, err := utilities.FileExists(cfg.CredentialsFile)
	if err != nil {
		return nil, fmt.Errorf("error checking for the credentials file: %w", err)
	}

	if exists {
		auth, err = authFromFile(cfg.CredentialsFile)
		if err != nil {
			return nil, fmt.Errorf("error getting the authentication details from the credentials file: %w", err)
		}
	} else {
		auth = config.Credentials{
			Instance:     "",
			ClientID:     "",
			ClientSecret: "",
			AccessToken:  "",
		}
	}

	gtsClient := GTSClient{
		authentication: auth,
		httpClient:     http.Client{},
		timeout:        time.Duration(cfg.GTSClient.Timeout) * time.Second,
		mediaTimeout:   time.Duration(cfg.GTSClient.MediaTimeout) * time.Second,
		userAgent:      info.ApplicationTitledName + "/" + info.BinaryVersion,
	}

	return &gtsClient, nil
}

func authFromFile(path string) (config.Credentials, error) {
	creds, err := config.NewCredentialsConfigFromFile(path)
	if err != nil {
		return config.Credentials{}, fmt.Errorf("error getting the credentials from the credentials file: %w", err)
	}

	auth, ok := creds.Credentials[creds.CurrentAccount]
	if !ok {
		return config.Credentials{}, Error{"the authentication details seems to be missing for the current account (" + creds.CurrentAccount + ")"}
	}

	return auth, nil
}

// UpdateAuthentication updates the authentication details for the GTSClient.
func (g *GTSClient) UpdateAuthentication(auth config.Credentials, _ *NoRPCResults) error {
	g.authentication = auth

	return nil
}
