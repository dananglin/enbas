// SPDX-FileCopyrightText: 2024 Dan Anglin <d.n.i.anglin@gmail.com>
//
// SPDX-License-Identifier: GPL-3.0-or-later

package executor

import (
	"flag"
	"fmt"

	"codeflow.dananglin.me.uk/apollo/enbas/internal/client"
)

type BlockExecutor struct {
	*flag.FlagSet

	topLevelFlags TopLevelFlags
	resourceType  string
	accountName   string
	unblock       bool
}

func NewBlockExecutor(tlf TopLevelFlags, name, summary string, unblock bool) *BlockExecutor {
	blockExe := BlockExecutor{
		FlagSet: flag.NewFlagSet(name, flag.ExitOnError),

		topLevelFlags: tlf,
		unblock:       unblock,
	}

	blockExe.StringVar(&blockExe.resourceType, flagType, "", "specify the type of resource to block or unblock")
	blockExe.StringVar(&blockExe.accountName, flagAccountName, "", "specify the account name in full (username@domain)")

	blockExe.Usage = commandUsageFunc(name, summary, blockExe.FlagSet)

	return &blockExe
}

func (b *BlockExecutor) Execute() error {
	funcMap := map[string]func(*client.Client) error{
		resourceAccount: b.blockAccount,
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

func (b *BlockExecutor) blockAccount(gtsClient *client.Client) error {
	accountID, err := getAccountID(gtsClient, false, b.accountName, b.topLevelFlags.ConfigDir)
	if err != nil {
		return fmt.Errorf("received an error while getting the account ID: %w", err)
	}

	if b.unblock {
		return b.unblockAccount(gtsClient, accountID)
	}

	if err := gtsClient.BlockAccount(accountID); err != nil {
		return fmt.Errorf("unable to block the account: %w", err)
	}

	fmt.Println("Successfully blocked the account.")

	return nil
}

func (b *BlockExecutor) unblockAccount(gtsClient *client.Client, accountID string) error {
	if err := gtsClient.UnblockAccount(accountID); err != nil {
		return fmt.Errorf("unable to unblock the account: %w", err)
	}

	fmt.Println("Successfully unblocked the account.")

	return nil
}
