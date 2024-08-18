package utilities_test

import (
	"testing"

	"codeflow.dananglin.me.uk/apollo/enbas/internal/utilities"
)

func TestGetFQDN(t *testing.T) {
	cases := []struct {
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

	for _, c := range cases {
		got := utilities.GetFQDN(c.instance)
		if c.want != got {
			t.Errorf("Unexpected result: want %s, got %s", c.want, got)
		} else {
			t.Logf("Expected result: got %s", got)
		}
	}
}
