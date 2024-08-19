package flag_test

import (
	"slices"
	"testing"

	internalFlag "codeflow.dananglin.me.uk/apollo/enbas/internal/flag"
)

func TestBoolPtrValue(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{
			input: "True",
			want:  true,
		},
		{
			input: "false",
			want:  false,
		},
	}

	value := internalFlag.NewBoolPtrValue()

	for _, test := range slices.All(tests) {
		if err := value.Set(test.input); err != nil {
			t.Fatalf(
				"Unable to parse %s as a BoolPtrValue: %v",
				test.input,
				err,
			)
		}

		got := *value.Value

		if got != test.want {
			t.Errorf(
				"Unexpected bool parsed from %s: want %t, got %t",
				test.input,
				test.want,
				got,
			)
		} else {
			t.Logf(
				"Expected bool parsed from %s: got %t",
				test.input,
				got,
			)
		}
	}
}

func TestNilBoolPtrValue(t *testing.T) {
	value := internalFlag.NewBoolPtrValue()
	want := "NOT SET"
	got := value.String()

	if got != want {
		t.Errorf("Unexpected string returned from the nil value; want %s, got %s", want, got)
	} else {
		t.Logf("Expected string returned from the nil value; got %s", got)
	}
}
