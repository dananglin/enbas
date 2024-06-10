// SPDX-FileCopyrightText: 2024 Dan Anglin <d.n.i.anglin@gmail.com>
//
// SPDX-License-Identifier: GPL-3.0-or-later

package executor

import (
	"flag"
	"fmt"

	"codeflow.dananglin.me.uk/apollo/enbas/internal/client"
)

type BlockOrUnblockExecutor struct {
	*flag.FlagSet

	topLevelFlags TopLevelFlags
	resourceType  string
	accountName   string
	command       string
}

func NewBlockOrUnblockExecutor(tlf TopLevelFlags, name, summary string) *BlockOrUnblockExecutor {
	blockExe := BlockOrUnblockExecutor{
		FlagSet: flag.NewFlagSet(name, flag.ExitOnError),

		topLevelFlags: tlf,
		command:       name,
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

	gtsClient, err := client.NewClientFromConfig(b.topLevelFlags.ConfigDir)
	if err != nil {
		return fmt.Errorf("unable to create the GoToSocial client: %w", err)
	}

	return doFunc(gtsClient)
}

func (b *BlockOrUnblockExecutor) blockOrUnblockAccount(gtsClient *client.Client) error {
	accountID, err := getAccountID(gtsClient, false, b.accountName, b.topLevelFlags.ConfigDir)
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

	fmt.Println("Successfully blocked the account.")

	return nil
}

func (b *BlockOrUnblockExecutor) unblockAccount(gtsClient *client.Client, accountID string) error {
	if err := gtsClient.UnblockAccount(accountID); err != nil {
		return fmt.Errorf("unable to unblock the account: %w", err)
	}

	fmt.Println("Successfully unblocked the account.")

	return nil
}
