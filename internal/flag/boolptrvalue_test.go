package flag_test

import (
	"flag"
	"slices"
	"testing"

	internalFlag "codeflow.dananglin.me.uk/apollo/enbas/internal/flag"
)

func TestBoolPtrValue(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{
			input: "True",
			want:  "true",
		},
		{
			input: "true",
			want:  "true",
		},
		{
			input: "1",
			want:  "true",
		},
		{
			input: "False",
			want:  "false",
		},
		{
			input: "false",
			want:  "false",
		},
		{
			input: "0",
			want:  "false",
		},
	}

	for _, test := range slices.All(tests) {
		args := []string{"--boolean-value=" + test.input}

		t.Run("Flag parsing test: "+test.input, testBoolPtrValueParsing(args, test.want))
	}
}

func testBoolPtrValueParsing(args []string, want string) func(t *testing.T) {
	return func(t *testing.T) {
		flagset := flag.NewFlagSet("test", flag.ExitOnError)
		boolVal := internalFlag.NewBoolPtrValue()

		flagset.Var(&boolVal, "boolean-value", "Boolean value")

		if err := flagset.Parse(args); err != nil {
			t.Fatalf("Received an error parsing the flag: %v", err)
		}

		got := boolVal.String()

		if got != want {
			t.Errorf(
				"Unexpected boolean value found after parsing BoolPtrValue: want %s, got %s",
				want,
				got,
			)
		} else {
			t.Logf(
				"Expected boolean value found after parsing BoolPtrValue: got %s",
				got,
			)
		}
	}
}

func TestNotSetBoolPtrValue(t *testing.T) {
	flagset := flag.NewFlagSet("test", flag.ExitOnError)
	boolVal := internalFlag.NewBoolPtrValue()

	var otherVal string

	flagset.Var(&boolVal, "boolean-value", "Boolean value")
	flagset.StringVar(&otherVal, "other-value", "", "Another value")

	args := []string{"--other-value", "other-value"}

	if err := flagset.Parse(args); err != nil {
		t.Fatalf("Received an error parsing the flag: %v", err)
	}

	want := "NOT SET"
	got := boolVal.String()

	if got != want {
		t.Errorf("Unexpected string returned from the nil value; want %s, got %s", want, got)
	} else {
		t.Logf("Expected string returned from the nil value; got %s", got)
	}
}
