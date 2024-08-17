package utilities

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"codeflow.dananglin.me.uk/apollo/enbas/internal"
)

const (
	cacheMediaDir    = "media"
	cacheStatusesDir = "statuses"
)

func CalculateConfigDir(configDir string) string {
	if configDir != "" {
		return configDir
	}

	configRoot, err := os.UserConfigDir()
	if err != nil {
		return filepath.Join(os.Getenv("HOME"), "."+internal.ApplicationName, "config")
	}

	return filepath.Join(configRoot, internal.ApplicationName)
}

func CalculateMediaCacheDir(cacheRoot, instance string) (string, error) {
	cacheDir, err := calculateCacheDir(cacheRoot, instance)
	if err != nil {
		return "", fmt.Errorf("unable to calculate the cache directory: %w", err)
	}

	return filepath.Join(cacheDir, cacheMediaDir), nil
}

func CalculateStatusesCacheDir(cacheRoot, instance string) (string, error) {
	cacheDir, err := calculateCacheDir(cacheRoot, instance)
	if err != nil {
		return "", fmt.Errorf("unable to calculate the cache directory: %w", err)
	}

	return filepath.Join(cacheDir, cacheStatusesDir), nil
}

func calculateCacheDir(cacheRoot, instance string) (string, error) {
	fqdn := GetFQDN(instance)

	if cacheRoot != "" {
		return filepath.Join(cacheRoot, fqdn), nil
	}

	cacheRoot, err := os.UserCacheDir()
	if err != nil {
		return "", fmt.Errorf("unable to get your default cache directory: %w", err)
	}

	return filepath.Join(cacheRoot, internal.ApplicationName, fqdn), nil
}

func EnsureDirectory(dir string) error {
	if _, err := os.Stat(dir); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			if err := os.MkdirAll(dir, 0o750); err != nil {
				return fmt.Errorf("unable to create %s: %w", dir, err)
			}
		} else {
			return fmt.Errorf(
				"received an unknown error after getting the directory information: %w",
				err,
			)
		}
	}

	return nil
}
