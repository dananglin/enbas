package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"codeflow.dananglin.me.uk/apollo/enbas/internal/config"
)

func TestConfigFile(t *testing.T) {
	t.Log("Testing saving and loading the configuration")

	configFilepath := filepath.Join("testdata", "config", "config.json")

	t.Run("Save the default configuration to file", testSaveDefaultConfigToFile(configFilepath))
	t.Run("Load the configuration from file", testLoadConfigFromFile(configFilepath))

	if err := os.Remove(configFilepath); err != nil {
		t.Fatalf(
			"received an error after trying to clean up the test configuration at %q: %v",
			configFilepath,
			err,
		)
	}
}

func testSaveDefaultConfigToFile(configFilepath string) func(t *testing.T) {
	return func(t *testing.T) {
		if err := config.SaveInitialConfigToFile(configFilepath); err != nil {
			t.Fatalf("Unable to save the default configuration within %q: %v", configFilepath, err)
		}

		fileExists, err := config.FileExists(configFilepath)
		if err != nil {
			t.Fatalf("Unable to determine if the configuration exists within %q: %v", configFilepath, err)
		}

		if !fileExists {
			t.Fatalf("The configuration does not appear to exist within %q", configFilepath)
		}
	}
}

func testLoadConfigFromFile(configFilepath string) func(t *testing.T) {
	return func(t *testing.T) {
		config, err := config.NewConfigFromFile(configFilepath)
		if err != nil {
			t.Fatalf("Unable to load the configuration from file: %v", err)
		}

		expectedDefaultHTTPTimeout := 5

		if config.GTSClient.Timeout != 5 {
			t.Errorf(
				"Unexpected HTTP Timeout settings loaded from the configuration: want %d, got %d",
				expectedDefaultHTTPTimeout,
				config.GTSClient.Timeout,
			)
		} else {
			t.Logf(
				"Expected HTTP Timeout settings loaded from the configuration: got %d",
				config.GTSClient.Timeout,
			)
		}
	}
}
