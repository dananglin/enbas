// SPDX-FileCopyrightText: 2024 Dan Anglin <d.n.i.anglin@gmail.com>
//
// SPDX-License-Identifier: GPL-3.0-or-later

package executor

import (
	"flag"
	"fmt"

	"codeflow.dananglin.me.uk/apollo/enbas/internal/client"
)

type RemoveExecutor struct {
	*flag.FlagSet

	topLevelFlags    TopLevelFlags
	resourceType     string
	fromResourceType string
	listID           string
	statusID         string
	accountNames     AccountNames
}

func NewRemoveExecutor(tlf TopLevelFlags, name, summary string) *RemoveExecutor {
	emptyArr := make([]string, 0, 3)

	removeExe := RemoveExecutor{
		FlagSet:       flag.NewFlagSet(name, flag.ExitOnError),
		accountNames:  AccountNames(emptyArr),
		topLevelFlags: tlf,
	}

	removeExe.StringVar(&removeExe.resourceType, flagType, "", "Specify the resource type to remove (e.g. account, note)")
	removeExe.StringVar(&removeExe.fromResourceType, flagFrom, "", "Specify the resource type to remove from (e.g. list, account, etc)")
	removeExe.StringVar(&removeExe.listID, flagListID, "", "The ID of the list to remove from")
	removeExe.StringVar(&removeExe.statusID, flagStatusID, "", "The ID of the status")
	removeExe.Var(&removeExe.accountNames, flagAccountName, "The name of the account to remove from the resource")

	removeExe.Usage = commandUsageFunc(name, summary, removeExe.FlagSet)

	return &removeExe
}

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

	gtsClient, err := client.NewClientFromConfig(r.topLevelFlags.ConfigDir)
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

	if len(r.accountNames) == 0 {
		return NoAccountSpecifiedError{}
	}

	accountIDs := make([]string, len(r.accountNames))

	for ind := range r.accountNames {
		accountID, err := getTheirAccountID(gtsClient, r.accountNames[ind])
		if err != nil {
			return fmt.Errorf("unable to get the account ID for %s: %w", r.accountNames[ind], err)
		}

		accountIDs[ind] = accountID
	}

	if err := gtsClient.RemoveAccountsFromList(r.listID, accountIDs); err != nil {
		return fmt.Errorf("unable to remove the accounts from the list: %w", err)
	}

	fmt.Println("Successfully removed the account(s) from the list.")

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
	if len(r.accountNames) != 1 {
		return fmt.Errorf("unexpected number of accounts specified: want 1, got %d", len(r.accountNames))
	}

	accountID, err := getAccountID(gtsClient, false, r.accountNames[0], r.topLevelFlags.ConfigDir)
	if err != nil {
		return fmt.Errorf("received an error while getting the account ID: %w", err)
	}

	if err := gtsClient.SetPrivateNote(accountID, ""); err != nil {
		return fmt.Errorf("unable to remove the private note from the account: %w", err)
	}

	fmt.Println("Successfully removed the private note from the account.")

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

	fmt.Println("Successfully removed the status from your bookmarks.")

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

	fmt.Printf("Successfully removed the %s from the status.\n", r.resourceType)

	return nil
}

func (r *RemoveExecutor) removeBoostFromStatus(gtsClient *client.Client) error {
	if err := gtsClient.UnreblogStatus(r.statusID); err != nil {
		return fmt.Errorf("unable to remove the boost from the status: %w", err)
	}

	fmt.Println("Successfully removed the boost from the status.")

	return nil
}
