// SPDX-FileCopyrightText: 2024 Dan Anglin <d.n.i.anglin@gmail.com>
//
// SPDX-License-Identifier: GPL-3.0-or-later

package executor

import (
	"flag"
	"fmt"

	"codeflow.dananglin.me.uk/apollo/enbas/internal/client"
)

type FollowExecutor struct {
	*flag.FlagSet

	topLevelFlags TopLevelFlags
	resourceType  string
	accountName   string
	showReposts   bool
	notify        bool
	unfollow      bool
}

func NewFollowExecutor(tlf TopLevelFlags, name, summary string, unfollow bool) *FollowExecutor {
	command := FollowExecutor{
		FlagSet:       flag.NewFlagSet(name, flag.ExitOnError),
		unfollow:      unfollow,
		topLevelFlags: tlf,
	}

	command.StringVar(&command.resourceType, flagType, "", "Specify the type of resource to follow")
	command.StringVar(&command.accountName, flagAccountName, "", "Specify the account name in full (username@domain)")
	command.BoolVar(&command.showReposts, flagShowReposts, true, "Show reposts from the account you want to follow")
	command.BoolVar(&command.notify, flagNotify, false, "Get notifications when the account you want to follow posts a status")

	command.Usage = commandUsageFunc(name, summary, command.FlagSet)

	return &command
}

func (c *FollowExecutor) Execute() error {
	funcMap := map[string]func(*client.Client) error{
		resourceAccount: c.followAccount,
	}

	doFunc, ok := funcMap[c.resourceType]
	if !ok {
		return UnsupportedTypeError{resourceType: c.resourceType}
	}

	gtsClient, err := client.NewClientFromConfig(c.topLevelFlags.ConfigDir)
	if err != nil {
		return fmt.Errorf("unable to create the GoToSocial client: %w", err)
	}

	return doFunc(gtsClient)
}

func (c *FollowExecutor) followAccount(gtsClient *client.Client) error {
	accountID, err := getAccountID(gtsClient, false, c.accountName, c.topLevelFlags.ConfigDir)
	if err != nil {
		return fmt.Errorf("received an error while getting the account ID: %w", err)
	}

	if c.unfollow {
		return c.unfollowAccount(gtsClient, accountID)
	}

	form := client.FollowAccountForm{
		AccountID:   accountID,
		ShowReposts: c.showReposts,
		Notify:      c.notify,
	}

	if err := gtsClient.FollowAccount(form); err != nil {
		return fmt.Errorf("unable to follow the account: %w", err)
	}

	fmt.Println("The follow request was sent successfully.")

	return nil
}

func (c *FollowExecutor) unfollowAccount(gtsClient *client.Client, accountID string) error {
	if err := gtsClient.UnfollowAccount(accountID); err != nil {
		return fmt.Errorf("unable to unfollow the account: %w", err)
	}

	fmt.Println("Successfully unfollowed the account.")

	return nil
}
