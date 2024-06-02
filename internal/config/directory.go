// SPDX-FileCopyrightText: 2024 Dan Anglin <d.n.i.anglin@gmail.com>
//
// SPDX-License-Identifier: GPL-3.0-or-later

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
				return fmt.Errorf("unable to create %s: %w", configDir, err)
			}
		} else {
			return fmt.Errorf("unknown error received after getting the config directory information: %w", err)
		}
	}

	return nil
}
