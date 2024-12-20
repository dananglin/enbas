package utilities

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
)

const filePrefix string = "file@"

func ReadContents(text string) (string, error) {
	if !strings.HasPrefix(text, filePrefix) {
		return text, nil
	}

	return ReadTextFile(strings.TrimPrefix(text, filePrefix))
}

func ReadTextFile(path string) (string, error) {
	file, err := OpenFile(path)
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

func SaveTextToFile(path, text string) error {
	file, err := CreateFile(path)
	if err != nil {
		return fmt.Errorf("unable to open %q: %w", path, err)
	}
	defer file.Close()

	if _, err := fmt.Fprint(file, text); err != nil {
		return fmt.Errorf("received an error writing the text to the file: %w", err)
	}

	return nil
}

func OpenFile(path string) (*os.File, error) {
	absPath, err := AbsolutePath(path)
	if err != nil {
		return nil, fmt.Errorf("error calculating the absolute path: %w", err)
	}

	file, err := os.Open(absPath) // #nosec G304 -- The path is cleaned when calculating the absolute path.
	if err != nil {
		return nil, fmt.Errorf("error opening the file: %w", err)
	}

	return file, nil
}

func CreateFile(path string) (*os.File, error) {
	absPath, err := AbsolutePath(path)
	if err != nil {
		return nil, fmt.Errorf("error calculating the absolute path: %w", err)
	}

	file, err := os.Create(absPath) // #nosec G304 -- The path is cleaned when calculating the absolute path.
	if err != nil {
		return nil, fmt.Errorf("error opening the file: %w", err)
	}

	return file, nil
}
