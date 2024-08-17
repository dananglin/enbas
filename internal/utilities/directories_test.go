package utilities_test

import (
	"path/filepath"
	"testing"

	"codeflow.dananglin.me.uk/apollo/enbas/internal/utilities"
)

func TestCalculateMediaCacheDir(t *testing.T) {
	t.Parallel()

	projectDir, err := projectRoot()
	if err != nil {
		t.Fatalf("Unable to get the project root directory: %v", err)
	}

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

func TestCalculateMediaCacheDirWithXDG(t *testing.T) {
	t.Setenv("XDG_CACHE_HOME", "/home/enbas/.cache")

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

func TestCalculateStatusesCacheDir(t *testing.T) {
	t.Parallel()

	projectDir, err := projectRoot()
	if err != nil {
		t.Fatalf("Unable to get the project root directory: %v", err)
	}

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
