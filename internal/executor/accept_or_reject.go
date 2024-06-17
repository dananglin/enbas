// SPDX-FileCopyrightText: 2024 Dan Anglin <d.n.i.anglin@gmail.com>
//
// SPDX-License-Identifier: GPL-3.0-or-later

package executor

import (
	"flag"
	"fmt"

	"codeflow.dananglin.me.uk/apollo/enbas/internal/client"
	"codeflow.dananglin.me.uk/apollo/enbas/internal/printer"
)

type AcceptOrRejectExecutor struct {
	*flag.FlagSet

	printer      *printer.Printer
	configDir    string
	resourceType string
	accountName  string
	command      string
}

func NewAcceptOrRejectExecutor(enbasPrinter *printer.Printer, configDir, name, summary string) *AcceptOrRejectExecutor {
	acceptExe := AcceptOrRejectExecutor{
		FlagSet: flag.NewFlagSet(name, flag.ExitOnError),

		printer:   enbasPrinter,
		configDir: configDir,
		command:   name,
	}

	acceptExe.StringVar(&acceptExe.resourceType, flagType, "", "Specify the type of resource to accept or reject")
	acceptExe.StringVar(&acceptExe.accountName, flagAccountName, "", "Specify the account name in full (username@domain)")

	acceptExe.Usage = commandUsageFunc(name, summary, acceptExe.FlagSet)

	return &acceptExe
}

func (a *AcceptOrRejectExecutor) Execute() error {
	funcMap := map[string]func(*client.Client) error{
		resourceFollowRequest: a.acceptOrRejectFollowRequest,
	}

	doFunc, ok := funcMap[a.resourceType]
	if !ok {
		return UnsupportedTypeError{resourceType: a.resourceType}
	}

	gtsClient, err := client.NewClientFromConfig(a.configDir)
	if err != nil {
		return fmt.Errorf("unable to create the GoToSocial client: %w", err)
	}

	return doFunc(gtsClient)
}

func (a *AcceptOrRejectExecutor) acceptOrRejectFollowRequest(gtsClient *client.Client) error {
	accountID, err := getAccountID(gtsClient, false, a.accountName, a.configDir)
	if err != nil {
		return fmt.Errorf("received an error while getting the account ID: %w", err)
	}

	switch a.command {
	case CommandAccept:
		return a.acceptFollowRequest(gtsClient, accountID)
	case CommandReject:
		return a.rejectFollowRequest(gtsClient, accountID)
	default:
		return nil
	}
}

func (a *AcceptOrRejectExecutor) acceptFollowRequest(gtsClient *client.Client, accountID string) error {
	if err := gtsClient.AcceptFollowRequest(accountID); err != nil {
		return fmt.Errorf("unable to accept the follow request: %w", err)
	}

	a.printer.PrintSuccess("Successfully accepted the follow request.")

	return nil
}

func (a *AcceptOrRejectExecutor) rejectFollowRequest(gtsClient *client.Client, accountID string) error {
	if err := gtsClient.RejectFollowRequest(accountID); err != nil {
		return fmt.Errorf("unable to reject the follow request: %w", err)
	}

	a.printer.PrintSuccess("Successfully rejected the follow request.")

	return nil
}
