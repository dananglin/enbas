package main

import (
	"flag"
	"fmt"

	"codeflow.dananglin.me.uk/apollo/enbas/internal/client"
)

type blockCommand struct {
	*flag.FlagSet

	resourceType string
	accountName  string
	unblock      bool
}

func newBlockCommand(name, summary string, unblock bool) *blockCommand {
	command := blockCommand{
		FlagSet: flag.NewFlagSet(name, flag.ExitOnError),

		unblock: unblock,
	}

	command.StringVar(&command.resourceType, resourceTypeFlag, "", "specify the type of resource to block or unblock")
	command.StringVar(&command.accountName, accountNameFlag, "", "specify the account name in full (username@domain)")

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

func (c *blockCommand) blockAccount(gtsClient *client.Client) error {
	accountID, err := getAccountID(gtsClient, false, c.accountName)
	if err != nil {
		return fmt.Errorf("received an error while getting the account ID; %w", err)
	}

	if c.unblock {
		return c.unblockAccount(gtsClient, accountID)
	}

	if err := gtsClient.BlockAccount(accountID); err != nil {
		return fmt.Errorf("unable to block the account; %w", err)
	}

	fmt.Println("Successfully blocked the account.")

	return nil
}

func (c *blockCommand) unblockAccount(gtsClient *client.Client, accountID string) error {
	if err := gtsClient.UnblockAccount(accountID); err != nil {
		return fmt.Errorf("unable to unblock the account; %w", err)
	}

	fmt.Println("Successfully unblocked the account.")

	return nil
}
