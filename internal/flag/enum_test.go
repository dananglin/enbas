package flag_test

import (
	"flag"
	"reflect"
	"slices"
	"testing"

	internalFlag "codeflow.dananglin.me.uk/apollo/enbas/internal/flag"
)

func TestEnumValue(t *testing.T) {
	cases := []struct {
		name          string
		enumFlagValue internalFlag.EnumValue
		input         string
		want          string
	}{
		{
			name: "The input is a valid value",
			enumFlagValue: internalFlag.NewEnumValue(
				[]string{
					"black",
					"red",
					"green",
					"yellow",
					"blue",
					"magenta",
					"cyan",
					"white",
				},
				"red",
			),
			input: "blue",
			want:  "blue",
		},
		{
			name: "The default value is used",
			enumFlagValue: internalFlag.NewEnumValue(
				[]string{
					"black",
					"red",
					"green",
					"yellow",
					"blue",
					"magenta",
					"cyan",
					"white",
				},
				"red",
			),
			input: "",
			want:  "red",
		},
		{
			name: "The default value is invalid",
			enumFlagValue: internalFlag.NewEnumValue(
				[]string{
					"black",
					"red",
					"green",
					"yellow",
					"blue",
					"magenta",
					"cyan",
					"white",
				},
				"orange",
			),
			input: "",
			want:  "black",
		},
		{
			name: "The default and input values are empty",
			enumFlagValue: internalFlag.NewEnumValue(
				[]string{
					"black",
					"red",
					"green",
					"yellow",
					"blue",
					"magenta",
					"cyan",
					"white",
				},
				"",
			),
			input: "",
			want:  "",
		},
	}

	for _, tc := range slices.All(cases) {
		var args []string

		if tc.input != "" {
			args = []string{"--colour", tc.input}
		} else {
			args = []string{}
		}

		t.Run(
			tc.name, testEnumValue(
				tc.name,
				tc.enumFlagValue,
				args,
				tc.want,
			),
		)
	}
}

func testEnumValue(
	testName string,
	flagValue internalFlag.EnumValue,
	args []string,
	want string,
) func(t *testing.T) {
	return func(t *testing.T) {
		t.Parallel()

		flagset := flag.NewFlagSet("test", flag.ExitOnError)
		flagset.Var(&flagValue, "colour", "Colour")

		if err := flagset.Parse(args); err != nil {
			t.Fatalf("FAILED test %q: received an error after attempting to parse the flag: %v",
				testName,
				err,
			)
		}

		got := flagValue.Value()

		if got != want {
			t.Errorf(
				"FAILED test %q: unexpected value received from parsed flag\nwant: %s\n got: %s",
				testName,
				want,
				got,
			)

			return
		}

		t.Logf(
			"GOOD result from %q: expected value received from parsed flag\ngot: %s",
			testName,
			got,
		)
	}
}

func TestMultiEnumValue(t *testing.T) {
	multiEnumFlagValue := internalFlag.NewMultiEnumValue(
		[]string{
			"black",
			"red",
			"green",
			"yellow",
			"blue",
			"magenta",
			"cyan",
			"white",
		},
	)

	cases := []struct {
		name   string
		inputs []string
		want   []string
	}{
		{
			name:   "All input values are valid",
			inputs: []string{"black", "blue", "white"},
			want:   []string{"black", "blue", "white"},
		},
		{
			name:   "No inputs provided",
			inputs: []string{},
			want:   []string{},
		},
	}

	for _, tc := range slices.All(cases) {
		args := make([]string, 0)

		for _, col := range slices.All(tc.inputs) {
			arg := []string{"--colour", col}
			args = append(args, arg...)
		}

		t.Run(
			tc.name, testMultiEnumValue(
				tc.name,
				multiEnumFlagValue,
				args,
				tc.want,
			),
		)
	}
}

func testMultiEnumValue(
	testName string,
	flagValue internalFlag.MultiEnumValue,
	args []string,
	want []string,
) func(t *testing.T) {
	return func(t *testing.T) {
		t.Parallel()

		flagset := flag.NewFlagSet("test", flag.ExitOnError)
		flagset.Var(&flagValue, "colour", "Color")

		if err := flagset.Parse(args); err != nil {
			t.Fatalf("FAILED test %q: received an error after attempting to parse the flag: %v",
				testName,
				err,
			)
		}

		got := flagValue.Values()

		if !reflect.DeepEqual(got, want) {
			t.Errorf(
				"FAILED test %q: unexpected values received from parsed flags\nwant: %v\n got: %v",
				testName,
				want,
				got,
			)

			return
		}

		t.Logf(
			"GOOD result from %q: expected values received from parsed flags\ngot: %v",
			testName,
			got,
		)
	}
}
