package main

import (
	"errors"
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

	command.StringVar(&command.resourceType, "type", "", "specify the type of resource to delete")
	command.StringVar(&command.listID, "list-id", "", "specify the ID of the list to delete")

	command.Usage = commandUsageFunc(name, summary, command.FlagSet)

	return &command
}

func (c *deleteCommand) Execute() error {
	if c.resourceType == "" {
		return errors.New("the type field is not set")
	}

	gtsClient, err := client.NewClientFromConfig()
	if err != nil {
		return fmt.Errorf("unable to create the GoToSocial client; %w", err)
	}

	funcMap := map[string]func(*client.Client) error{
		"lists": c.deleteList,
	}

	doFunc, ok := funcMap[c.resourceType]
	if !ok {
		return fmt.Errorf("unsupported resource type %q", c.resourceType)
	}

	return doFunc(gtsClient)
}

func (c *deleteCommand) deleteList(gtsClient *client.Client) error {
	if c.listID == "" {
		return errors.New("the list-id flag is not set")
	}

	if err := gtsClient.DeleteList(c.listID); err != nil {
		return fmt.Errorf("unable to delete the list; %w", err)
	}

	fmt.Println("The list was successfully deleted.")

	return nil
}
