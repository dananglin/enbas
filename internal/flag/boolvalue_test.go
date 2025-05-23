package flag_test

import (
	"flag"
	"slices"
	"testing"

	internalFlag "codeflow.dananglin.me.uk/apollo/enbas/internal/flag"
)

func TestBoolValue(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{
			input: "True",
			want:  true,
		},
		{
			input: "true",
			want:  true,
		},
		{
			input: "1",
			want:  true,
		},
		{
			input: "False",
			want:  false,
		},
		{
			input: "false",
			want:  false,
		},
		{
			input: "0",
			want:  false,
		},
	}

	for _, test := range slices.All(tests) {
		args := []string{"--boolean-value=" + test.input}

		t.Run("Flag parsing test: "+test.input, testBoolPtrValueParsing(args, test.want))
	}
}

func testBoolPtrValueParsing(args []string, want bool) func(t *testing.T) {
	return func(t *testing.T) {
		flagset := flag.NewFlagSet("", flag.ExitOnError)
		boolVal := internalFlag.NewBoolValue(false)

		flagset.Var(&boolVal, "boolean-value", "Boolean value")

		if err := flagset.Parse(args); err != nil {
			t.Fatalf("Received an error parsing the flag: %v", err)
		}

		if !boolVal.IsSet() {
			t.Errorf("Unexpected result received from IsSet() method.\nwant: true\n got: %t", boolVal.IsSet())
		}

		got := boolVal.Value()

		if got != want {
			t.Errorf(
				"Unexpected boolean value found after parsing BoolPtrValue: want %t, got %t",
				want,
				got,
			)

			return
		}

		t.Logf(
			"Expected boolean value found after parsing BoolPtrValue: got %t",
			got,
		)
	}
}

func TestNotSetBoolPtrValue(t *testing.T) {
	flagset := flag.NewFlagSet("", flag.ExitOnError)
	boolVal := internalFlag.NewBoolValue(false)

	var otherVal string

	flagset.Var(&boolVal, "boolean-value", "Boolean value")
	flagset.StringVar(&otherVal, "other-value", "", "Another value")

	args := []string{"--other-value", "other-value"}

	if err := flagset.Parse(args); err != nil {
		t.Fatalf("Received an error parsing the flag: %v", err)
	}

	if boolVal.IsSet() {
		t.Errorf("Unexpected result received from IsSet() method.\nwant: false\n got: %t", boolVal.IsSet())
	}

	t.Logf("Expected result received from IsSet() method.\ngot: %t", boolVal.IsSet())
}
