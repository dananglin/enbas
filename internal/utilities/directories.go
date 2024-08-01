package utilities

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"codeflow.dananglin.me.uk/apollo/enbas/internal"
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

func CalculateCacheDir(cacheDir, instanceFQDN string) string {
	if cacheDir != "" {
		return cacheDir
	}

	cacheRoot, err := os.UserCacheDir()
	if err != nil {
		return filepath.Join(os.Getenv("HOME"), "."+internal.ApplicationName, "cache")
	}

	return filepath.Join(cacheRoot, internal.ApplicationName, instanceFQDN)
}

func EnsureDirectory(dir string) error {
	if _, err := os.Stat(dir); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			if err := os.MkdirAll(dir, 0o750); err != nil {
				return fmt.Errorf("unable to create %s: %w", dir, err)
			}
		} else {
			return fmt.Errorf("received an unknown error after getting the directory information: %w", err)
		}
	}

	return nil
}
