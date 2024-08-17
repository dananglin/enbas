package utilities_test

import (
	"fmt"
	"os"
	"path/filepath"
)

func projectRoot() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("unable to get the current working directory, %w", err)
	}

	return filepath.Join(cwd, "..", ".."), nil
}
