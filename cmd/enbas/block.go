package main

import (
	"flag"
	"fmt"

	"codeflow.dananglin.me.uk/apollo/enbas/internal/client"
)

type blockCommand struct {
	*flag.FlagSet

	resourceType string
	accountID    string
	unblock      bool
}

func newBlockCommand(name, summary string, unblock bool) *blockCommand {
	command := blockCommand{
		FlagSet: flag.NewFlagSet(name, flag.ExitOnError),

		unblock: unblock,
	}

	command.StringVar(&command.resourceType, resourceTypeFlag, "", "specify the type of resource to block or unblock")
	command.StringVar(&command.accountID, accountIDFlag, "", "specify the ID of the account you want to block or unblock")

	command.Usage = commandUsageFunc(name, summary, command.FlagSet)

	return &command
}

func (c *blockCommand) Execute() error {
	funcMap := map[string]func(*client.Client) error{
		accountResource: c.blockAccount,
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

func (c *blockCommand) blockAccount(gts *client.Client) error {
	if c.accountID == "" {
		return flagNotSetError{flagText: accountIDFlag}
	}

	if c.unblock {
		return c.unblockAccount(gts)
	}

	if err := gts.BlockAccount(c.accountID); err != nil {
		return fmt.Errorf("unable to block the account; %w", err)
	}

	fmt.Println("Successfully blocked the account.")

	return nil
}

func (c *blockCommand) unblockAccount(gts *client.Client) error {
	if err := gts.UnblockAccount(c.accountID); err != nil {
		return fmt.Errorf("unable to unblock the account; %w", err)
	}

	fmt.Println("Successfully unblocked the account.")

	return nil
}
