package flag_test

import (
	"flag"
	"reflect"
	"testing"

	internalFlag "codeflow.dananglin.me.uk/apollo/enbas/internal/flag"
)

func TestMultiIntValue(t *testing.T) {
	flagset := flag.NewFlagSet("test", flag.ExitOnError)
	flagValue := internalFlag.NewMultiIntValue()

	if !flagValue.Empty() {
		t.Fatalf(
			"FAILED test %q: The initialised MultiIntValue is not empty",
			t.Name(),
		)
	}

	flagset.Var(&flagValue, "value", "Integer value")

	args := []string{
		"--value", "0",
		"--value", "1",
		"--value", "2",
		"--value", "3",
		"--value", "4",
	}

	if err := flagset.Parse(args); err != nil {
		t.Fatalf(
			"FAILED test %q: flag parsing error: %v",
			t.Name(),
			err,
		)
	}

	wantLength := 5
	if !flagValue.ExpectedLength(wantLength) {
		t.Errorf(
			"FAILED test %q: unexpected number of values found in the MultiIntValue\nwant: %d\n got: %d",
			t.Name(),
			wantLength,
			flagValue.Length(),
		)
	}

	t.Logf(
		"GOOD result from %q: expected number of values found in the parsed MultiIntValue\ngot: %d",
		t.Name(),
		flagValue.Length(),
	)

	want := []int{0, 1, 2, 3, 4}
	got := flagValue.Values()

	if !reflect.DeepEqual(got, want) {
		t.Errorf(
			"FAILED test %q: unexpected result found in parsed MultiIntValue\nwant: %v\n got %v",
			t.Name(),
			want,
			got,
		)

		return
	}

	t.Logf(
		"GOOD result from %q: expected result found in parsed MultiIntValue\ngot %v",
		t.Name(),
		got,
	)
}
