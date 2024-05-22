package main

import (
	"flag"
	"fmt"

	"codeflow.dananglin.me.uk/apollo/enbas/internal/client"
)

type followCommand struct {
	*flag.FlagSet

	topLevelFlags topLevelFlags
	resourceType  string
	accountName   string
	showReposts   bool
	notify        bool
	unfollow      bool
}

func newFollowCommand(tlf topLevelFlags, name, summary string, unfollow bool) *followCommand {
	command := followCommand{
		FlagSet:       flag.NewFlagSet(name, flag.ExitOnError),
		unfollow:      unfollow,
		topLevelFlags: tlf,
	}

	command.StringVar(&command.resourceType, resourceTypeFlag, "", "specify the type of resource to follow")
	command.StringVar(&command.accountName, accountNameFlag, "", "specify the account name in full (username@domain)")
	command.BoolVar(&command.showReposts, showRepostsFlag, true, "show reposts from the account you want to follow")
	command.BoolVar(&command.notify, notifyFlag, false, "get notifications when the account you want to follow posts a status")

	command.Usage = commandUsageFunc(name, summary, command.FlagSet)

	return &command
}

func (c *followCommand) Execute() error {
	funcMap := map[string]func(*client.Client) error{
		accountResource: c.followAccount,
	}

	doFunc, ok := funcMap[c.resourceType]
	if !ok {
		return unsupportedResourceTypeError{resourceType: c.resourceType}
	}

	gtsClient, err := client.NewClientFromConfig(c.topLevelFlags.configDir)
	if err != nil {
		return fmt.Errorf("unable to create the GoToSocial client; %w", err)
	}

	return doFunc(gtsClient)
}

func (c *followCommand) followAccount(gtsClient *client.Client) error {
	accountID, err := getAccountID(gtsClient, false, c.accountName, c.topLevelFlags.configDir)
	if err != nil {
		return fmt.Errorf("received an error while getting the account ID; %w", err)
	}

	if c.unfollow {
		return c.unfollowAccount(gtsClient, accountID)
	}

	if err := gtsClient.FollowAccount(accountID, c.showReposts, c.notify); err != nil {
		return fmt.Errorf("unable to follow the account; %w", err)
	}

	fmt.Println("The follow request was sent successfully.")

	return nil
}

func (c *followCommand) unfollowAccount(gtsClient *client.Client, accountID string) error {
	if err := gtsClient.UnfollowAccount(accountID); err != nil {
		return fmt.Errorf("unable to unfollow the account; %w", err)
	}

	fmt.Println("Successfully unfollowed the account.")

	return nil
}
