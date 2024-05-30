package utilities

import (
	"fmt"
	"os"
)

func ReadFile(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("unable to read the data from the file; %w", err)
	}

	return string(data), nil
}
