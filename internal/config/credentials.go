package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"codeflow.dananglin.me.uk/apollo/enbas/internal"
)

const (
	credentialsFileName = "credentials.json"
)

type CredentialsConfig struct {
	CurrentAccount string                 `json:"currentAccount"`
	Credentials    map[string]Credentials `json:"credentials"`
}

type Credentials struct {
	Instance     string `json:"instance"`
	ClientID     string `json:"clientId"`
	ClientSecret string `json:"clientSecret"`
	AccessToken  string `json:"accessToken"`
}

func SaveCredentials(username string, credentials Credentials) (string, error) {
	if err := ensureConfigDir(); err != nil {
		return "", fmt.Errorf("unable to ensure the configuration directory; %w", err)
	}

	var authConfig CredentialsConfig

	filepath := credentialsConfigFile()

	if _, err := os.Stat(filepath); err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return "", fmt.Errorf("unknown error received when running stat on %s; %w", filepath, err)
		}

		authConfig.Credentials = make(map[string]Credentials)
	} else {
		authConfig, err = NewAuthenticationConfigFromFile()
		if err != nil {
			return "", fmt.Errorf("unable to retrieve the existing authentication configuration; %w", err)
		}
	}

	instance := ""

	if strings.HasPrefix(credentials.Instance, "https://") {
		instance = strings.TrimPrefix(credentials.Instance, "https://")
	} else if strings.HasPrefix(credentials.Instance, "http://") {
		instance = strings.TrimPrefix(credentials.Instance, "http://")
	}

	authenticationName := username + "@" + instance

	authConfig.CurrentAccount = authenticationName

	authConfig.Credentials[authenticationName] = credentials

	if err := saveCredentialsConfigFile(authConfig); err != nil {
		return "", fmt.Errorf("unable to save the authentication configuration to file; %w", err)
	}

	return authenticationName, nil
}

func NewAuthenticationConfigFromFile() (CredentialsConfig, error) {
	path := credentialsConfigFile()

	file, err := os.Open(path)
	if err != nil {
		return CredentialsConfig{}, fmt.Errorf("unable to open %s, %w", path, err)
	}
	defer file.Close()

	var authConfig CredentialsConfig

	if err := json.NewDecoder(file).Decode(&authConfig); err != nil {
		return CredentialsConfig{}, fmt.Errorf("unable to decode the JSON data; %w", err)
	}

	return authConfig, nil
}

func UpdateCurrentAccount(account string) error {
	authConfig, err := NewAuthenticationConfigFromFile()
	if err != nil {
		return fmt.Errorf("unable to retrieve the existing authentication configuration; %w", err)
	}

	if _, ok := authConfig.Credentials[account]; !ok {
		return fmt.Errorf("account %s is not found", account)
	}

	authConfig.CurrentAccount = account

	if err := saveCredentialsConfigFile(authConfig); err != nil {
		return fmt.Errorf("unable to save the authentication configuration to file; %w", err)
	}

	return nil
}

func credentialsConfigFile() string {
	return filepath.Join(configDir(), credentialsFileName)
}

func configDir() string {
	rootDir, err := os.UserConfigDir()
	if err != nil {
		rootDir = "."
	}

	return filepath.Join(rootDir, internal.ApplicationName)
}

func ensureConfigDir() error {
	dir := configDir()

	if _, err := os.Stat(dir); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			if err := os.MkdirAll(dir, 0o750); err != nil {
				return fmt.Errorf("unable to create %s; %w", dir, err)
			}
		} else {
			return fmt.Errorf("unknown error received when running stat on %s; %w", dir, err)
		}
	}

	return nil
}

func saveCredentialsConfigFile(authConfig CredentialsConfig) error {
	file, err := os.Create(credentialsConfigFile())
	if err != nil {
		return fmt.Errorf("unable to open the config file; %w", err)
	}

	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "    ")

	if err := encoder.Encode(authConfig); err != nil {
		return fmt.Errorf("unable to save the JSON data to the authentication config file; %w", err)
	}

	return nil
}
