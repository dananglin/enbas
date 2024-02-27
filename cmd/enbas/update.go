package main

import (
	"errors"
	"flag"
	"fmt"

	"codeflow.dananglin.me.uk/apollo/enbas/internal/client"
	"codeflow.dananglin.me.uk/apollo/enbas/internal/model"
)

type updateCommand struct {
	*flag.FlagSet

	resourceType      string
	listID            string
	listTitle         string
	listRepliesPolicy string
}

func newUpdateCommand(name, summary string) *updateCommand {
	command := updateCommand{
		FlagSet: flag.NewFlagSet(name, flag.ExitOnError),
	}

	command.StringVar(&command.resourceType, "type", "", "specify the type of resource to update")
	command.StringVar(&command.listID, "list-id", "", "specify the ID of the list to update")
	command.StringVar(&command.listTitle, "list-title", "", "specify the title of the list")
	command.StringVar(&command.listRepliesPolicy, "list-replies-policy", "", "specify the policy of the replies for this list (valid values are followed, list and none)")

	command.Usage = commandUsageFunc(name, summary, command.FlagSet)

	return &command
}

func (c *updateCommand) Execute() error {
	if c.resourceType == "" {
		return errors.New("the type field is not set")
	}

	gtsClient, err := client.NewClientFromConfig()
	if err != nil {
		return fmt.Errorf("unable to create the GoToSocial client; %w", err)
	}

	funcMap := map[string]func(*client.Client) error{
		"lists": c.updateList,
	}

	doFunc, ok := funcMap[c.resourceType]
	if !ok {
		return fmt.Errorf("unsupported resource type %q", c.resourceType)
	}

	return doFunc(gtsClient)
}

func (c *updateCommand) updateList(gtsClient *client.Client) error {
	if c.listID == "" {
		return errors.New("the list-id flag is not set")
	}

	list, err := gtsClient.GetList(c.listID)
	if err != nil {
		return fmt.Errorf("unable to get the list; %w", err)
	}

	if c.listTitle != "" {
		list.Title = c.listTitle
	}

	if c.listRepliesPolicy != "" {
		repliesPolicy, err := model.ParseListRepliesPolicy(c.listRepliesPolicy)
		if err != nil {
			return fmt.Errorf("unable to parse the list replies policy; %w", err)
		}

		list.RepliesPolicy = repliesPolicy
	}

	updatedList, err := gtsClient.UpdateList(list)
	if err != nil {
		return fmt.Errorf("unable to update the list; %w", err)
	}

	fmt.Println("Successfully updated the list.")
	fmt.Printf("\n%s\n", updatedList)

	return nil
}
