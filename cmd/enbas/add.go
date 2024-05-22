package main

import (
	"errors"
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
	content        string
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
	command.StringVar(&command.content, contentFlag, "", "the content of the note")

	command.Usage = commandUsageFunc(name, summary, command.FlagSet)

	return &command
}

func (c *addCommand) Execute() error {
	if c.toResourceType == "" {
		return flagNotSetError{flagText: addToFlag}
	}

	funcMap := map[string]func(*client.Client) error{
		listResource:    c.addToList,
		accountResource: c.addToAccount,
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
		return unsupportedAddOperationError{
			ResourceType:      c.resourceType,
			AddToResourceType: c.toResourceType,
		}
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

func (c *addCommand) addToAccount(gtsClient *client.Client) error {
	funcMap := map[string]func(*client.Client) error{
		noteResource: c.addNoteToAccount,
	}

	doFunc, ok := funcMap[c.resourceType]
	if !ok {
		return unsupportedAddOperationError{
			ResourceType:      c.resourceType,
			AddToResourceType: c.toResourceType,
		}
	}

	return doFunc(gtsClient)
}

func (c *addCommand) addNoteToAccount(gtsClient *client.Client) error {
	if len(c.accountNames) != 1 {
		return fmt.Errorf("unexpected number of accounts specified; want 1, got %d", len(c.accountNames))
	}

	accountID, err := getAccountID(gtsClient, false, c.accountNames[0])
	if err != nil {
		return fmt.Errorf("received an error while getting the account ID; %w", err)
	}

	if c.content == "" {
		return errors.New("the note content should not be empty")
	}

	if err := gtsClient.SetPrivateNote(accountID, c.content); err != nil {
		return fmt.Errorf("unable to add the private note to the account; %w", err)
	}

	fmt.Println("Successfully added the private note to the account.")

	return nil
}
