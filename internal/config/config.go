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

type AuthenticationConfig struct {
	CurrentAccount  string                    `json:"currentAccount"`
	Authentications map[string]Authentication `json:"authentications"`
}

type Authentication struct {
	Instance     string `json:"instance"`
	ClientID     string `json:"clientId"`
	ClientSecret string `json:"clientSecret"`
	AccessToken  string `json:"accessToken"`
}

func SaveAuthentication(username string, authentication Authentication) (string, error) {
	if err := ensureConfigDir(); err != nil {
		return "", fmt.Errorf("unable to ensure the configuration directory; %w", err)
	}

	var authConfig AuthenticationConfig

	filepath := authenticationConfigFile()

	if _, err := os.Stat(filepath); err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return "", fmt.Errorf("unknown error received when running stat on %s; %w", filepath, err)
		}

		authConfig.Authentications = make(map[string]Authentication)
	} else {
		authConfig, err = NewAuthenticationConfigFromFile()
		if err != nil {
			return "", fmt.Errorf("unable to retrieve the existing authentication configuration; %w", err)
		}
	}

	instance := ""

	if strings.HasPrefix(authentication.Instance, "https://") {
		instance = strings.TrimPrefix(authentication.Instance, "https://")
	} else if strings.HasPrefix(authentication.Instance, "http://") {
		instance = strings.TrimPrefix(authentication.Instance, "http://")
	}

	authenticationName := username + "@" + instance

	authConfig.CurrentAccount = authenticationName

	authConfig.Authentications[authenticationName] = authentication

	if err := saveAuthenticationFile(authConfig); err != nil {
		return "", fmt.Errorf("unable to save the authentication configuration to file; %w", err)
	}

	return authenticationName, nil
}

func NewAuthenticationConfigFromFile() (AuthenticationConfig, error) {
	path := authenticationConfigFile()

	file, err := os.Open(path)
	if err != nil {
		return AuthenticationConfig{}, fmt.Errorf("unable to open %s, %w", path, err)
	}
	defer file.Close()

	var authConfig AuthenticationConfig

	if err := json.NewDecoder(file).Decode(&authConfig); err != nil {
		return AuthenticationConfig{}, fmt.Errorf("unable to decode the JSON data; %w", err)
	}

	return authConfig, nil
}

func UpdateCurrentAccount(account string) error {
	authConfig, err := NewAuthenticationConfigFromFile()
	if err != nil {
		return fmt.Errorf("unable to retrieve the existing authentication configuration; %w", err)
	}

	if _, ok := authConfig.Authentications[account]; !ok {
		return fmt.Errorf("account %s is not found", account)
	}

	authConfig.CurrentAccount = account

	if err := saveAuthenticationFile(authConfig); err != nil {
		return fmt.Errorf("unable to save the authentication configuration to file; %w", err)
	}

	return nil
}

func authenticationConfigFile() string {
	return filepath.Join(configDir(), "authentications.json")
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

func saveAuthenticationFile(authConfig AuthenticationConfig) error {
	file, err := os.Create(authenticationConfigFile())
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
