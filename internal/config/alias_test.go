package config_test

import (
	"errors"
	"os"
	"path/filepath"
	"slices"
	"testing"

	"codeflow.dananglin.me.uk/apollo/enbas/internal/config"
)

func TestAliases(t *testing.T) {
	testName := "Testing interactions with user defined aliases"
	t.Log(testName)

	configFilepath := filepath.Join("testdata", "config", t.Name()+".json")

	if err := config.SaveInitialConfigToFile(configFilepath); err != nil {
		t.Fatalf(
			"FAILED test %q: error saving the initial configuration to %q: %v",
			testName,
			configFilepath,
			err,
		)
	}

	defer func() {
		if err := os.Remove(configFilepath); err != nil {
			t.Fatalf(
				"received an error after attempting to clean up the test configuration file at %q: %v",
				configFilepath,
				err,
			)
		}
	}()

	testName = "Creating aliases"
	t.Run(testName, testCreateAlias(testName, configFilepath))

	testName = "Editing aliases"
	t.Run(testName, testEditAlias(testName, configFilepath))

	testName = "Renaming aliases"
	t.Run(testName, testRenameAlias(testName, configFilepath))

	testName = "Deleting aliases"
	t.Run(testName, testDeleteAlias(testName, configFilepath))
}

func testCreateAlias(testName, configFilepath string) func(t *testing.T) {
	return func(t *testing.T) {
		t.Logf("Creating some aliases")

		aliases := []struct {
			name string
			args string
		}{
			{
				name: "boost",
				args: "reblog status --status-id",
			},
			{
				name: "quick-toot",
				args: "create status --content-type plain",
			},
			{
				name: "fave",
				args: "favourite status --status-id",
			},
		}

		for _, alias := range slices.All(aliases) {
			if err := config.CreateAlias(configFilepath, alias.name, alias.args); err != nil {
				t.Fatalf(
					"FAILED test %q: received an error after attempting to create the alias %q: %v",
					testName,
					alias.name,
					err,
				)
			}

			t.Logf("Successfully created the alias %q", alias.name)
		}

		checkAlias(t, configFilepath, aliases[1].name, aliases[1].args)

		t.Logf("Ensuring an existing alias is not created again")

		err := config.CreateAlias(configFilepath, aliases[0].name, aliases[0].args)
		if err == nil {
			t.Errorf(
				"FAILED test %q: No error received after attempting to create an existing alias",
				testName,
			)

			return
		}

		wantErr := config.NewAliasAlreadyPresentError(aliases[0].name)

		if !errors.Is(err, wantErr) {
			t.Errorf(
				"FAILED test %q: Unexpected error received after attempting to create an existing alias\nwant: %v\n got: %v",
				testName,
				wantErr,
				err,
			)

			return
		}

		t.Logf(
			"GOOD result from %q: Expected error received after attempting to create an existing alias\ngot: %v",
			testName,
			err,
		)
	}
}

func testEditAlias(testName, configFilepath string) func(t *testing.T) {
	return func(t *testing.T) {
		t.Logf("Editing an existing alias")

		alias := "quick-toot"
		newArgs := "create status --content-type plain --visibility public --content"

		if err := config.EditAlias(
			configFilepath,
			alias,
			newArgs,
		); err != nil {
			t.Fatalf(
				"FAILED test %q: received an error after attempting to edit the alias %q: %v",
				testName,
				alias,
				err,
			)
		}

		t.Logf("Successfully edited the alias %q", alias)

		checkAlias(t, configFilepath, alias, newArgs)

		t.Log("Ensuring a non-existing alias is not edited")

		badAlias := "repost"
		badAliasArgs := "reblog status --status-id"

		err := config.EditAlias(configFilepath, badAlias, badAliasArgs)

		if err == nil {
			t.Errorf(
				"FAILED test %q: No error received after attempting to edit a non-existing alias",
				testName,
			)

			return
		}

		wantErr := config.NewAliasNotPresentError(badAlias)

		if !errors.Is(err, wantErr) {
			t.Errorf(
				"FAILED test %q: Unexpected error received after attempting to edit a non-existing alias\nwant: %v\n got: %v",
				testName,
				wantErr,
				err,
			)

			return
		}

		t.Logf(
			"GOOD result from %q: Expected error received after attempting to edit a non-existing alias\ngot: %v",
			testName,
			err,
		)
	}
}

func testRenameAlias(testName, configFilepath string) func(t *testing.T) {
	return func(t *testing.T) {
		t.Logf("Renaming an existing alias")

		oldName := "fave"
		newName := "like"
		newNameArgs := "favourite status --status-id"

		if err := config.RenameAlias(
			configFilepath,
			oldName,
			newName,
		); err != nil {
			t.Fatalf(
				"FAILED test %q: received an error after attempting to rename %q to %q: %v",
				testName,
				oldName,
				newName,
				err,
			)
		}

		t.Logf("Successfully renamed the alias from %q to %q", oldName, newName)

		checkAlias(t, configFilepath, newName, newNameArgs)
		checkAliasIsRemoved(t, configFilepath, oldName)

		t.Log("Ensuring a non-existing alias is not renamed")

		badOldName := "toot"
		badNewName := "post"

		err := config.RenameAlias(
			configFilepath,
			badOldName,
			badNewName,
		)

		if err == nil {
			t.Errorf(
				"FAILED test %q: No error received after attempting to rename a non-existing alias",
				testName,
			)

			return
		}

		wantErr := config.NewAliasNotPresentError(badOldName)

		if !errors.Is(err, wantErr) {
			t.Errorf(
				"FAILED test %q: Unexpected error received after attempting to rename a non-existing alias\nwant: %v\n got: %v",
				testName,
				wantErr,
				err,
			)

			return
		}

		t.Logf(
			"GOOD result from %q: Expected error received after attempting to rename a non-existing alias\ngot: %v",
			testName,
			err,
		)

		t.Logf("Ensuring that renaming an alias does not overwrite another")

		oldName = "boost"
		newName = "quick-toot"

		err = config.RenameAlias(
			configFilepath,
			oldName,
			newName,
		)

		if err == nil {
			t.Errorf(
				"FAILED test %q: No error received after attempting to rename an alias that overwrites another.",
				testName,
			)

			return
		}

		wantErr = config.NewAliasAlreadyPresentError(newName)

		if !errors.Is(err, wantErr) {
			t.Errorf(
				"FAILED test %q: Unexpected error received after attempting to rename an alias that overwrites another\nwant: %v\n got: %v",
				testName,
				wantErr,
				err,
			)

			return
		}

		t.Logf(
			"GOOD result from %q: Expected error received after attempting to rename an alias that overwrites another\ngot: %v",
			testName,
			err,
		)
	}
}

func testDeleteAlias(testName, configFilepath string) func(t *testing.T) {
	return func(t *testing.T) {
		alias := "quick-toot"

		if err := config.DeleteAlias(
			configFilepath,
			alias,
		); err != nil {
			t.Fatalf(
				"FAILED test %q: received an error after attempting to delete the alias %q: %v",
				testName,
				alias,
				err,
			)
		}

		t.Logf("Successfully deleted the alias %q", alias)

		checkAliasIsRemoved(t, configFilepath, alias)

		t.Logf("Ensuring a non-existing alias is not deleted")

		badAlias := "toot"

		err := config.DeleteAlias(configFilepath, badAlias)

		if err == nil {
			t.Errorf(
				"FAILED test %q: No error received after attempting to delete a non-existing alias",
				testName,
			)

			return
		}

		wantErr := config.NewAliasNotPresentError(badAlias)

		if !errors.Is(err, wantErr) {
			t.Errorf(
				"FAILED test %q: Unexpected error received after attempting to delete a non-existing alias\nwant: %v\n got: %v",
				testName,
				wantErr,
				err,
			)

			return
		}

		t.Logf(
			"GOOD result from %q: Expected error received after attempting to delete a non-existing alias\ngot: %v",
			testName,
			err,
		)
	}
}

func checkAlias(t *testing.T, configFilepath, alias, wantArgs string) {
	t.Helper()

	helperName := "check alias"

	cfg, err := config.NewConfigFromFile(configFilepath)
	if err != nil {
		t.Fatalf(
			"FAILED(%s): Received an error after attempting to load the configuration: %v",
			helperName,
			err,
		)
	}

	gotArgs, ok := cfg.Aliases[alias]
	if !ok {
		t.Errorf(
			"FAILED(%s): The alias %q does not appear in the configuration.",
			helperName,
			alias,
		)

		return
	}

	if wantArgs != gotArgs {
		t.Errorf(
			"BAD RESULTS(%s): Unexpected arguments received after checking alias %q\nwant: %q\n got: %q",
			helperName,
			alias,
			wantArgs,
			gotArgs,
		)

		return
	}

	t.Logf(
		"GOOD RESULTS(%s): Expected arguments received after checking alias %q\ngot: %q",
		helperName,
		alias,
		gotArgs,
	)
}

func checkAliasIsRemoved(t *testing.T, configFilepath, alias string) {
	t.Helper()

	helperName := "check alias is removed"

	cfg, err := config.NewConfigFromFile(configFilepath)
	if err != nil {
		t.Fatalf(
			"FAILED(%s): Received an error after attempting to load the configuration: %v",
			helperName,
			err,
		)
	}

	if _, exists := cfg.Aliases[alias]; exists {
		t.Errorf(
			"BAD RESULTS(%s): The alias %q is still present in the configuration file",
			helperName,
			alias,
		)
	}

	t.Logf(
		"GOOD RESULT(%s): The alias %q is not present in the configuration file",
		helperName,
		alias,
	)
}
