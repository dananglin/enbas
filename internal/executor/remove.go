package executor

import (
	"fmt"

	"codeflow.dananglin.me.uk/apollo/enbas/internal/client"
)

func (r *RemoveExecutor) Execute() error {
	if r.fromResourceType == "" {
		return FlagNotSetError{flagText: flagFrom}
	}

	funcMap := map[string]func(*client.Client) error{
		resourceList:      r.removeFromList,
		resourceAccount:   r.removeFromAccount,
		resourceBookmarks: r.removeFromBookmarks,
		resourceStatus:    r.removeFromStatus,
	}

	doFunc, ok := funcMap[r.fromResourceType]
	if !ok {
		return UnsupportedTypeError{resourceType: r.fromResourceType}
	}

	gtsClient, err := client.NewClientFromFile(r.config.CredentialsFile)
	if err != nil {
		return fmt.Errorf("unable to create the GoToSocial client: %w", err)
	}

	return doFunc(gtsClient)
}

func (r *RemoveExecutor) removeFromList(gtsClient *client.Client) error {
	funcMap := map[string]func(*client.Client) error{
		resourceAccount: r.removeAccountsFromList,
	}

	doFunc, ok := funcMap[r.resourceType]
	if !ok {
		return UnsupportedRemoveOperationError{
			ResourceType:           r.resourceType,
			RemoveFromResourceType: r.fromResourceType,
		}
	}

	return doFunc(gtsClient)
}

func (r *RemoveExecutor) removeAccountsFromList(gtsClient *client.Client) error {
	if r.listID == "" {
		return FlagNotSetError{flagText: flagListID}
	}

	if r.accountNames.Empty() {
		return NoAccountSpecifiedError{}
	}

	accounts, err := getOtherAccounts(gtsClient, r.accountNames)
	if err != nil {
		return fmt.Errorf("unable to get the accounts: %w", err)
	}

	accountIDs := make([]string, len(accounts))

	for ind := range accounts {
		accountIDs[ind] = accounts[ind].ID
	}

	if err := gtsClient.RemoveAccountsFromList(r.listID, accountIDs); err != nil {
		return fmt.Errorf("unable to remove the accounts from the list: %w", err)
	}

	r.printer.PrintSuccess("Successfully removed the account(s) from the list.")

	return nil
}

func (r *RemoveExecutor) removeFromAccount(gtsClient *client.Client) error {
	funcMap := map[string]func(*client.Client) error{
		resourceNote: r.removeNoteFromAccount,
	}

	doFunc, ok := funcMap[r.resourceType]
	if !ok {
		return UnsupportedRemoveOperationError{
			ResourceType:           r.resourceType,
			RemoveFromResourceType: r.fromResourceType,
		}
	}

	return doFunc(gtsClient)
}

func (r *RemoveExecutor) removeNoteFromAccount(gtsClient *client.Client) error {
	accountID, err := getAccountID(gtsClient, false, r.accountNames, r.config.CredentialsFile)
	if err != nil {
		return fmt.Errorf("received an error while getting the account ID: %w", err)
	}

	if err := gtsClient.SetPrivateNote(accountID, ""); err != nil {
		return fmt.Errorf("unable to remove the private note from the account: %w", err)
	}

	r.printer.PrintSuccess("Successfully removed the private note from the account.")

	return nil
}

func (r *RemoveExecutor) removeFromBookmarks(gtsClient *client.Client) error {
	funcMap := map[string]func(*client.Client) error{
		resourceStatus: r.removeStatusFromBookmarks,
	}

	doFunc, ok := funcMap[r.resourceType]
	if !ok {
		return UnsupportedRemoveOperationError{
			ResourceType:           r.resourceType,
			RemoveFromResourceType: r.fromResourceType,
		}
	}

	return doFunc(gtsClient)
}

func (r *RemoveExecutor) removeStatusFromBookmarks(gtsClient *client.Client) error {
	if r.statusID == "" {
		return FlagNotSetError{flagText: flagStatusID}
	}

	if err := gtsClient.RemoveStatusFromBookmarks(r.statusID); err != nil {
		return fmt.Errorf("unable to remove the status from your bookmarks: %w", err)
	}

	r.printer.PrintSuccess("Successfully removed the status from your bookmarks.")

	return nil
}

func (r *RemoveExecutor) removeFromStatus(gtsClient *client.Client) error {
	if r.statusID == "" {
		return FlagNotSetError{flagText: flagStatusID}
	}

	funcMap := map[string]func(*client.Client) error{
		resourceStar:  r.removeStarFromStatus,
		resourceLike:  r.removeStarFromStatus,
		resourceBoost: r.removeBoostFromStatus,
	}

	doFunc, ok := funcMap[r.resourceType]
	if !ok {
		return UnsupportedRemoveOperationError{
			ResourceType:           r.resourceType,
			RemoveFromResourceType: r.fromResourceType,
		}
	}

	return doFunc(gtsClient)
}

func (r *RemoveExecutor) removeStarFromStatus(gtsClient *client.Client) error {
	if err := gtsClient.UnlikeStatus(r.statusID); err != nil {
		return fmt.Errorf("unable to remove the %s from the status: %w", r.resourceType, err)
	}

	r.printer.PrintSuccess("Successfully removed the " + r.resourceType + " from the status.")

	return nil
}

func (r *RemoveExecutor) removeBoostFromStatus(gtsClient *client.Client) error {
	if err := gtsClient.UnreblogStatus(r.statusID); err != nil {
		return fmt.Errorf("unable to remove the boost from the status: %w", err)
	}

	r.printer.PrintSuccess("Successfully removed the boost from the status.")

	return nil
}
