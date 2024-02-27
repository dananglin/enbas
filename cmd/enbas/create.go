package main

import (
	"errors"
	"flag"
	"fmt"

	"codeflow.dananglin.me.uk/apollo/enbas/internal/client"
	"codeflow.dananglin.me.uk/apollo/enbas/internal/model"
)

type createCommand struct {
	*flag.FlagSet

	resourceType      string
	listTitle         string
	listRepliesPolicy string
}

func newCreateCommand(name, summary string) *createCommand {
	command := createCommand{
		FlagSet: flag.NewFlagSet(name, flag.ExitOnError),
	}

	command.StringVar(&command.resourceType, "type", "", "specify the type of resource to create")
	command.StringVar(&command.listTitle, "list-title", "", "specify the title of the list")
	command.StringVar(&command.listRepliesPolicy, "list-replies-policy", "list", "specify the policy of the replies for this list (valid values are followed, list and none)")

	command.Usage = commandUsageFunc(name, summary, command.FlagSet)

	return &command
}

func (c *createCommand) Execute() error {
	if c.resourceType == "" {
		return errors.New("the type field is not set")
	}

	gtsClient, err := client.NewClientFromConfig()
	if err != nil {
		return fmt.Errorf("unable to create the GoToSocial client; %w", err)
	}

	funcMap := map[string]func(*client.Client) error{
		"lists": c.createLists,
	}

	doFunc, ok := funcMap[c.resourceType]
	if !ok {
		return fmt.Errorf("unsupported type %q", c.resourceType)
	}

	return doFunc(gtsClient)
}

func (c *createCommand) createLists(gtsClient *client.Client) error {
	if c.listTitle == "" {
		return errors.New("the list-title flag is not set")
	}

	repliesPolicy, err := model.ParseListRepliesPolicy(c.listRepliesPolicy)
	if err != nil {
		return fmt.Errorf("unable to parse the list replies policy; %w", err)
	}

	list, err := gtsClient.CreateList(c.listTitle, repliesPolicy)
	if err != nil {
		return fmt.Errorf("unable to create the list; %w", err)
	}

	fmt.Println("Successfully created the following list:")
	fmt.Printf("\n%s\n", list)

	return nil
}
