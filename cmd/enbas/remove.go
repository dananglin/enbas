package main

import (
	"flag"
	"fmt"

	"codeflow.dananglin.me.uk/apollo/enbas/internal/client"
)

type removeCommand struct {
	*flag.FlagSet

	fromResourceType string
	listID           string
	accountIDs       accountIDs
}

func newRemoveCommand(name, summary string) *removeCommand {
	emptyArr := make([]string, 0, 3)

	command := removeCommand{
		FlagSet:    flag.NewFlagSet(name, flag.ExitOnError),
		accountIDs: accountIDs(emptyArr),
	}

	command.StringVar(&command.fromResourceType, removeFromFlag, "", "specify the type of resource to remove from")
	command.StringVar(&command.listID, listIDFlag, "", "the ID of the list to remove from")
	command.Var(&command.accountIDs, accountIDFlag, "the ID of the account to remove from the list")

	command.Usage = commandUsageFunc(name, summary, command.FlagSet)

	return &command
}

func (c *removeCommand) Execute() error {
	if c.fromResourceType == "" {
		return flagNotSetError{flagText: "remove-from"}
	}

	funcMap := map[string]func(*client.Client) error{
		listResource: c.removeAccountsFromList,
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

func (c *removeCommand) removeAccountsFromList(gtsClient *client.Client) error {
	if c.listID == "" {
		return flagNotSetError{flagText: listIDFlag}
	}

	if len(c.accountIDs) == 0 {
		return noAccountIDsSpecifiedError{}
	}

	if err := gtsClient.RemoveAccountsFromList(c.listID, []string(c.accountIDs)); err != nil {
		return fmt.Errorf("unable to remove the accounts from the list; %w", err)
	}

	fmt.Println("Successfully removed the account(s) from the list.")

	return nil
}
