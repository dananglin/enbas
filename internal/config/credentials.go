package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"codeflow.dananglin.me.uk/apollo/enbas/internal/utilities"
)

const (
	defaultCredentialsFileName = "credentials.json"
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
func SaveCredentials(filePath, username string, credentials Credentials) (string, error) {
	part := filepath.Dir(filePath)

	// ensure that the directory exists.
	credentialsDir, err := utilities.CalculateConfigDir(part)
	if err != nil {
		return "", fmt.Errorf("unable to calculate the directory to your credentials file: %w", err)
	}

	if err := utilities.EnsureDirectory(credentialsDir); err != nil {
		return "", fmt.Errorf("unable to ensure the configuration directory: %w", err)
	}

	var authConfig CredentialsConfig

	if _, err := os.Stat(filePath); err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return "", fmt.Errorf("unknown error received when running stat on %s: %w", filePath, err)
		}

		authConfig.Credentials = make(map[string]Credentials)
	} else {
		authConfig, err = NewCredentialsConfigFromFile(filePath)
		if err != nil {
			return "", fmt.Errorf("unable to retrieve the existing authentication configuration: %w", err)
		}
	}

	instance := utilities.GetFQDN(credentials.Instance)

	authenticationName := username + "@" + instance

	authConfig.CurrentAccount = authenticationName

	authConfig.Credentials[authenticationName] = credentials

	if err := saveCredentialsConfigFile(authConfig, filePath); err != nil {
		return "", fmt.Errorf("unable to save the authentication configuration to file: %w", err)
	}

	return authenticationName, nil
}

// UpdateCurrentAccount updates the name of the current account in the credentials config file.
func UpdateCurrentAccount(account string, filePath string) error {
	credentialsConfig, err := NewCredentialsConfigFromFile(filePath)
	if err != nil {
		return fmt.Errorf("unable to retrieve the existing authentication configuration: %w", err)
	}

	if _, ok := credentialsConfig.Credentials[account]; !ok {
		return CredentialsNotFoundError{account}
	}

	credentialsConfig.CurrentAccount = account

	if err := saveCredentialsConfigFile(credentialsConfig, filePath); err != nil {
		return fmt.Errorf("unable to save the authentication configuration to file: %w", err)
	}

	return nil
}

// NewCredentialsConfigFromFile creates a new CredentialsConfig value from reading
// the credentials file.
func NewCredentialsConfigFromFile(path string) (CredentialsConfig, error) {
	file, err := utilities.OpenFile(path)
	if err != nil {
		return CredentialsConfig{}, fmt.Errorf("unable to open %s: %w", path, err)
	}
	defer file.Close()

	var authConfig CredentialsConfig

	if err := json.NewDecoder(file).Decode(&authConfig); err != nil {
		return CredentialsConfig{}, fmt.Errorf("unable to decode the JSON data: %w", err)
	}

	return authConfig, nil
}

func saveCredentialsConfigFile(authConfig CredentialsConfig, filePath string) error {
	file, err := utilities.CreateFile(filePath)
	if err != nil {
		return fmt.Errorf("unable to create the file at %s: %w", filePath, err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "    ")

	if err := encoder.Encode(authConfig); err != nil {
		return fmt.Errorf("unable to save the JSON data to the authentication config file: %w", err)
	}

	return nil
}

func defaultCredentialsConfigFile(configDir string) (string, error) {
	dir, err := utilities.CalculateConfigDir(configDir)
	if err != nil {
		return "", fmt.Errorf("unable to calculate the config directory: %w", err)
	}

	path, err := utilities.AbsolutePath(filepath.Join(dir, defaultCredentialsFileName))
	if err != nil {
		return "", fmt.Errorf("unable to get the absolute path to the credentials config file: %w", err)
	}

	return path, nil
}
