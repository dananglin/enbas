// SPDX-FileCopyrightText: 2024 Dan Anglin <d.n.i.anglin@gmail.com>
//
// SPDX-License-Identifier: GPL-3.0-or-later

package executor

import (
	"flag"
	"fmt"

	"codeflow.dananglin.me.uk/apollo/enbas/internal/client"
)

type DeleteExecutor struct {
	*flag.FlagSet

	topLevelFlags TopLevelFlags
	resourceType  string
	listID        string
}

func NewDeleteExecutor(tlf TopLevelFlags, name, summary string) *DeleteExecutor {
	deleteExe := DeleteExecutor{
		FlagSet:       flag.NewFlagSet(name, flag.ExitOnError),
		topLevelFlags: tlf,
	}

	deleteExe.StringVar(&deleteExe.resourceType, flagType, "", "specify the type of resource to delete")
	deleteExe.StringVar(&deleteExe.listID, flagListID, "", "specify the ID of the list to delete")

	deleteExe.Usage = commandUsageFunc(name, summary, deleteExe.FlagSet)

	return &deleteExe
}

func (d *DeleteExecutor) Execute() error {
	if d.resourceType == "" {
		return FlagNotSetError{flagText: flagType}
	}

	funcMap := map[string]func(*client.Client) error{
		resourceList: d.deleteList,
	}

	doFunc, ok := funcMap[d.resourceType]
	if !ok {
		return UnsupportedTypeError{resourceType: d.resourceType}
	}

	gtsClient, err := client.NewClientFromConfig(d.topLevelFlags.ConfigDir)
	if err != nil {
		return fmt.Errorf("unable to create the GoToSocial client: %w", err)
	}

	return doFunc(gtsClient)
}

func (d *DeleteExecutor) deleteList(gtsClient *client.Client) error {
	if d.listID == "" {
		return FlagNotSetError{flagText: flagListID}
	}

	if err := gtsClient.DeleteList(d.listID); err != nil {
		return fmt.Errorf("unable to delete the list: %w", err)
	}

	fmt.Println("The list was successfully deleted.")

	return nil
}
