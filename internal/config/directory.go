package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"codeflow.dananglin.me.uk/apollo/enbas/internal"
)

func calculateConfigDir(configDir string) string {
	if configDir != "" {
		return configDir
	}

	rootDir, err := os.UserConfigDir()
	if err != nil {
		rootDir = "."
	}

	return filepath.Join(rootDir, internal.ApplicationName)
}

func ensureConfigDir(configDir string) error {
	if _, err := os.Stat(configDir); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			if err := os.MkdirAll(configDir, 0o750); err != nil {
				return fmt.Errorf("unable to create %s; %w", configDir, err)
			}
		} else {
			return fmt.Errorf("unknown error received when running stat on %s; %w", configDir, err)
		}
	}

	return nil
}