package command_test

import (
	"errors"
	"reflect"
	"slices"
	"strings"
	"testing"

	"codeflow.dananglin.me.uk/apollo/enbas/internal/command"
)

func TestExtractArgsFromAlias(t *testing.T) {
	testUserDefinedAliases := map[string]string{
		"boost":         "reblog status --status-id",
		"home-timeline": "show timeline --timeline-category home",
		"my-followers":  "show followers from account --my-account",
	}

	cases := []struct {
		testName string
		args     []string
		want     []string
	}{
		{
			testName: "The alias is a built-in action word",
			args:     []string{"show", "status", "--status-id", "01JFERJPVNLUD9VDSVAIC3BEVZQSC1LJ"},
			want:     []string{"show", "status", "--status-id", "01JFERJPVNLUD9VDSVAIC3BEVZQSC1LJ"},
		},
		{
			testName: "The alias is a built-in alias",
			args:     []string{"whoami"},
			want:     []string{"verify", "access"},
		},
		{
			testName: "The alias is a user-defined alias",
			args:     []string{"boost", "01FILLVHHONFJLCJRROAARMSPH1JUV5L"},
			want:     []string{"reblog", "status", "--status-id", "01FILLVHHONFJLCJRROAARMSPH1JUV5L"},
		},
		{
			testName: "The alias is unknown",
			args:     []string{"turn", "status", "upside", "down", "--status-id", "01KDMQNWAKXVE1JUMVJN8PSBPO6A2VRN"},
			want:     []string{"turn", "status", "upside", "down", "--status-id", "01KDMQNWAKXVE1JUMVJN8PSBPO6A2VRN"},
		},
	}

	for _, tc := range slices.All(cases) {
		t.Run(
			tc.testName,
			testExtractArgsFromAlias(
				tc.testName,
				testUserDefinedAliases,
				tc.args,
				tc.want,
			),
		)
	}
}

func testExtractArgsFromAlias(
	testName string,
	userDeinedAliases map[string]string,
	cmd []string,
	want []string,
) func(t *testing.T) {
	return func(t *testing.T) {
		t.Parallel()

		got, err := command.ExtractArgsFromAlias(cmd, userDeinedAliases)
		if err != nil {
			t.Fatalf(
				"FAILED test %s: received an error after parsing the alias from command %q: %v",
				testName,
				strings.Join(cmd, " "),
				err,
			)
		}

		if !reflect.DeepEqual(want, got) {
			t.Errorf(
				"FAILED test %s: received unexpected result from parsing %q.\nwant: %+v\n got: %+v",
				testName,
				strings.Join(cmd, " "),
				want,
				got,
			)

			return
		}

		t.Logf(
			"GOOD result from %q: received expected result from parsing %q.\ngot: %+v",
			testName,
			strings.Join(cmd, " "),
			got,
		)
	}
}

func TestInvalidAliasErrorHandling(t *testing.T) {
	cases := []struct {
		name string
		args []string
	}{
		{
			name: "The alias has upper case letters",
			args: []string{"BoOsT", "01R0OEDWVCKRAXEEPLNBDSLB4GUTRKCJ"},
		},
		{
			name: "The alias is too short",
			args: []string{"bo", "01R0OEDWVCKRAXEEPLNBDSLB4GUTRKCJ"},
		},
		{
			name: "The alias contains an unsupported character",
			args: []string{"boost_status", "01R0OEDWVCKRAXEEPLNBDSLB4GUTRKCJ"},
		},
		{
			name: "The alias contains no characters",
			args: []string{"", "--status-id", "01R0OEDWVCKRAXEEPLNBDSLB4GUTRKCJ"},
		},
		{
			name: "The alias starts with a hyphen",
			args: []string{"-boost-status", "01R0OEDWVCKRAXEEPLNBDSLB4GUTRKCJ"},
		},
		{
			name: "The alias ends with a hyphen",
			args: []string{"boost-status-", "01R0OEDWVCKRAXEEPLNBDSLB4GUTRKCJ"},
		},
	}

	for _, tc := range slices.All(cases) {
		t.Run(
			tc.name,
			testInvalidAliasErrorHandling(
				tc.name,
				tc.args,
			),
		)
	}
}

func testInvalidAliasErrorHandling(
	testName string,
	args []string,
) func(t *testing.T) {
	return func(t *testing.T) {
		t.Parallel()

		_, err := command.ExtractArgsFromAlias(
			args,
			make(map[string]string),
		)

		if err == nil {
			t.Errorf(
				"FAILED test %q: No error received after parsing invalid alias",
				testName,
			)

			return
		}

		wantErr := command.NewInvalidAliasError(args[0])

		if !errors.Is(err, wantErr) {
			t.Errorf(
				"FAILED test %q: Unexpected error received after parsing invalid alias from %q\nwant: %v\n got: %v",
				testName,
				strings.Join(args, " "),
				wantErr,
				err,
			)

			return
		}

		t.Logf(
			"GOOD result from %q: Expected error received after parsing invalid alias\ngot: %v",
			testName,
			err,
		)
	}
}
