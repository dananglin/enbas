package utilities_test

import (
	"path/filepath"
	"testing"

	"codeflow.dananglin.me.uk/apollo/enbas/internal/utilities"
)

func TestDirectoryCalculations(t *testing.T) {
	t.Log("Testing the directory calculations")

	t.Setenv("XDG_CACHE_HOME", "/home/enbas/.cache")
	t.Setenv("XDG_CONFIG_HOME", "/home/enbas/.config")

	t.Run("Media Cache Directory Calculation", testCalculateMediaCacheDir)
	t.Run("Media Cache Directory Calculation (with XDG_CACHE_HOME)", testCalculateMediaCacheDirWithXDG)
	t.Run("Statuses Cache Directory Calculation", testCalculateStatusesCacheDir)
	t.Run("Configuration Directory Calculation", testCalculateConfigDir)
	t.Run("Configuration Directory Calculation (with XDG_CONFIG_HOME)", testCalculateConfigCacheDirWithXDG)
}

func testCalculateMediaCacheDir(t *testing.T) {
	cacheRoot := filepath.Join("testdata", "test", "cache")

	absCacheRoot, err := utilities.AbsolutePath(cacheRoot)
	if err != nil {
		t.Fatalf(
			"FAILED test %s: Unable to calculate the absolute path of the root cache directory: %v",
			t.Name(),
			err,
		)
	} else {
		t.Logf("Absolute path of cache root: %s", absCacheRoot)
	}

	instance := "http://gotosocial.yellow-desert.social"

	got, err := utilities.CalculateMediaCacheDir(absCacheRoot, instance)
	if err != nil {
		t.Fatalf(
			"FAILED test %s: Unable to calculate the media cache directory: %v",
			t.Name(),
			err,
		)
	}

	want := absCacheRoot + "/gotosocial.yellow-desert.social/media"

	if got != want {
		t.Errorf(
			"FAILED test %s: Unexpected media cache directory calculated: want %s, got %s",
			t.Name(),
			want,
			got,
		)
	} else {
		t.Logf("Expected media cache directory calculated: got %s", got)
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

func testCalculateStatusesCacheDir(t *testing.T) {
	cacheRoot := filepath.Join("testdata", "test", "cache")

	absCacheRoot, err := utilities.AbsolutePath(cacheRoot)
	if err != nil {
		t.Fatalf(
			"FAILED test %s: Unable to calculate the absolute path of the root cache directory: %v",
			t.Name(),
			err,
		)
	} else {
		t.Logf("Absolute path of cache root: %s", absCacheRoot)
	}

	instance := "https://fedi.blue-mammoth.party"

	got, err := utilities.CalculateStatusesCacheDir(absCacheRoot, instance)
	if err != nil {
		t.Fatalf(
			"FAILED test %s: Unable to calculate the statuses cache directory: %v",
			t.Name(),
			err,
		)
	}

	want := absCacheRoot + "/fedi.blue-mammoth.party/statuses"

	if got != want {
		t.Errorf(
			"FAILED test %s: Unexpected statuses cache directory calculated: want %s, got %s",
			t.Name(),
			want,
			got,
		)
	} else {
		t.Logf("Expected statuses cache directory calculated: got %s", got)
	}
}

func testCalculateConfigDir(t *testing.T) {
	configDir := filepath.Join("testdata", "test", "config")

	absConfigDirPath, err := utilities.AbsolutePath(configDir)
	if err != nil {
		t.Fatalf(
			"FAILED test %s: Unable to calculate the absolute path of the config directory: %v",
			t.Name(),
			err,
		)
	} else {
		t.Logf("Absolute path of the config directory: %s", absConfigDirPath)
	}

	got, err := utilities.CalculateConfigDir(absConfigDirPath)
	if err != nil {
		t.Fatalf(
			"FAILED test %s: Unable to calculate the config directory: %v",
			t.Name(),
			err,
		)
	}

	if got != absConfigDirPath {
		t.Errorf(
			"FAILED test %s: Unexpected config directory calculated: want %s, got %s",
			t.Name(),
			absConfigDirPath,
			got,
		)
	} else {
		t.Logf("Expected config directory calculated: got %s", got)
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
