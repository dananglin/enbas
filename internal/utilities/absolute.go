package utilities

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func AbsolutePath(path string) (string, error) {
	if strings.HasPrefix(path, "~") {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("unable to get user's home directory; %w", err)
		}

		path = filepath.Join(homeDir, path[1:])
	}

	absPath, err := filepath.Abs(path)
	if err != nil {
		return "", fmt.Errorf("unable to get the absolute path: %w", err)
	}

	return absPath, nil
}
