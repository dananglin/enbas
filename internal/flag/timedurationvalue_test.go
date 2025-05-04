package flag_test

import (
	"flag"
	"slices"
	"testing"

	internalFlag "codeflow.dananglin.me.uk/apollo/enbas/internal/flag"
)

func TestTimeDurationValue(t *testing.T) {
	parsingTests := []struct {
		input string
		want  string
	}{
		{
			input: `"1 day"`,
			want:  "24h0m0s",
		},
		{
			input: `"3 days, 5 hours, 39 minutes and 6 seconds"`,
			want:  "77h39m6s",
		},
		{
			input: `"1 minute and 30 seconds"`,
			want:  "1m30s",
		},
		{
			input: `"(7 seconds) (21 hours) (41 days)"`,
			want:  "1005h0m7s",
		},
	}

	for _, test := range slices.All(parsingTests) {
		args := []string{"--duration", test.input}

		t.Run("Flag parsing test: "+test.input, testTimeDurationValueParsing(args, test.want))
	}
}

func testTimeDurationValueParsing(args []string, want string) func(t *testing.T) {
	return func(t *testing.T) {
		flagset := flag.NewFlagSet("test", flag.ExitOnError)
		duration := internalFlag.NewTimeDurationValue(0)

		flagset.Var(&duration, "duration", "Duration value")

		if err := flagset.Parse(args); err != nil {
			t.Fatalf("Received an error parsing the flag: %v", err)
		}

		got := duration.String()

		if got != want {
			t.Errorf(
				"Unexpected duration parsed from the flag: want %s, got %s",
				want,
				got,
			)
		} else {
			t.Logf(
				"Expected duration parsed from the flag: got %s",
				got,
			)
		}
	}
}
