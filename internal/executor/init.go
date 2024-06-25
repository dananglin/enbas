// SPDX-FileCopyrightText: 2024 Dan Anglin <d.n.i.anglin@gmail.com>
//
// SPDX-License-Identifier: GPL-3.0-or-later

package executor

import (
	"flag"
	"fmt"

	"codeflow.dananglin.me.uk/apollo/enbas/internal/config"
	"codeflow.dananglin.me.uk/apollo/enbas/internal/printer"
	"codeflow.dananglin.me.uk/apollo/enbas/internal/utilities"
)

type InitExecutor struct {
	*flag.FlagSet

	printer   *printer.Printer
	configDir string
}

func NewInitExecutor(printer *printer.Printer, configDir, name, summary string) *InitExecutor {
	initExe := InitExecutor{
		FlagSet: flag.NewFlagSet(name, flag.ExitOnError),

		printer:   printer,
		configDir: configDir,
	}

	initExe.Usage = commandUsageFunc(name, summary, initExe.FlagSet)

	return &initExe
}

func (i *InitExecutor) Execute() error {
	if err := utilities.EnsureDirectory(i.configDir); err != nil {
		return fmt.Errorf("unable to ensure that the configuration directory is present: %w", err)
	}

	i.printer.PrintSuccess("The configuration directory is present.")

	fileExists, err := config.FileExists(i.configDir)
	if err != nil {
		return fmt.Errorf("unable to check if the config file exists: %w", err)
	}

	if fileExists {
		i.printer.PrintInfo("The configuration file is already present in " + i.configDir + "\n")

		return nil
	}

	if err := config.SaveDefaultConfigToFile(i.configDir); err != nil {
		return fmt.Errorf("unable to create a new configuration file in %s: %w", i.configDir, err)
	}

	i.printer.PrintSuccess("Successfully created a new configuration file in " + i.configDir)

	return nil
}
