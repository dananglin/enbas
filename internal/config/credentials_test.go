package config_test

import (
	"maps"
	"os"
	"path/filepath"
	"testing"

	"codeflow.dananglin.me.uk/apollo/enbas/internal/config"
)

func TestCredentialsFile(t *testing.T) {
	t.Log("Testing saving and loading credentials from file")

	credentialsFile := filepath.Join("testdata", "config", "credentials.json")

	credentialsMap := map[string]config.Credentials{
		"admin": {
			Instance:     "https://gts.red-crow.private",
			ClientID:     "01EOB91DVQGPA364QK32TM3LXR1998BMXSZE4",
			ClientSecret: "ffd76025-4b23-4ce6-b8ea-077ce3cadf5a",
			AccessToken:  "C9VDXGGRPZ0448SH562N6N6893VNPGJMGJ336TXLMH8RXGWF4",
		},
		"bobby": {
			Instance:     "https://gts.red-crow.private",
			ClientID:     "01CUVHR6LIST7Q6R25Z9Y14WZK780V91S9VQB",
			ClientSecret: "379fc272-c7cc-4ccb-8461-f3f71207f798",
			AccessToken:  "F0YWQG1R4DDAMXGBZ514BCW7ATWN6JRGLDRUZO4RFAMTT6J38",
		},
		"app": {
			Instance:     "https://gts.red-crow.private",
			ClientID:     "01HLZY7XCD60564OP3RG6FZTOAD3LGF0R8SEK",
			ClientSecret: "dfd8e954-53b1-4f00-9c09-0b181f44bb79",
			AccessToken:  "JZ2PZ4YNE1BB38VMRIQ7DNWXKZE6B1EBV310RNC53KQCVHXGB",
		},
	}

	t.Run("Saving credentials to file", testSaveCredentials(credentialsFile, credentialsMap))

	expectedCurrentAccount := "bobby@gts.red-crow.private"
	t.Run("Updating the current account in the credentials file", testUpdateCurrentAccount(expectedCurrentAccount, credentialsFile))

	t.Run("Loading the credentials from file", testLoadCredentialsConfigFromFile(credentialsFile, expectedCurrentAccount))

	if err := os.Remove(credentialsFile); err != nil {
		t.Fatalf(
			"received an error after trying to clean up the test configuration at %q: %v",
			credentialsFile,
			err,
		)
	}
}

func testSaveCredentials(credentialsFile string, credentialsMap map[string]config.Credentials) func(t *testing.T) {
	return func(t *testing.T) {
		for username, credentials := range maps.All(credentialsMap) {
			if _, err := config.SaveCredentials(credentialsFile, username, credentials); err != nil {
				t.Fatalf(
					"Unable to save the credentials for %s to %q: %v",
					username,
					credentialsFile,
					err,
				)
			}
		}

		t.Log("All credentials saved to file.")
	}
}

func testUpdateCurrentAccount(account, credentialsFile string) func(t *testing.T) {
	return func(t *testing.T) {
		if err := config.UpdateCurrentAccount(account, credentialsFile); err != nil {
			t.Fatalf("Unable to update the current account to %q: %v", account, err)
		}

		t.Logf("Successfully updated the current account.")
	}
}

func testLoadCredentialsConfigFromFile(credentialsFile string, expectedCurrentAccount string) func(t *testing.T) {
	return func(t *testing.T) {
		credentials, err := config.NewCredentialsConfigFromFile(credentialsFile)
		if err != nil {
			t.Fatalf(
				"Unable to load the credentials configuration from %q: %v",
				credentialsFile,
				err,
			)
		}

		if credentials.CurrentAccount != expectedCurrentAccount {
			t.Errorf(
				"Unexpected current account found in the credentials configuration file: want %s, got %s",
				expectedCurrentAccount,
				credentials.CurrentAccount,
			)
		} else {
			t.Logf("Expected current account found in the credentials configuration file: got %s", credentials.CurrentAccount)
		}
	}
}
