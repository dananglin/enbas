package executor

import (
	"flag"
	"fmt"

	"codeflow.dananglin.me.uk/apollo/enbas/internal/client"
	"codeflow.dananglin.me.uk/apollo/enbas/internal/model"
	"codeflow.dananglin.me.uk/apollo/enbas/internal/utilities"
)

type EditExecutor struct {
	*flag.FlagSet

	topLevelFlags     TopLevelFlags
	resourceType      string
	listID            string
	listTitle         string
	listRepliesPolicy string
}

func NewEditExecutor(tlf TopLevelFlags, name, summary string) *EditExecutor {
	editExe := EditExecutor{
		FlagSet:       flag.NewFlagSet(name, flag.ExitOnError),
		topLevelFlags: tlf,
	}

	editExe.StringVar(&editExe.resourceType, flagType, "", "specify the type of resource to update")
	editExe.StringVar(&editExe.listID, flagListID, "", "specify the ID of the list to update")
	editExe.StringVar(&editExe.listTitle, flagListTitle, "", "specify the title of the list")
	editExe.StringVar(&editExe.listRepliesPolicy, flagListRepliesPolicy, "", "specify the policy of the replies for this list (valid values are followed, list and none)")

	editExe.Usage = commandUsageFunc(name, summary, editExe.FlagSet)

	return &editExe
}

func (e *EditExecutor) Execute() error {
	if e.resourceType == "" {
		return FlagNotSetError{flagText: flagType}
	}

	funcMap := map[string]func(*client.Client) error{
		resourceList: e.editList,
	}

	doFunc, ok := funcMap[e.resourceType]
	if !ok {
		return UnsupportedTypeError{resourceType: e.resourceType}
	}

	gtsClient, err := client.NewClientFromConfig(e.topLevelFlags.ConfigDir)
	if err != nil {
		return fmt.Errorf("unable to create the GoToSocial client; %w", err)
	}

	return doFunc(gtsClient)
}

func (e *EditExecutor) editList(gtsClient *client.Client) error {
	if e.listID == "" {
		return FlagNotSetError{flagText: flagListID}
	}

	list, err := gtsClient.GetList(e.listID)
	if err != nil {
		return fmt.Errorf("unable to get the list; %w", err)
	}

	if e.listTitle != "" {
		list.Title = e.listTitle
	}

	if e.listRepliesPolicy != "" {
		parsedListRepliesPolicy := model.ParseListRepliesPolicy(e.listRepliesPolicy)
		if parsedListRepliesPolicy == model.ListRepliesPolicyUnknown {
			return InvalidListRepliesPolicyError{Policy: e.listRepliesPolicy}
		}

		list.RepliesPolicy = parsedListRepliesPolicy
	}

	updatedList, err := gtsClient.UpdateList(list)
	if err != nil {
		return fmt.Errorf("unable to update the list; %w", err)
	}

	fmt.Println("Successfully updated the list.")
	utilities.Display(updatedList, *e.topLevelFlags.NoColor)

	return nil
}
