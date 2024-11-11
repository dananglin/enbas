package utilities_test

import (
	"slices"
	"testing"

	"codeflow.dananglin.me.uk/apollo/enbas/internal/utilities"
)

func TestGetFQDN(t *testing.T) {
	testCases := []struct {
		instance string
		want     string
	}{
		{
			instance: "https://gts.red-crow.private",
			want:     "gts.red-crow.private",
		},
		{
			instance: "http://gotosocial.yellow-desert.social",
			want:     "gotosocial.yellow-desert.social",
		},
		{
			instance: "fedi.blue-mammoth.party",
			want:     "fedi.blue-mammoth.party",
		},
	}

	for _, tc := range slices.All(testCases) {
		got := utilities.GetFQDN(tc.instance)
		if tc.want != got {
			t.Errorf("Unexpected result received: want %q, got %q", tc.want, got)
		} else {
			t.Logf("Expected result received: got %q", got)
		}
	}
}
