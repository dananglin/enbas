package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"codeflow.dananglin.me.uk/apollo/enbas/internal/config"
)

func TestConfigFile(t *testing.T) {
	t.Log("Testing saving and loading the configuration")

	configDir := filepath.Join("testdata", "config")

	t.Run("Save the default configuration to file", testSaveDefaultConfigToFile(configDir))
	t.Run("Load the configuration from file", testLoadConfigFromFile(configDir))

	expectedConfigFile := filepath.Join(configDir, "config.json")
	if err := os.Remove(expectedConfigFile); err != nil {
		t.Fatalf(
			"received an error after trying to clean up the test configuration at %q: %v",
			expectedConfigFile,
			err,
		)
	}
}

func testSaveDefaultConfigToFile(configDir string) func(t *testing.T) {
	return func(t *testing.T) {
		if err := config.SaveDefaultConfigToFile(configDir); err != nil {
			t.Fatalf("Unable to save the default configuration within %q: %v", configDir, err)
		}

		fileExists, err := config.FileExists(configDir)
		if err != nil {
			t.Fatalf("Unable to determine if the configuration exists within %q: %v", configDir, err)
		}

		if !fileExists {
			t.Fatalf("The configuration does not appear to exist within %q", configDir)
		}
	}
}

func testLoadConfigFromFile(configDir string) func(t *testing.T) {
	return func(t *testing.T) {
		config, err := config.NewConfigFromFile(configDir)
		if err != nil {
			t.Fatalf("Unable to load the configuration from file: %v", err)
		}

		expectedDefaultHTTPTimeout := 5

		if config.HTTP.Timeout != expectedDefaultHTTPTimeout {
			t.Errorf(
				"Unexpected HTTP Timeout settings loaded from the configuration: want %d, got %d",
				expectedDefaultHTTPTimeout,
				config.HTTP.Timeout,
			)
		} else {
			t.Logf(
				"Expected HTTP Timeout settings loaded from the configuration: got %d",
				config.HTTP.Timeout,
			)
		}
	}
}
