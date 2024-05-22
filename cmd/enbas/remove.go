package main

import (
	"flag"
	"fmt"

	"codeflow.dananglin.me.uk/apollo/enbas/internal/client"
)

type removeCommand struct {
	*flag.FlagSet

	topLevelFlags topLevelFlags
	resourceType     string
	fromResourceType string
	listID           string
	accountNames     accountNames
}

func newRemoveCommand(tlf topLevelFlags, name, summary string) *removeCommand {
	emptyArr := make([]string, 0, 3)

	command := removeCommand{
		FlagSet:      flag.NewFlagSet(name, flag.ExitOnError),
		accountNames: accountNames(emptyArr),
		topLevelFlags: tlf,
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
		listResource:    c.removeFromList,
		accountResource: c.removeFromAccount,
	}

	doFunc, ok := funcMap[c.fromResourceType]
	if !ok {
		return unsupportedResourceTypeError{resourceType: c.fromResourceType}
	}

	gtsClient, err := client.NewClientFromConfig(c.topLevelFlags.configDir)
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
		return unsupportedRemoveOperationError{
			ResourceType:           c.resourceType,
			RemoveFromResourceType: c.fromResourceType,
		}
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

func (c *removeCommand) removeFromAccount(gtsClient *client.Client) error {
	funcMap := map[string]func(*client.Client) error{
		noteResource: c.removeNoteFromAccount,
	}

	doFunc, ok := funcMap[c.resourceType]
	if !ok {
		return unsupportedRemoveOperationError{
			ResourceType:           c.resourceType,
			RemoveFromResourceType: c.fromResourceType,
		}
	}

	return doFunc(gtsClient)
}

func (c *removeCommand) removeNoteFromAccount(gtsClient *client.Client) error {
	if len(c.accountNames) != 1 {
		return fmt.Errorf("unexpected number of accounts specified; want 1, got %d", len(c.accountNames))
	}

	accountID, err := getAccountID(gtsClient, false, c.accountNames[0], c.topLevelFlags.configDir)
	if err != nil {
		return fmt.Errorf("received an error while getting the account ID; %w", err)
	}

	if err := gtsClient.SetPrivateNote(accountID, ""); err != nil {
		return fmt.Errorf("unable to remove the private note from the account; %w", err)
	}

	fmt.Println("Successfully removed the private note from the account.")

	return nil
}
