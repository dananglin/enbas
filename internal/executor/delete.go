package executor

import (
	"fmt"

	"codeflow.dananglin.me.uk/apollo/enbas/internal/client"
)

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

	gtsClient, err := client.NewClientFromFile(d.config.CredentialsFile)
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

	d.printer.PrintSuccess("The list was successfully deleted.")

	return nil
}
