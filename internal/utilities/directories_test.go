package utilities_test

import (
	"path/filepath"
	"testing"

	"codeflow.dananglin.me.uk/apollo/enbas/internal/utilities"
)

func TestDirectoryCalculations(t *testing.T) {
	t.Log("Testing the directory calculations")

	projectDir, err := projectRoot()
	if err != nil {
		t.Fatalf("Unable to get the project root directory: %v", err)
	}

	t.Setenv("XDG_CACHE_HOME", "/home/enbas/.cache")
	t.Setenv("XDG_CONFIG_HOME", "/home/enbas/.config")

	t.Run("Media Cache Directory Calculation", testCalculateMediaCacheDir(projectDir))
	t.Run("Media Cache Directory Calculation (with XDG_CACHE_HOME)", testCalculateMediaCacheDirWithXDG)
	t.Run("Statuses Cache Directory Calculation", testCalculateStatusesCacheDir(projectDir))
	t.Run("Configuration Directory Calculation", testCalculateConfigDir(projectDir))
	t.Run("Configuration Directory Calculation (with XDG_CONFIG_HOME)", testCalculateConfigCacheDirWithXDG)
}

func testCalculateMediaCacheDir(projectDir string) func(t *testing.T) {
	return func(t *testing.T) {
		cacheRoot := filepath.Join(projectDir, "test", "cache")
		instance := "http://gotosocial.yellow-desert.social"

		got, err := utilities.CalculateMediaCacheDir(cacheRoot, instance)
		if err != nil {
			t.Fatalf("Unable to calculate the media cache directory: %v", err)
		}

		want := projectDir + "/test/cache/gotosocial.yellow-desert.social/media"

		if got != want {
			t.Errorf("Unexpected media cache directory calculated: want %s, got %s", want, got)
		} else {
			t.Logf("Expected media cache directory calculated: got %s", got)
		}
	}
}

func testCalculateMediaCacheDirWithXDG(t *testing.T) {
	cacheRoot := ""
	instance := "https://gotosocial.yellow-desert.social"

	got, err := utilities.CalculateMediaCacheDir(cacheRoot, instance)
	if err != nil {
		t.Fatalf("Unable to calculate the media cache directory: %v", err)
	}

	want := "/home/enbas/.cache/enbas/gotosocial.yellow-desert.social/media"

	if got != want {
		t.Errorf("Unexpected media cache directory calculated: want %s, got %s", want, got)
	} else {
		t.Logf("Expected media cache directory calculated: got %s", got)
	}
}

func testCalculateStatusesCacheDir(projectDir string) func(t *testing.T) {
	return func(t *testing.T) {
		cacheRoot := filepath.Join(projectDir, "test", "cache")
		instance := "https://fedi.blue-mammoth.party"

		got, err := utilities.CalculateStatusesCacheDir(cacheRoot, instance)
		if err != nil {
			t.Fatalf("Unable to calculate the statuses cache directory: %v", err)
		}

		want := projectDir + "/test/cache/fedi.blue-mammoth.party/statuses"

		if got != want {
			t.Errorf("Unexpected statuses cache directory calculated: want %s, got %s", want, got)
		} else {
			t.Logf("Expected statuses cache directory calculated: got %s", got)
		}
	}
}

func testCalculateConfigDir(projectDir string) func(t *testing.T) {
	return func(t *testing.T) {
		configDir := filepath.Join(projectDir, "test", "config")

		got, err := utilities.CalculateConfigDir(configDir)
		if err != nil {
			t.Fatalf("Unable to calculate the config directory: %v", err)
		}

		want := projectDir + "/test/config"

		if got != want {
			t.Errorf("Unexpected config directory calculated: want %s, got %s", want, got)
		} else {
			t.Logf("Expected config directory calculated: got %s", got)
		}
	}
}

func testCalculateConfigCacheDirWithXDG(t *testing.T) {
	configDir := ""

	got, err := utilities.CalculateConfigDir(configDir)
	if err != nil {
		t.Fatalf("Unable to calculate the config directory: %v", err)
	}

	want := "/home/enbas/.config/enbas"

	if got != want {
		t.Errorf("Unexpected config directory calculated: want %s, got %s", want, got)
	} else {
		t.Logf("Expected config directory calculated: got %s", got)
	}
}
