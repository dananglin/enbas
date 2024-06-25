// SPDX-FileCopyrightText: 2024 Dan Anglin <d.n.i.anglin@gmail.com>
//
// SPDX-License-Identifier: GPL-3.0-or-later

package executor

import (
	"flag"
	"fmt"

	"codeflow.dananglin.me.uk/apollo/enbas/internal/config"
	"codeflow.dananglin.me.uk/apollo/enbas/internal/printer"
)

type WhoAmIExecutor struct {
	*flag.FlagSet

	printer *printer.Printer
	config  *config.Config
}

func NewWhoAmIExecutor(printer *printer.Printer, config *config.Config, name, summary string) *WhoAmIExecutor {
	whoExe := WhoAmIExecutor{
		FlagSet: flag.NewFlagSet(name, flag.ExitOnError),

		printer: printer,
		config:  config,
	}

	whoExe.Usage = commandUsageFunc(name, summary, whoExe.FlagSet)

	return &whoExe
}

func (c *WhoAmIExecutor) Execute() error {
	config, err := config.NewCredentialsConfigFromFile(c.config.CredentialsFile)
	if err != nil {
		return fmt.Errorf("unable to load the credential config: %w", err)
	}

	c.printer.PrintInfo("You are logged in as '" + config.CurrentAccount + "'.\n")

	return nil
}
