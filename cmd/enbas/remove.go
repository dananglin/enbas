package main

import (
	"flag"
	"fmt"

	"codeflow.dananglin.me.uk/apollo/enbas/internal/client"
)

type removeCommand struct {
	*flag.FlagSet

	resourceType     string
	fromResourceType string
	listID           string
	accountNames     accountNames
}

func newRemoveCommand(name, summary string) *removeCommand {
	emptyArr := make([]string, 0, 3)

	command := removeCommand{
		FlagSet:      flag.NewFlagSet(name, flag.ExitOnError),
		accountNames: accountNames(emptyArr),
	}

	command.StringVar(&command.resourceType, resourceTypeFlag, "", "specify the resource type to remove (e.g. account, note)")
	command.StringVar(&command.fromResourceType, removeFromFlag, "", "specify the resource type to remove from (e.g. list, account, etc)")
	command.StringVar(&command.listID, listIDFlag, "", "the ID of the list to remove from")
	command.Var(&command.accountNames, accountNameFlag, "the name of the account to remove from the resource")

	command.Usage = commandUsageFunc(name, summary, command.FlagSet)

	return &command
}

func (c *removeCommand) Execute() error {
	if c.fromResourceType == "" {
		return flagNotSetError{flagText: removeFromFlag}
	}

	funcMap := map[string]func(*client.Client) error{
		listResource: c.removeFromList,
	}

	doFunc, ok := funcMap[c.fromResourceType]
	if !ok {
		return unsupportedResourceTypeError{resourceType: c.fromResourceType}
	}

	gtsClient, err := client.NewClientFromConfig()
	if err != nil {
		return fmt.Errorf("unable to create the GoToSocial client; %w", err)
	}

	return doFunc(gtsClient)
}

func (c *removeCommand) removeFromList(gtsClient *client.Client) error {
	funcMap := map[string]func(*client.Client) error{
		accountResource: c.removeAccountsFromList,
	}

	doFunc, ok := funcMap[c.resourceType]
	if !ok {
		return unsupportedResourceTypeError{resourceType: c.resourceType}
	}

	return doFunc(gtsClient)
}

func (c *removeCommand) removeAccountsFromList(gtsClient *client.Client) error {
	if c.listID == "" {
		return flagNotSetError{flagText: listIDFlag}
	}

	if len(c.accountNames) == 0 {
		return noAccountSpecifiedError{}
	}

	accountIDs := make([]string, len(c.accountNames))

	for i := range c.accountNames {
		accountID, err := getTheirAccountID(gtsClient, c.accountNames[i])
		if err != nil {
			return fmt.Errorf("unable to get the account ID for %s, %w", c.accountNames[i], err)
		}

		accountIDs[i] = accountID
	}

	if err := gtsClient.RemoveAccountsFromList(c.listID, accountIDs); err != nil {
		return fmt.Errorf("unable to remove the accounts from the list; %w", err)
	}

	fmt.Println("Successfully removed the account(s) from the list.")

	return nil
}
