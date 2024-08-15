package utilities

import (
	"bufio"
	"errors"
	"fmt"
	"os"
)

func ReadFile(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", fmt.Errorf("unable to open %q: %w", path, err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var lines []string

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("received an error after scanning the contents from %q: %w", path, err)
	}

	if lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}

	return strings.Join(lines, "\n"), nil
}

func FileExists(path string) (bool, error) {
	if _, err := os.Stat(path); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return false, nil
		}

		return false, fmt.Errorf("unable to check if the file exists: %w", err)
	}

	return true, nil
}
