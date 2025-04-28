package config

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"

	"codeflow.dananglin.me.uk/apollo/enbas/internal/info"
	"codeflow.dananglin.me.uk/apollo/enbas/internal/utilities"
)

const (
	defaultConfigFileName      string = "config.json"
	defaultCredentialsFileName string = "credentials.json"
)

func ConfigPath(configFilepath string) (string, error) {
	return configPath(configFilepath)
}

func configPath(configFilepath string) (string, error) {
	if configFilepath != "" {
		return configFilepath, nil
	}

	return defaultConfigPath()
}

func defaultConfigPath() (string, error) {
	dir, err := defaultConfigDir()
	if err != nil {
		return "", fmt.Errorf("error calculating the default config directory: %w", err)
	}

	return filepath.Join(dir, defaultConfigFileName), nil
}

func credentialsPath(credentialsFilepath string) (string, error) {
	if credentialsFilepath != "" {
		return credentialsFilepath, nil
	}

	return defaultCredentialsFilepath()
}

func defaultCredentialsFilepath() (string, error) {
	dir, err := defaultConfigDir()
	if err != nil {
		return "", fmt.Errorf("error calculating the default config directory: %w", err)
	}

	return filepath.Join(dir, defaultCredentialsFileName), nil
}

func defaultConfigDir() (string, error) {
	configHome, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("error calculating your config home directory: %w", err)
	}

	return filepath.Join(configHome, info.ApplicationName), nil
}

func defaultCacheDir() (string, error) {
	cacheHome, err := os.UserCacheDir()
	if err != nil {
		return "", fmt.Errorf("error calculating your cache home directory: %w", err)
	}

	return filepath.Join(cacheHome, info.ApplicationName), nil
}

func newSocketPath() (string, error) {
	runtimeDir := os.Getenv("XDG_RUNTIME_DIR")
	if runtimeDir == "" {
		return "", nil
	}

	randBytes := make([]byte, 4)

	if _, err := rand.Read(randBytes); err != nil {
		return "", fmt.Errorf("error creating random bytes: %w", err)
	}

	path, err := utilities.AbsolutePath(filepath.Join(
		runtimeDir,
		info.ApplicationName,
		"server."+hex.EncodeToString(randBytes)+".socket",
	))
	if err != nil {
		return "", fmt.Errorf(
			"error calculating the absolute path to the socket file: %w",
			err,
		)
	}

	return path, nil
}
