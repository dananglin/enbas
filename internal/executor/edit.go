// SPDX-FileCopyrightText: 2024 Dan Anglin <d.n.i.anglin@gmail.com>
//
// SPDX-License-Identifier: GPL-3.0-or-later

package executor

import (
	"flag"
	"fmt"

	"codeflow.dananglin.me.uk/apollo/enbas/internal/client"
	"codeflow.dananglin.me.uk/apollo/enbas/internal/config"
	"codeflow.dananglin.me.uk/apollo/enbas/internal/model"
	"codeflow.dananglin.me.uk/apollo/enbas/internal/printer"
)

type EditExecutor struct {
	*flag.FlagSet

	printer           *printer.Printer
	config            *config.Config
	resourceType      string
	listID            string
	listTitle         string
	listRepliesPolicy string
}

func NewEditExecutor(printer *printer.Printer, config *config.Config, name, summary string) *EditExecutor {
	editExe := EditExecutor{
		FlagSet: flag.NewFlagSet(name, flag.ExitOnError),

		printer: printer,
		config:  config,
	}

	editExe.StringVar(&editExe.resourceType, flagType, "", "Specify the type of resource to update")
	editExe.StringVar(&editExe.listID, flagListID, "", "Specify the ID of the list to update")
	editExe.StringVar(&editExe.listTitle, flagListTitle, "", "Specify the title of the list")
	editExe.StringVar(&editExe.listRepliesPolicy, flagListRepliesPolicy, "", "Specify the policy of the replies for this list (valid values are followed, list and none)")

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
