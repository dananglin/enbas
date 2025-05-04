package flag_test

import (
	"flag"
	"reflect"
	"testing"

	internalFlag "codeflow.dananglin.me.uk/apollo/enbas/internal/flag"
)

func TestMultiStringValue(t *testing.T) {
	flagset := flag.NewFlagSet("test", flag.ExitOnError)
	flagValue := internalFlag.NewMultiStringValue()

	if !flagValue.Empty() {
		t.Fatalf(
			"FAILED test %q: the initialised MultiStringValue is not empty",
			t.Name(),
		)
	}

	flagset.Var(&flagValue, "colour", "String value")

	args := []string{
		"--colour", "orange",
		"--colour", "blue",
		"--colour", "magenta",
		"--colour", "red",
		"--colour", "green",
		"--colour", "silver",
	}

	if err := flagset.Parse(args); err != nil {
		t.Fatalf(
			"FAILED test %q: flag parsing error: %v",
			t.Name(),
			err,
		)
	}

	wantLength := 6
	if !flagValue.ExpectedLength(wantLength) {
		t.Errorf(
			"FAILED test %q: unexpected number of values found in the MultiStringValue\nwant: %d\n got: %d",
			t.Name(),
			wantLength,
			flagValue.Length(),
		)
	} else {
		t.Logf(
			"GOOD result from %q: expected number of values found in the parsed MultiStringValue\ngot: %d",
			t.Name(),
			flagValue.Length(),
		)
	}

	want := []string{"orange", "blue", "magenta", "red", "green", "silver"}
	got := flagValue.Values()

	if !reflect.DeepEqual(got, want) {
		t.Errorf(
			"FAILED test %q: unexpected result found in parsed MultiStringValue\nwant: %v\n got: %v",
			t.Name(),
			want,
			got,
		)

		return
	}

	t.Logf(
		"GOOD result from %q: expected result found in parsed MultiStringValue\ngot %v",
		t.Name(),
		got,
	)
}
