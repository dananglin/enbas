package executor

import (
	"fmt"

	"codeflow.dananglin.me.uk/apollo/enbas/internal/config"
	"codeflow.dananglin.me.uk/apollo/enbas/internal/printer"
	"codeflow.dananglin.me.uk/apollo/enbas/internal/utilities"
)

func (i *InitExecutor) Execute() error {
	if err := utilities.EnsureDirectory(i.configDir); err != nil {
		return fmt.Errorf("unable to ensure that the configuration directory is present: %w", err)
	}

	printer.PrintSuccess(i.printSettings, "The configuration directory is present.")

	fileExists, err := config.FileExists(i.configDir)
	if err != nil {
		return fmt.Errorf("unable to check if the config file exists: %w", err)
	}

	if fileExists {
		printer.PrintInfo("The configuration file is already present in " + i.configDir + "\n")

		return nil
	}

	if err := config.SaveInitialConfigToFile(i.configDir); err != nil {
		return fmt.Errorf("unable to create a new configuration file in %s: %w", i.configDir, err)
	}

	printer.PrintSuccess(i.printSettings, "Successfully created a new configuration file in "+i.configDir)

	return nil
}
