package executor

import (
	"flag"
	"fmt"

	"codeflow.dananglin.me.uk/apollo/enbas/internal/client"
	"codeflow.dananglin.me.uk/apollo/enbas/internal/model"
)

type CreateExecutor struct {
	*flag.FlagSet

	topLevelFlags     TopLevelFlags
	resourceType      string
	listTitle         string
	listRepliesPolicy string
}

func NewCreateExecutor(tlf TopLevelFlags, name, summary string) *CreateExecutor {
	createExe := CreateExecutor{
		FlagSet: flag.NewFlagSet(name, flag.ExitOnError),

		topLevelFlags: tlf,
	}

	createExe.StringVar(&createExe.resourceType, flagType, "", "specify the type of resource to create")
	createExe.StringVar(&createExe.listTitle, flagListTitle, "", "specify the title of the list")
	createExe.StringVar(&createExe.listRepliesPolicy, flagListRepliesPolicy, "list", "specify the policy of the replies for this list (valid values are followed, list and none)")

	createExe.Usage = commandUsageFunc(name, summary, createExe.FlagSet)

	return &createExe
}

func (c *CreateExecutor) Execute() error {
	if c.resourceType == "" {
		return FlagNotSetError{flagText: flagType}
	}

	gtsClient, err := client.NewClientFromConfig(c.topLevelFlags.ConfigDir)
	if err != nil {
		return fmt.Errorf("unable to create the GoToSocial client; %w", err)
	}

	funcMap := map[string]func(*client.Client) error{
		resourceList: c.createList,
	}

	doFunc, ok := funcMap[c.resourceType]
	if !ok {
		return UnsupportedTypeError{resourceType: c.resourceType}
	}

	return doFunc(gtsClient)
}

func (c *CreateExecutor) createList(gtsClient *client.Client) error {
	if c.listTitle == "" {
		return FlagNotSetError{flagText: flagListTitle}
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
