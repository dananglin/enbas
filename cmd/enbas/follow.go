package main

import (
	"flag"
	"fmt"

	"codeflow.dananglin.me.uk/apollo/enbas/internal/client"
)

type followCommand struct {
	*flag.FlagSet

	resourceType string
	accountID    string
	showReposts  bool
	notify       bool
	unfollow     bool
}

func newFollowCommand(name, summary string, unfollow bool) *followCommand {
	command := followCommand{
		FlagSet: flag.NewFlagSet(name, flag.ExitOnError),

		unfollow: unfollow,
	}

	command.StringVar(&command.resourceType, resourceTypeFlag, "", "specify the type of resource to follow")
	command.StringVar(&command.accountID, accountIDFlag, "", "specify the ID of the account you want to follow")
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

	gtsClient, err := client.NewClientFromConfig()
	if err != nil {
		return fmt.Errorf("unable to create the GoToSocial client; %w", err)
	}

	return doFunc(gtsClient)
}

func (c *followCommand) followAccount(gts *client.Client) error {
	if c.accountID == "" {
		return flagNotSetError{flagText: accountIDFlag}
	}

	if c.unfollow {
		return c.unfollowAccount(gts)
	}

	if err := gts.FollowAccount(c.accountID, c.showReposts, c.notify); err != nil {
		return fmt.Errorf("unable to follow the account; %w", err)
	}

	fmt.Println("The follow request was sent successfully.")

	return nil
}

func (c *followCommand) unfollowAccount(gts *client.Client) error {
	if err := gts.UnfollowAccount(c.accountID); err != nil {
		return fmt.Errorf("unable to unfollow the account; %w", err)
	}

	fmt.Println("Successfully unfollowed the account.")

	return nil
}
