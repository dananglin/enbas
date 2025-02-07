package executor

import (
	"fmt"
	"net/rpc"

	"codeflow.dananglin.me.uk/apollo/enbas/internal/gtsclient"
	"codeflow.dananglin.me.uk/apollo/enbas/internal/printer"
	"codeflow.dananglin.me.uk/apollo/enbas/internal/server"
)

func (r *RemoveExecutor) Execute() error {
	if r.fromResourceType == "" {
		return FlagNotSetError{flagText: flagFrom}
	}

	funcMap := map[string]func(*rpc.Client) error{
		resourceList:      r.removeFromList,
		resourceAccount:   r.removeFromAccount,
		resourceBookmarks: r.removeFromBookmarks,
		resourceStatus:    r.removeFromStatus,
	}

	doFunc, ok := funcMap[r.fromResourceType]
	if !ok {
		return UnsupportedTypeError{resourceType: r.fromResourceType}
	}

	client, err := server.Connect(r.config.Server, r.configDir)
	if err != nil {
		return fmt.Errorf("error creating the client for the daemon process: %w", err)
	}
	defer client.Close()

	return doFunc(client)
}

func (r *RemoveExecutor) removeFromList(client *rpc.Client) error {
	funcMap := map[string]func(*rpc.Client) error{
		resourceAccount: r.removeAccountsFromList,
	}

	doFunc, ok := funcMap[r.resourceType]
	if !ok {
		return UnsupportedRemoveOperationError{
			resourceType:           r.resourceType,
			removeFromResourceType: r.fromResourceType,
		}
	}

	return doFunc(client)
}

func (r *RemoveExecutor) removeAccountsFromList(client *rpc.Client) error {
	if r.listID == "" {
		return MissingIDError{
			resource: resourceList,
			action:   "remove from",
		}
	}

	if r.accountNames.Empty() {
		return NoAccountSpecifiedError{}
	}

	accounts, err := getMultipleAccounts(client, r.accountNames)
	if err != nil {
		return fmt.Errorf("unable to get the accounts: %w", err)
	}

	accountIDs := make([]string, len(accounts))

	for ind := range accounts {
		accountIDs[ind] = accounts[ind].ID
	}

	if err := client.Call(
		"GTSClient.RemoveAccountsFromList",
		gtsclient.RemoveAccountsFromListArgs{
			ListID:     r.listID,
			AccountIDs: accountIDs,
		},
		nil,
	); err != nil {
		return fmt.Errorf("error removing the accounts from the list: %w", err)
	}

	printer.PrintSuccess(r.printSettings, "Successfully removed the account(s) from the list.")

	return nil
}

func (r *RemoveExecutor) removeFromAccount(client *rpc.Client) error {
	funcMap := map[string]func(*rpc.Client) error{
		resourceNote: r.removeNoteFromAccount,
	}

	doFunc, ok := funcMap[r.resourceType]
	if !ok {
		return UnsupportedRemoveOperationError{
			resourceType:           r.resourceType,
			removeFromResourceType: r.fromResourceType,
		}
	}

	return doFunc(client)
}

func (r *RemoveExecutor) removeNoteFromAccount(client *rpc.Client) error {
	accountID, err := getAccountID(client, r.accountNames)
	if err != nil {
		return fmt.Errorf("received an error while getting the account ID: %w", err)
	}

	if err := client.Call(
		"GTSClient.SetPrivateNote",
		gtsclient.SetPrivateNoteArgs{
			AccountID: accountID,
			Note:      "",
		},
		nil,
	); err != nil {
		return fmt.Errorf("unable to remove the private note from the account: %w", err)
	}

	printer.PrintSuccess(r.printSettings, "Successfully removed the private note from the account.")

	return nil
}

func (r *RemoveExecutor) removeFromBookmarks(client *rpc.Client) error {
	funcMap := map[string]func(*rpc.Client) error{
		resourceStatus: r.removeStatusFromBookmarks,
	}

	doFunc, ok := funcMap[r.resourceType]
	if !ok {
		return UnsupportedRemoveOperationError{
			resourceType:           r.resourceType,
			removeFromResourceType: r.fromResourceType,
		}
	}

	return doFunc(client)
}

func (r *RemoveExecutor) removeStatusFromBookmarks(client *rpc.Client) error {
	if r.statusID == "" {
		return MissingIDError{
			resource: resourceStatus,
			action:   "remove",
		}
	}

	if err := client.Call("GTSClient.RemoveStatusFromBookmarks", r.statusID, nil); err != nil {
		return fmt.Errorf("error removing the status from your bookmarks: %w", err)
	}

	printer.PrintSuccess(r.printSettings, "Successfully removed the status from your bookmarks.")

	return nil
}

func (r *RemoveExecutor) removeFromStatus(client *rpc.Client) error {
	if r.statusID == "" {
		return MissingIDError{
			resource: resourceStatus,
			action:   "remove from",
		}
	}

	funcMap := map[string]func(*rpc.Client) error{
		resourceStar:  r.removeStarFromStatus,
		resourceLike:  r.removeStarFromStatus,
		resourceBoost: r.removeBoostFromStatus,
	}

	doFunc, ok := funcMap[r.resourceType]
	if !ok {
		return UnsupportedRemoveOperationError{
			resourceType:           r.resourceType,
			removeFromResourceType: r.fromResourceType,
		}
	}

	return doFunc(client)
}

func (r *RemoveExecutor) removeStarFromStatus(client *rpc.Client) error {
	if err := client.Call("GTSClient.UnlikeStatus", r.statusID, nil); err != nil {
		return fmt.Errorf("error removing the %s from the status: %w", r.resourceType, err)
	}

	printer.PrintSuccess(r.printSettings, "Successfully removed the "+r.resourceType+" from the status.")

	return nil
}

func (r *RemoveExecutor) removeBoostFromStatus(client *rpc.Client) error {
	if err := client.Call("GTSClient.UnreblogStatus", r.statusID, nil); err != nil {
		return fmt.Errorf("unable to remove the boost from the status: %w", err)
	}

	printer.PrintSuccess(r.printSettings, "Successfully removed the boost from the status.")

	return nil
}
