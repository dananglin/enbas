package command_test

import (
	"errors"
	"fmt"
	"reflect"
	"slices"
	"strings"
	"testing"

	"codeflow.dananglin.me.uk/apollo/enbas/internal/command"
)

func TestParseCommand(t *testing.T) {
	cases := []struct {
		args []string
		want command.Command
	}{
		{
			args: []string{"show", "version"},
			want: command.Command{
				Action:             "show",
				FocusedTarget:      "version",
				FocusedTargetFlags: []string{},
				Preposition:        "",
				RelatedTarget:      "",
				RelatedTargetFlags: []string{},
			},
		},
		{
			args: []string{"edit", "list", "--list-id", "HSRVaWEfhpiZl1ESzFHm", "-title", "Edited list title"},
			want: command.Command{
				Action:             "edit",
				FocusedTarget:      "list",
				FocusedTargetFlags: []string{"--list-id", "HSRVaWEfhpiZl1ESzFHm", "-title", "Edited list title"},
				Preposition:        "",
				RelatedTarget:      "",
				RelatedTargetFlags: []string{},
			},
		},
		{
			args: []string{"add", "note", "to", "account", "--content", "This is a private note about Joe", "--account-name", "@joe@gts.social.example"},
			want: command.Command{
				Action:             "add",
				FocusedTarget:      "note",
				FocusedTargetFlags: []string{},
				Preposition:        "to",
				RelatedTarget:      "account",
				RelatedTargetFlags: []string{"--content", "This is a private note about Joe", "--account-name", "@joe@gts.social.example"},
			},
		},
		{
			args: []string{"remove", "note", "from", "account", "--account-name", "@joe@gts.social.example"},
			want: command.Command{
				Action:             "remove",
				FocusedTarget:      "note",
				FocusedTargetFlags: []string{},
				Preposition:        "from",
				RelatedTarget:      "account",
				RelatedTargetFlags: []string{"--account-name", "@joe@gts.social.example"},
			},
		},
		{
			args: []string{"add", "status", "to", "bookmarks", "--status-id", "D40ED90BA3400B5A5CC3D22C971E697A"},
			want: command.Command{
				Action:             "add",
				FocusedTarget:      "status",
				FocusedTargetFlags: []string{},
				Preposition:        "to",
				RelatedTarget:      "bookmarks",
				RelatedTargetFlags: []string{"--status-id", "D40ED90BA3400B5A5CC3D22C971E697A"},
			},
		},
		{
			args: []string{"remove", "status", "from", "bookmarks", "--status-id", "AA4C6509A58E78DA4FFA9718AE9F193F"},
			want: command.Command{
				Action:             "remove",
				FocusedTarget:      "status",
				FocusedTargetFlags: []string{},
				Preposition:        "from",
				RelatedTarget:      "bookmarks",
				RelatedTargetFlags: []string{"--status-id", "AA4C6509A58E78DA4FFA9718AE9F193F"},
			},
		},
		{
			args: []string{"create", "access", "--url", "https://gts.example-host.example"},
			want: command.Command{
				Action:             "create",
				FocusedTarget:      "access",
				FocusedTargetFlags: []string{"--url", "https://gts.example-host.example"},
				Preposition:        "",
				RelatedTarget:      "",
				RelatedTargetFlags: []string{},
			},
		},
	}

	for ind, tc := range slices.All(cases) {
		t.Run(
			fmt.Sprintf("Test case: %d", ind+1),
			testParseCommand(tc.args, tc.want),
		)
	}
}

func testParseCommand(args []string, want command.Command) func(t *testing.T) {
	cmd := strings.Join(args, " ")

	return func(t *testing.T) {
		t.Parallel()

		got, err := command.Parse(args)
		if err != nil {
			t.Fatalf("FAILED: received parsing error for %q: %v", cmd, err)
		}

		if !reflect.DeepEqual(want, got) {
			t.Errorf(
				"FAILED: received unexpected result from parsing %q.\nwant: %+v\n got: %+v",
				cmd,
				want,
				got,
			)

			return
		}

		t.Logf(
			"Received expected result from parsing %q.\ngot: %+v",
			cmd,
			got,
		)

		// Ensure that the command is valid
		if err := got.Validate(); err != nil {
			t.Errorf("FAILED: validation failed for %q: %v", cmd, err)

			return
		}

		t.Logf("Validation passed for %q.", cmd)
	}
}

func TestCommandValidationErrors(t *testing.T) {
	cases := []struct {
		name      string
		cmd       command.Command
		wantError error
	}{
		{
			name: "No arguments provided (NoActionError expected)",
			cmd: command.Command{
				Action:             "",
				FocusedTarget:      "",
				FocusedTargetFlags: []string{},
				Preposition:        "",
				RelatedTarget:      "",
				RelatedTargetFlags: []string{},
			},
			wantError: command.NewNoActionError(),
		},
		{
			name: "No focused target specified",
			cmd: command.Command{
				Action:             "show",
				FocusedTarget:      "",
				FocusedTargetFlags: []string{},
				Preposition:        "",
				RelatedTarget:      "",
				RelatedTargetFlags: []string{},
			},
			wantError: command.NewNoFocusedTargetError("show"),
		},
		{
			name: "No related target specified if preposition is specified",
			cmd: command.Command{
				Action:             "add",
				FocusedTarget:      "note",
				FocusedTargetFlags: []string{"--content", "Test private note"},
				Preposition:        "to",
				RelatedTarget:      "",
				RelatedTargetFlags: []string{"--account-name", "@eric@gts.social.example"},
			},
			wantError: command.NewNoRelatedTargetError("add", "note", "to"),
		},
		{
			name: "-h flag detected",
			cmd: command.Command{
				Action:             "show",
				FocusedTarget:      "status",
				FocusedTargetFlags: []string{"-h"},
				Preposition:        "",
				RelatedTarget:      "",
				RelatedTargetFlags: []string{},
			},
			wantError: command.NewHelpFlagDetectedError("show", "status"),
		},
		{
			name: "-help flag detected",
			cmd: command.Command{
				Action:             "show",
				FocusedTarget:      "status",
				FocusedTargetFlags: []string{"-help", "me"},
				Preposition:        "",
				RelatedTarget:      "",
				RelatedTargetFlags: []string{},
			},
			wantError: command.NewHelpFlagDetectedError("show", "status"),
		},
		{
			name: "--h flag detected",
			cmd: command.Command{
				Action:             "add",
				FocusedTarget:      "account",
				FocusedTargetFlags: []string{"--account-name", "@bob@gts.social.example"},
				Preposition:        "to",
				RelatedTarget:      "list",
				RelatedTargetFlags: []string{"--h"},
			},
			wantError: command.NewHelpFlagDetectedError("add", "account"),
		},
		{
			name: "--help flag detected",
			cmd: command.Command{
				Action:             "add",
				FocusedTarget:      "account",
				FocusedTargetFlags: []string{"--account-name", "@bob@gts.social.example"},
				Preposition:        "to",
				RelatedTarget:      "list",
				RelatedTargetFlags: []string{"--help"},
			},
			wantError: command.NewHelpFlagDetectedError("add", "account"),
		},
	}

	for _, tc := range slices.All(cases) {
		t.Run(
			fmt.Sprintf("Test case: %s", tc.name),
			testCommandValidationError(tc.name, tc.cmd, tc.wantError),
		)
	}
}

func testCommandValidationError(testName string, cmd command.Command, wantErr error) func(t *testing.T) {
	return func(t *testing.T) {
		t.Parallel()

		err := cmd.Validate()
		if err == nil {
			t.Errorf(
				"FAILED test %q: No error received after expected validation failure",
				testName,
			)

			return
		}

		if !errors.Is(err, wantErr) {
			t.Errorf(
				"FAILED test %q: Unexpected error received after validation failure\nwant: %v\n got: %v",
				testName,
				wantErr,
				err,
			)

			return
		}

		t.Logf(
			"GOOD result from %q: Expected error received after validation failure\ngot: %v",
			testName,
			err,
		)
	}
}

func TestHelpCommand(t *testing.T) {
	help := command.HelpCommand()

	if err := help.Validate(); err != nil {
		t.Errorf("FAILED: Received an error after attempting to validate the help command: got %v", err)

		return
	}

	t.Log("PASSED: Validation passed for the help command")
}

func TestMissingPrepositionKeywordError(t *testing.T) {
	cases := []struct {
		name                   string
		args                   []string
		missingPrepositionWord string
	}{
		{
			name:                   "Preposition word missing between the focused and related targets",
			args:                   []string{"show", "followings", "account", "--account-name", "bob"},
			missingPrepositionWord: "from",
		},
		{
			name:                   "No more arguments after the focused target",
			args:                   []string{"add", "accounts"},
			missingPrepositionWord: "to",
		},
	}

	for _, tc := range slices.All(cases) {
		t.Run(
			tc.name,
			testMissingPrepositionKeywordError(
				tc.name,
				tc.args,
				tc.missingPrepositionWord,
			),
		)
	}
}

func testMissingPrepositionKeywordError(
	testName string,
	args []string,
	missingPrepositionWord string,
) func(t *testing.T) {
	return func(t *testing.T) {
		t.Parallel()

		wantErr := command.NewPrepositionKeywordMissingError(missingPrepositionWord)

		_, err := command.Parse(args)
		if err == nil {
			t.Errorf(
				"FAILED test %q: No error received after expected command parsing failure",
				testName,
			)

			return
		}

		if !errors.Is(err, wantErr) {
			t.Errorf(
				"Failed test %q: Unexpected error received after expected command parsing failure\nwant: %v\n got: %v",
				testName,
				wantErr,
				err,
			)

			return
		}

		t.Logf(
			"Expected error received after expected command parsing failure\ngot: %v",
			err,
		)
	}
}
