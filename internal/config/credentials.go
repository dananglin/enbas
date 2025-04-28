package config

import (
	"encoding/json"
	"fmt"
	"path/filepath"

	"codeflow.dananglin.me.uk/apollo/enbas/internal/utilities"
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

// SaveCredentials saves the credentials into the specified credentials file. If the path to the credentials
// file is not supplied the default path is used which is $XDG_CONFIG_HOME/enbas/credentials.json.
// If the file's parent directory is not present it will be created.
func SaveCredentials(credentialsFilepath, username string, credentials Credentials) (string, error) {
	path, err := credentialsPath(credentialsFilepath)
	if err != nil {
		return "", fmt.Errorf("error calculating the path to the credentials file: %w", err)
	}

	// ensure that the parent directory exists.
	parentDir := filepath.Dir(path)

	if err := utilities.EnsureDirectory(parentDir); err != nil {
		return "", fmt.Errorf("error ensuring the presence of the parent directory of the credentials file: %w", err)
	}

	var authConfig CredentialsConfig

	exists, err := utilities.FileExists(path)
	if err != nil {
		return "", fmt.Errorf("unexpected error received after checking for the existence of the credentials file at %q: %w", path, err)
	}

	if !exists {
		authConfig.Credentials = make(map[string]Credentials)
	} else {
		authConfig, err = NewCredentialsConfigFromFile(path)
		if err != nil {
			return "", fmt.Errorf("error retrieving the exsting credentials from the credentials file: %w", err)
		}
	}

	instance := utilities.GetFQDN(credentials.Instance)

	authenticationName := username + "@" + instance

	authConfig.CurrentAccount = authenticationName

	authConfig.Credentials[authenticationName] = credentials

	if err := saveCredentialsConfigFile(authConfig, path); err != nil {
		return "", fmt.Errorf("unable to save the authentication configuration to file: %w", err)
	}

	return authenticationName, nil
}

// UpdateCurrentAccount updates the name of the current account in the credentials config file.
func UpdateCurrentAccount(account string, filePath string) error {
	credentialsConfig, err := NewCredentialsConfigFromFile(filePath)
	if err != nil {
		return fmt.Errorf("error retrieving the existing credentials from the credentials file: %w", err)
	}

	if _, ok := credentialsConfig.Credentials[account]; !ok {
		return CredentialsNotFoundError{account}
	}

	credentialsConfig.CurrentAccount = account

	if err := saveCredentialsConfigFile(credentialsConfig, filePath); err != nil {
		return fmt.Errorf("error saving the credentials to file: %w", err)
	}

	return nil
}

// NewCredentialsConfigFromFile creates a new CredentialsConfig value from reading
// the credentials file.
func NewCredentialsConfigFromFile(path string) (CredentialsConfig, error) {
	file, err := utilities.OpenFile(path)
	if err != nil {
		return CredentialsConfig{}, fmt.Errorf("error opening %s: %w", path, err)
	}
	defer file.Close()

	var authConfig CredentialsConfig

	if err := json.NewDecoder(file).Decode(&authConfig); err != nil {
		return CredentialsConfig{}, fmt.Errorf("error decoding the JSON data: %w", err)
	}

	return authConfig, nil
}

func saveCredentialsConfigFile(authConfig CredentialsConfig, filePath string) error {
	file, err := utilities.CreateFile(filePath)
	if err != nil {
		return fmt.Errorf("error creating the file at %s: %w", filePath, err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "    ")

	if err := encoder.Encode(authConfig); err != nil {
		return fmt.Errorf("error saving the JSON data to the credentials file: %w", err)
	}

	return nil
}
