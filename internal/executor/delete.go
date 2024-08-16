package executor

import (
	"errors"
	"fmt"
	"path/filepath"

	"codeflow.dananglin.me.uk/apollo/enbas/internal/client"
	"codeflow.dananglin.me.uk/apollo/enbas/internal/utilities"
)

func (d *DeleteExecutor) Execute() error {
	if d.resourceType == "" {
		return FlagNotSetError{flagText: flagType}
	}

	funcMap := map[string]func(*client.Client) error{
		resourceList:   d.deleteList,
		resourceStatus: d.deleteStatus,
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

func (d *DeleteExecutor) deleteStatus(gtsClient *client.Client) error {
	if d.statusID == "" {
		return FlagNotSetError{flagText: flagStatusID}
	}

	status, err := gtsClient.GetStatus(d.statusID)
	if err != nil {
		return fmt.Errorf("unable to get the status: %w", err)
	}

	myAccountID, err := getAccountID(gtsClient, true, nil)
	if err != nil {
		return fmt.Errorf("unable to get your account ID: %w", err)
	}

	if status.Account.ID != myAccountID {
		return errors.New("unable to delete the status because the status does not belong to you")
	}

	text, err := gtsClient.DeleteStatus(d.statusID)
	if err != nil {
		return fmt.Errorf("unable to delete the status: %w", err)
	}

	d.printer.PrintSuccess("The status was successfully deleted.")

	if d.saveText {
		cacheDir := filepath.Join(
			utilities.CalculateCacheDir(
				d.config.CacheDirectory,
				utilities.GetFQDN(gtsClient.Authentication.Instance),
			),
			"statuses",
		)

		if err := utilities.EnsureDirectory(cacheDir); err != nil {
			return fmt.Errorf("unable to ensure the existence of the directory %q: %w", cacheDir, err)
		}

		path := filepath.Join(cacheDir, fmt.Sprintf("deleted-status-%s.txt", d.statusID))

		if err := utilities.SaveTextToFile(path, text); err != nil {
			return fmt.Errorf("unable to save the text to %q: %w", path, err)
		}

		d.printer.PrintSuccess("The text was successfully saved to '" + path + "'.")
	}

	return nil
}
