// SPDX-FileCopyrightText: 2024 Dan Anglin <d.n.i.anglin@gmail.com>
//
// SPDX-License-Identifier: GPL-3.0-or-later

package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
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

type CredentialsNotFoundError struct {
	AccountName string
}

func (e CredentialsNotFoundError) Error() string {
	return "unable to find the credentials for the account '" + e.AccountName + "'"
}

// SaveCredentials saves the credentials into the credentials file within the specified configuration
// directory. If the directory is not specified then the default directory is used. If the directory
// is not present, it will be created.
func SaveCredentials(configDir, username string, credentials Credentials) (string, error) {
	if err := ensureConfigDir(calculateConfigDir(configDir)); err != nil {
		return "", fmt.Errorf("unable to ensure the configuration directory; %w", err)
	}

	var authConfig CredentialsConfig

	filepath := credentialsConfigFile(configDir)

	if _, err := os.Stat(filepath); err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return "", fmt.Errorf("unknown error received when running stat on %s; %w", filepath, err)
		}

		authConfig.Credentials = make(map[string]Credentials)
	} else {
		authConfig, err = NewCredentialsConfigFromFile(configDir)
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

	if err := saveCredentialsConfigFile(authConfig, configDir); err != nil {
		return "", fmt.Errorf("unable to save the authentication configuration to file; %w", err)
	}

	return authenticationName, nil
}

func UpdateCurrentAccount(account string, configDir string) error {
	credentialsConfig, err := NewCredentialsConfigFromFile(configDir)
	if err != nil {
		return fmt.Errorf("unable to retrieve the existing authentication configuration; %w", err)
	}

	if _, ok := credentialsConfig.Credentials[account]; !ok {
		return CredentialsNotFoundError{account}
	}

	credentialsConfig.CurrentAccount = account

	if err := saveCredentialsConfigFile(credentialsConfig, configDir); err != nil {
		return fmt.Errorf("unable to save the authentication configuration to file; %w", err)
	}

	return nil
}

func NewCredentialsConfigFromFile(configDir string) (CredentialsConfig, error) {
	path := credentialsConfigFile(configDir)

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

func saveCredentialsConfigFile(authConfig CredentialsConfig, configDir string) error {
	path := credentialsConfigFile(configDir)

	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("unable to open %s; %w", path, err)
	}

	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "    ")

	if err := encoder.Encode(authConfig); err != nil {
		return fmt.Errorf("unable to save the JSON data to the authentication config file; %w", err)
	}

	return nil
}

func credentialsConfigFile(configDir string) string {
	return filepath.Join(calculateConfigDir(configDir), credentialsFileName)
}
