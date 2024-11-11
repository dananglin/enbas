package utilities_test

import (
	"os"
	"slices"
	"testing"

	"codeflow.dananglin.me.uk/apollo/enbas/internal/utilities"
)

func TestReadContents(t *testing.T) {
	t.Log("Testing reading contents from file")

	testCases := []struct {
		content string
		want    string
	}{
		{
			content: "A simple one-line status.",
			want:    "A simple one-line status.",
		},
		{
			content: "file@testdata/statuses/status.txt",
			want:    "Hello World! This is a test status from a text file.",
		},
	}

	for _, tc := range slices.All(testCases) {
		got, err := utilities.ReadContents(tc.content)
		if err != nil {
			t.Fatalf("Unable to read the contents from %q: %v", tc.content, err)
		}

		if got != tc.want {
			t.Errorf(
				"Unexpected content read from %q: want %q, got %q",
				tc.content,
				tc.want,
				got,
			)
		} else {
			t.Logf(
				"Expected content read from %q: got %q",
				tc.content,
				got,
			)
		}
	}
}

func TestSaveTextToFile(t *testing.T) {
	t.Log("Testing saving contents to a text file")

	content := "A test status"
	path := "testdata/statuses/" + t.Name() + ".golden"

	defer os.Remove(path)

	if err := utilities.SaveTextToFile(path, content); err != nil {
		t.Fatalf("Unable to save the contents to %q: %v", path, err)
	} else {
		t.Logf("Successfully saved the contents to %q", path)
	}

	path = "file@" + path

	got, err := utilities.ReadContents(path)
	if err != nil {
		t.Fatalf("Unable to read the contents from %q: %v", path, err)
	} else {
		t.Logf("Successfully read the contents from %q", path)
	}

	if got != content {
		t.Errorf(
			"Unexpected content read from %q: want %q, got %q",
			path,
			content,
			got,
		)
	} else {
		t.Logf(
			"Expected content read from %q: got %q",
			path,
			got,
		)
	}
}
