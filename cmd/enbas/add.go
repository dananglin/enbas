package main

import (
	"flag"
	"fmt"

	"codeflow.dananglin.me.uk/apollo/enbas/internal/client"
)

type addCommand struct {
	*flag.FlagSet

	toResourceType string
	listID         string
	accountIDs     accountIDs
}

func newAddCommand(name, summary string) *addCommand {
	emptyArr := make([]string, 0, 3)

	command := addCommand{
		FlagSet:    flag.NewFlagSet(name, flag.ExitOnError),
		accountIDs: accountIDs(emptyArr),
	}

	command.StringVar(&command.toResourceType, addToFlag, "", "specify the type of resource to add to")
	command.StringVar(&command.listID, listIDFlag, "", "the ID of the list to add to")
	command.Var(&command.accountIDs, accountIDFlag, "the ID of the account to add to the list")

	command.Usage = commandUsageFunc(name, summary, command.FlagSet)

	return &command
}

func (c *addCommand) Execute() error {
	if c.toResourceType == "" {
		return flagNotSetError{flagText: "add-to"}
	}

	funcMap := map[string]func(*client.Client) error{
		listResource: c.addAccountsToList,
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

func (c *addCommand) addAccountsToList(gtsClient *client.Client) error {
	if c.listID == "" {
		return flagNotSetError{flagText: listIDFlag}
	}

	if len(c.accountIDs) == 0 {
		return noAccountIDsSpecifiedError{}
	}

	if err := gtsClient.AddAccountsToList(c.listID, []string(c.accountIDs)); err != nil {
		return fmt.Errorf("unable to add the accounts to the list; %w", err)
	}

	fmt.Println("Successfully added the account(s) to the list.")

	return nil
}
