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

type BlockOrUnblockExecutor struct {
	*flag.FlagSet

	printer      *printer.Printer
	configDir    string
	resourceType string
	accountName  string
	command      string
}

func NewBlockOrUnblockExecutor(printer *printer.Printer, configDir, name, summary string) *BlockOrUnblockExecutor {
	blockExe := BlockOrUnblockExecutor{
		FlagSet: flag.NewFlagSet(name, flag.ExitOnError),

		printer:   printer,
		configDir: configDir,
		command:   name,
	}

	blockExe.StringVar(&blockExe.resourceType, flagType, "", "Specify the type of resource to block or unblock")
	blockExe.StringVar(&blockExe.accountName, flagAccountName, "", "Specify the account name in full (username@domain)")

	blockExe.Usage = commandUsageFunc(name, summary, blockExe.FlagSet)

	return &blockExe
}

func (b *BlockOrUnblockExecutor) Execute() error {
	funcMap := map[string]func(*client.Client) error{
		resourceAccount: b.blockOrUnblockAccount,
	}

	doFunc, ok := funcMap[b.resourceType]
	if !ok {
		return UnsupportedTypeError{resourceType: b.resourceType}
	}

	gtsClient, err := client.NewClientFromConfig(b.configDir)
	if err != nil {
		return fmt.Errorf("unable to create the GoToSocial client: %w", err)
	}

	return doFunc(gtsClient)
}

func (b *BlockOrUnblockExecutor) blockOrUnblockAccount(gtsClient *client.Client) error {
	if b.accountName == "" {
		return FlagNotSetError{flagText: flagAccountName}
	}

	accountID, err := getAccountID(gtsClient, false, b.accountName, b.configDir)
	if err != nil {
		return fmt.Errorf("received an error while getting the account ID: %w", err)
	}

	switch b.command {
	case CommandBlock:
		return b.blockAccount(gtsClient, accountID)
	case CommandUnblock:
		return b.unblockAccount(gtsClient, accountID)
	default:
		return nil
	}
}

func (b *BlockOrUnblockExecutor) blockAccount(gtsClient *client.Client, accountID string) error {
	if err := gtsClient.BlockAccount(accountID); err != nil {
		return fmt.Errorf("unable to block the account: %w", err)
	}

	b.printer.PrintSuccess("Successfully blocked the account.")

	return nil
}

func (b *BlockOrUnblockExecutor) unblockAccount(gtsClient *client.Client, accountID string) error {
	if err := gtsClient.UnblockAccount(accountID); err != nil {
		return fmt.Errorf("unable to unblock the account: %w", err)
	}

	b.printer.PrintSuccess("Successfully unblocked the account.")

	return nil
}
