package executor

import (
	"fmt"
	"net/rpc"
	"path/filepath"

	"codeflow.dananglin.me.uk/apollo/enbas/internal/gtsclient"
	"codeflow.dananglin.me.uk/apollo/enbas/internal/model"
	"codeflow.dananglin.me.uk/apollo/enbas/internal/printer"
	"codeflow.dananglin.me.uk/apollo/enbas/internal/server"
	"codeflow.dananglin.me.uk/apollo/enbas/internal/utilities"
)

func (d *DeleteExecutor) Execute() error {
	if d.resourceType == "" {
		return FlagNotSetError{flagText: flagType}
	}

	funcMap := map[string]func(*rpc.Client) error{
		resourceList:   d.deleteList,
		resourceStatus: d.deleteStatus,
	}

	doFunc, ok := funcMap[d.resourceType]
	if !ok {
		return UnsupportedTypeError{resourceType: d.resourceType}
	}

	client, err := server.Connect(d.config.Server, d.configDir)
	if err != nil {
		return fmt.Errorf("error creating the client for the daemon process: %w", err)
	}
	defer client.Close()

	return doFunc(client)
}

func (d *DeleteExecutor) deleteList(client *rpc.Client) error {
	if d.listID == "" {
		return MissingIDError{
			resource: resourceList,
			action:   "delete",
		}
	}

	if err := client.Call("GTSClient.DeleteList", d.listID, nil); err != nil {
		return fmt.Errorf("unable to delete the list: %w", err)
	}

	printer.PrintSuccess(d.printSettings, "The list was successfully deleted.")

	return nil
}

func (d *DeleteExecutor) deleteStatus(client *rpc.Client) error {
	if d.statusID == "" {
		return MissingIDError{
			resource: resourceStatus,
			action:   "delete",
		}
	}

	var status model.Status
	if err := client.Call("GTSClient.GetStatus", d.statusID, &status); err != nil {
		return fmt.Errorf("unable to retrieve the status: %w", err)
	}

	myAccountID, err := getAccountID(client, true, nil)
	if err != nil {
		return fmt.Errorf("unable to get your account ID: %w", err)
	}

	if status.Account.ID != myAccountID {
		return Error{"unable to delete the status because the status does not belong to you"}
	}

	var text string
	if err := client.Call("GTSClient.DeleteStatus", d.statusID, &text); err != nil {
		return fmt.Errorf("error deleting the status: %w", err)
	}

	printer.PrintSuccess(d.printSettings, "The status was successfully deleted.")

	if d.saveText {
		var instanceURL string
		if err := client.Call("GTSClient.GetInstanceURL", gtsclient.NoRPCArgs{}, &instanceURL); err != nil {
			return fmt.Errorf("unable to get the instance URL: %w", err)
		}

		cacheDir, err := utilities.CalculateStatusesCacheDir(d.config.CacheDirectory, instanceURL)
		if err != nil {
			return fmt.Errorf("unable to get the cache directory for the status: %w", err)
		}

		if err := utilities.EnsureDirectory(cacheDir); err != nil {
			return fmt.Errorf("unable to ensure the existence of the directory %q: %w", cacheDir, err)
		}

		path := filepath.Join(cacheDir, fmt.Sprintf("deleted-status-%s.txt", d.statusID))

		if err := utilities.SaveTextToFile(path, text); err != nil {
			return fmt.Errorf("unable to save the text to %q: %w", path, err)
		}

		printer.PrintSuccess(d.printSettings, "The text was successfully saved to '"+path+"'.")
	}

	return nil
}
