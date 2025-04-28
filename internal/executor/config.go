package executor

import (
	"fmt"

	"codeflow.dananglin.me.uk/apollo/enbas/internal/cli"
	"codeflow.dananglin.me.uk/apollo/enbas/internal/command"
	"codeflow.dananglin.me.uk/apollo/enbas/internal/config"
	"codeflow.dananglin.me.uk/apollo/enbas/internal/printer"
)

// configFunc is the function for the config target which interacts with the
// application's config file.
func configFunc(
	cfg config.Config,
	printSettings printer.Settings,
	cmd command.Command,
) error {
	switch cmd.Action {
	case cli.ActionCreate:
		return configCreate(printSettings, cfg.Path)
	default:
		return unsupportedActionError{action: cmd.Action, target: cli.TargetConfig}
	}
}

func configCreate(
	printSettings printer.Settings,
	configPath string,
) error {
	if err := config.EnsureParentDir(configPath); err != nil {
		return fmt.Errorf("error checking the existence of the configuration directory: %w", err)
	}

	printer.PrintSuccess(printSettings, "The configuration directory is present.")

	fileExists, err := config.FileExists(configPath)
	if err != nil {
		return fmt.Errorf("error checking the existence of the configuration file: %w", err)
	}

	if fileExists {
		printer.PrintInfo("A file or directory is already present at '" + configPath + "'.\n")

		return nil
	}

	if err := config.SaveInitialConfigToFile(configPath); err != nil {
		return fmt.Errorf("error creating the new configuration file: %w", err)
	}

	printer.PrintSuccess(
		printSettings,
		"Successfully created a new configuration file at '"+configPath+"'.",
	)

	return nil
}
