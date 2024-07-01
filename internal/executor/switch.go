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

type SwitchExecutor struct {
	*flag.FlagSet

	config         *config.Config
	printer        *printer.Printer
	toResourceType string
	accountName    string
}

func NewSwitchExecutor(printer *printer.Printer, config *config.Config, name, summary string) *SwitchExecutor {
	switchExe := SwitchExecutor{
		FlagSet: flag.NewFlagSet(name, flag.ExitOnError),

		config:  config,
		printer: printer,
	}

	switchExe.StringVar(&switchExe.toResourceType, flagTo, "", "The account to switch to")
	switchExe.StringVar(&switchExe.accountName, flagAccountName, "", "The name of the account to switch to")

	switchExe.Usage = commandUsageFunc(name, summary, switchExe.FlagSet)

	return &switchExe
}

func (s *SwitchExecutor) Execute() error {
	funcMap := map[string]func() error{
		resourceAccount: s.switchToAccount,
	}

	doFunc, ok := funcMap[s.toResourceType]
	if !ok {
		return UnsupportedTypeError{resourceType: s.toResourceType}
	}

	return doFunc()
}

func (s *SwitchExecutor) switchToAccount() error {
	if s.accountName == "" {
		return NoAccountSpecifiedError{}
	}

	if err := config.UpdateCurrentAccount(s.accountName, s.config.CredentialsFile); err != nil {
		return fmt.Errorf("unable to switch account to the account: %w", err)
	}

	s.printer.PrintSuccess("The current account is now set to '" + s.accountName + "'.")

	return nil
}
