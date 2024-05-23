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
	accountNames     AccountNames
}

func NewRemoveExecutor(tlf TopLevelFlags, name, summary string) *RemoveExecutor {
	emptyArr := make([]string, 0, 3)

	removeExe := RemoveExecutor{
		FlagSet:       flag.NewFlagSet(name, flag.ExitOnError),
		accountNames:  AccountNames(emptyArr),
		topLevelFlags: tlf,
	}

	removeExe.StringVar(&removeExe.resourceType, flagType, "", "specify the resource type to remove (e.g. account, note)")
	removeExe.StringVar(&removeExe.fromResourceType, flagFrom, "", "specify the resource type to remove from (e.g. list, account, etc)")
	removeExe.StringVar(&removeExe.listID, flagListID, "", "the ID of the list to remove from")
	removeExe.Var(&removeExe.accountNames, flagAccountName, "the name of the account to remove from the resource")

	removeExe.Usage = commandUsageFunc(name, summary, removeExe.FlagSet)

	return &removeExe
}

func (r *RemoveExecutor) Execute() error {
	if r.fromResourceType == "" {
		return FlagNotSetError{flagText: flagFrom}
	}

	funcMap := map[string]func(*client.Client) error{
		resourceList:    r.removeFromList,
		resourceAccount: r.removeFromAccount,
	}

	doFunc, ok := funcMap[r.fromResourceType]
	if !ok {
		return UnsupportedTypeError{resourceType: r.fromResourceType}
	}

	gtsClient, err := client.NewClientFromConfig(r.topLevelFlags.ConfigDir)
	if err != nil {
		return fmt.Errorf("unable to create the GoToSocial client; %w", err)
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

	for i := range r.accountNames {
		accountID, err := getTheirAccountID(gtsClient, r.accountNames[i])
		if err != nil {
			return fmt.Errorf("unable to get the account ID for %s, %w", r.accountNames[i], err)
		}

		accountIDs[i] = accountID
	}

	if err := gtsClient.RemoveAccountsFromList(r.listID, accountIDs); err != nil {
		return fmt.Errorf("unable to remove the accounts from the list; %w", err)
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
		return fmt.Errorf("unexpected number of accounts specified; want 1, got %d", len(r.accountNames))
	}

	accountID, err := getAccountID(gtsClient, false, r.accountNames[0], r.topLevelFlags.ConfigDir)
	if err != nil {
		return fmt.Errorf("received an error while getting the account ID; %w", err)
	}

	if err := gtsClient.SetPrivateNote(accountID, ""); err != nil {
		return fmt.Errorf("unable to remove the private note from the account; %w", err)
	}

	fmt.Println("Successfully removed the private note from the account.")

	return nil
}
