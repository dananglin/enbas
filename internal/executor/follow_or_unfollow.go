// SPDX-FileCopyrightText: 2024 Dan Anglin <d.n.i.anglin@gmail.com>
//
// SPDX-License-Identifier: GPL-3.0-or-later

package executor

import (
	"flag"
	"fmt"

	"codeflow.dananglin.me.uk/apollo/enbas/internal/client"
)

type FollowOrUnfollowExecutor struct {
	*flag.FlagSet

	topLevelFlags TopLevelFlags
	resourceType  string
	accountName   string
	showReposts   bool
	notify        bool
	action        string
}

func NewFollowOrUnfollowExecutor(tlf TopLevelFlags, name, summary string) *FollowOrUnfollowExecutor {
	command := FollowOrUnfollowExecutor{
		FlagSet: flag.NewFlagSet(name, flag.ExitOnError),

		topLevelFlags: tlf,
		action:        name,
	}

	command.StringVar(&command.resourceType, flagType, "", "Specify the type of resource to follow")
	command.StringVar(&command.accountName, flagAccountName, "", "Specify the account name in full (username@domain)")
	command.BoolVar(&command.showReposts, flagShowReposts, true, "Show reposts from the account you want to follow")
	command.BoolVar(&command.notify, flagNotify, false, "Get notifications when the account you want to follow posts a status")

	command.Usage = commandUsageFunc(name, summary, command.FlagSet)

	return &command
}

func (f *FollowOrUnfollowExecutor) Execute() error {
	funcMap := map[string]func(*client.Client) error{
		resourceAccount: f.followOrUnfollowAccount,
	}

	doFunc, ok := funcMap[f.resourceType]
	if !ok {
		return UnsupportedTypeError{resourceType: f.resourceType}
	}

	gtsClient, err := client.NewClientFromConfig(f.topLevelFlags.ConfigDir)
	if err != nil {
		return fmt.Errorf("unable to create the GoToSocial client: %w", err)
	}

	return doFunc(gtsClient)
}

func (f *FollowOrUnfollowExecutor) followOrUnfollowAccount(gtsClient *client.Client) error {
	accountID, err := getAccountID(gtsClient, false, f.accountName, f.topLevelFlags.ConfigDir)
	if err != nil {
		return fmt.Errorf("received an error while getting the account ID: %w", err)
	}

	switch f.action {
	case CommandFollow:
		return f.followAccount(gtsClient, accountID)
	case CommandUnfollow:
		return f.unfollowAccount(gtsClient, accountID)
	default:
		return nil
	}
}

func (f *FollowOrUnfollowExecutor) followAccount(gtsClient *client.Client, accountID string) error {
	form := client.FollowAccountForm{
		AccountID:   accountID,
		ShowReposts: f.showReposts,
		Notify:      f.notify,
	}

	if err := gtsClient.FollowAccount(form); err != nil {
		return fmt.Errorf("unable to follow the account: %w", err)
	}

	fmt.Println("The follow request was sent successfully.")

	return nil
}

func (f *FollowOrUnfollowExecutor) unfollowAccount(gtsClient *client.Client, accountID string) error {
	if err := gtsClient.UnfollowAccount(accountID); err != nil {
		return fmt.Errorf("unable to unfollow the account: %w", err)
	}

	fmt.Println("Successfully unfollowed the account.")

	return nil
}
