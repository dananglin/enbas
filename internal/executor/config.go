package executor

import (
	"fmt"

	"codeflow.dananglin.me.uk/apollo/enbas/internal/cli"
	"codeflow.dananglin.me.uk/apollo/enbas/internal/command"
	"codeflow.dananglin.me.uk/apollo/enbas/internal/config"
	"codeflow.dananglin.me.uk/apollo/enbas/internal/printer"
	"codeflow.dananglin.me.uk/apollo/enbas/internal/utilities"
)

// configFunc is the function for the config target which interacts with the
// application's config file.
func configFunc(
	opts topLevelOpts,
	cmd command.Command,
) error {
	// Create the print settings
	printSettings := printer.NewSettings(opts.noColor, "", 0)

	switch cmd.Action {
	case cli.ActionCreate:
		return configCreate(printSettings, opts.configDir)
	default:
		return unsupportedActionError{action: cmd.Action, target: cli.TargetConfig}
	}
}

func configCreate(
	printSettings printer.Settings,
	configDir string,
) error {
	if err := utilities.EnsureDirectory(configDir); err != nil {
		return fmt.Errorf("error checking the existence of the configuration directory: %w", err)
	}

	printer.PrintSuccess(printSettings, "The configuration directory is present.")

	fileExists, err := config.FileExists(configDir)
	if err != nil {
		return fmt.Errorf("error checking the existence of the configuration file: %w", err)
	}

	if fileExists {
		printer.PrintInfo("The configuration file is already present in " + configDir + "\n")

		return nil
	}

	if err := config.SaveInitialConfigToFile(configDir); err != nil {
		return fmt.Errorf("error creating the new configuration file: %w", err)
	}

	printer.PrintSuccess(
		printSettings,
		"Successfully created a new configuration file in "+configDir,
	)

	return nil
}
