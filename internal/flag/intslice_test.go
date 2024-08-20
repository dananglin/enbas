package flag_test

import (
	"flag"
	"testing"

	internalFlag "codeflow.dananglin.me.uk/apollo/enbas/internal/flag"
)

func TestIntSliceValue(t *testing.T) {
	flagset := flag.NewFlagSet("test", flag.ExitOnError)
	intSliceVal := internalFlag.NewIntSliceValue()

	if !intSliceVal.Empty() {
		t.Fatalf("The initialised IntSliceValue is not empty")
	}

	flagset.Var(&intSliceVal, "int-value", "Integer value")

	args := []string{
		"--int-value", "0",
		"--int-value", "1",
		"--int-value", "2",
		"--int-value", "3",
		"--int-value", "4",
	}

	if err := flagset.Parse(args); err != nil {
		t.Fatalf("Received an error parsing the flag: %v", err)
	}

	wantLength := 5
	if !intSliceVal.ExpectedLength(wantLength) {
		t.Fatalf(
			"Error: intSliceVal.ExpectedLength(%d) == false: actual length is %d",
			wantLength,
			len(intSliceVal),
		)
	}

	want := "0, 1, 2, 3, 4"
	got := intSliceVal.String()

	if got != want {
		t.Errorf(
			"Unexpected result after parsing IntSliceValue: want %s, got %s",
			want,
			got,
		)
	} else {
		t.Logf(
			"Expected result after parsing IntSliceValue: got %s",
			got,
		)
	}
}
