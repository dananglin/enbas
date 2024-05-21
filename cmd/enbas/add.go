package main

import (
	"flag"
	"fmt"

	"codeflow.dananglin.me.uk/apollo/enbas/internal/client"
)

type addCommand struct {
	*flag.FlagSet

	resourceType   string
	toResourceType string
	listID         string
	accountNames   accountNames
}

func newAddCommand(name, summary string) *addCommand {
	emptyArr := make([]string, 0, 3)

	command := addCommand{
		FlagSet:      flag.NewFlagSet(name, flag.ExitOnError),
		accountNames: accountNames(emptyArr),
	}

	command.StringVar(&command.resourceType, resourceTypeFlag, "", "specify the resource type to add (e.g. account, note)")
	command.StringVar(&command.toResourceType, addToFlag, "", "specify the target resource type to add to (e.g. list, account, etc)")
	command.StringVar(&command.listID, listIDFlag, "", "the ID of the list to add to")
	command.Var(&command.accountNames, accountNameFlag, "the name of the account to add to the resource")

	command.Usage = commandUsageFunc(name, summary, command.FlagSet)

	return &command
}

func (c *addCommand) Execute() error {
	if c.toResourceType == "" {
		return flagNotSetError{flagText: addToFlag}
	}

	funcMap := map[string]func(*client.Client) error{
		listResource: c.addToList,
	}

	doFunc, ok := funcMap[c.toResourceType]
	if !ok {
		return unsupportedResourceTypeError{resourceType: c.toResourceType}
	}

	gtsClient, err := client.NewClientFromConfig()
	if err != nil {
		return fmt.Errorf("unable to create the GoToSocial client; %w", err)
	}

	return doFunc(gtsClient)
}

func (c *addCommand) addToList(gtsClient *client.Client) error {
	funcMap := map[string]func(*client.Client) error{
		accountResource: c.addAccountsToList,
	}

	doFunc, ok := funcMap[c.resourceType]
	if !ok {
		return unsupportedResourceTypeError{resourceType: c.resourceType}
	}

	return doFunc(gtsClient)
}

func (c *addCommand) addAccountsToList(gtsClient *client.Client) error {
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

	if err := gtsClient.AddAccountsToList(c.listID, accountIDs); err != nil {
		return fmt.Errorf("unable to add the accounts to the list; %w", err)
	}

	fmt.Println("Successfully added the account(s) to the list.")

	return nil
}
