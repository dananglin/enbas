package executor

import (
	"fmt"

	"codeflow.dananglin.me.uk/apollo/enbas/internal/client"
	"codeflow.dananglin.me.uk/apollo/enbas/internal/model"
)

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

	gtsClient, err := client.NewClientFromFile(e.config.CredentialsFile)
	if err != nil {
		return fmt.Errorf("unable to create the GoToSocial client: %w", err)
	}

	return doFunc(gtsClient)
}

func (e *EditExecutor) editList(gtsClient *client.Client) error {
	if e.listID == "" {
		return FlagNotSetError{flagText: flagListID}
	}

	list, err := gtsClient.GetList(e.listID)
	if err != nil {
		return fmt.Errorf("unable to get the list: %w", err)
	}

	if e.listTitle != "" {
		list.Title = e.listTitle
	}

	if e.listRepliesPolicy != "" {
		parsedListRepliesPolicy, err := model.ParseListRepliesPolicy(e.listRepliesPolicy)
		if err != nil {
			return err
		}

		list.RepliesPolicy = parsedListRepliesPolicy
	}

	updatedList, err := gtsClient.UpdateList(list)
	if err != nil {
		return fmt.Errorf("unable to update the list: %w", err)
	}

	e.printer.PrintSuccess("Successfully updated the list.")
	e.printer.PrintList(updatedList)

	return nil
}
