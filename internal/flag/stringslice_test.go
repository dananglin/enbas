package flag_test

import (
	"flag"
	"testing"

	internalFlag "codeflow.dananglin.me.uk/apollo/enbas/internal/flag"
)

func TestStringSliceValue(t *testing.T) {
	flagset := flag.NewFlagSet("test", flag.ExitOnError)
	stringSliceVal := internalFlag.NewStringSliceValue()

	if !stringSliceVal.Empty() {
		t.Fatalf("The initialised StringSliceValue is not empty")
	}

	flagset.Var(&stringSliceVal, "colour", "String value")

	args := []string{
		"--colour", "orange",
		"--colour", "blue",
		"--colour", "magenta",
		"--colour", "red",
		"--colour", "green",
		"--colour", "silver",
	}

	if err := flagset.Parse(args); err != nil {
		t.Fatalf("Received an error parsing the flag: %v", err)
	}

	wantLength := 6
	if !stringSliceVal.ExpectedLength(wantLength) {
		t.Fatalf(
			"Error: intSliceVal.ExpectedLength(%d) == false: actual length is %d",
			wantLength,
			len(stringSliceVal),
		)
	}

	want := "orange, blue, magenta, red, green, silver"
	got := stringSliceVal.String()

	if got != want {
		t.Errorf(
			"Unexpected result after parsing StringSliceValue: want %s, got %s",
			want,
			got,
		)
	} else {
		t.Logf(
			"Expected result after parsing StringSliceValue: got %s",
			got,
		)
	}
}
