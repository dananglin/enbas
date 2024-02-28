package main

import (
	"flag"
	"fmt"

	"codeflow.dananglin.me.uk/apollo/enbas/internal/client"
)

type deleteCommand struct {
	*flag.FlagSet

	resourceType string
	listID       string
}

func newDeleteCommand(name, summary string) *deleteCommand {
	command := deleteCommand{
		FlagSet: flag.NewFlagSet(name, flag.ExitOnError),
	}

	command.StringVar(&command.resourceType, resourceTypeFlag, "", "specify the type of resource to delete")
	command.StringVar(&command.listID, listIDFlag, "", "specify the ID of the list to delete")

	command.Usage = commandUsageFunc(name, summary, command.FlagSet)

	return &command
}

func (c *deleteCommand) Execute() error {
	if c.resourceType == "" {
		return flagNotSetError{flagText: resourceTypeFlag}
	}

	funcMap := map[string]func(*client.Client) error{
		listResource: c.deleteList,
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

func (c *deleteCommand) deleteList(gtsClient *client.Client) error {
	if c.listID == "" {
		return flagNotSetError{flagText: listIDFlag}
	}

	if err := gtsClient.DeleteList(c.listID); err != nil {
		return fmt.Errorf("unable to delete the list; %w", err)
	}

	fmt.Println("The list was successfully deleted.")

	return nil
}
