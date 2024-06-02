// SPDX-FileCopyrightText: 2024 Dan Anglin <d.n.i.anglin@gmail.com>
//
// SPDX-License-Identifier: GPL-3.0-or-later

package executor

import (
	"flag"
	"fmt"

	"codeflow.dananglin.me.uk/apollo/enbas/internal/client"
)

type AddExecutor struct {
	*flag.FlagSet

	topLevelFlags  TopLevelFlags
	resourceType   string
	toResourceType string
	listID         string
	accountNames   AccountNames
	content        string
}

func NewAddExecutor(tlf TopLevelFlags, name, summary string) *AddExecutor {
	emptyArr := make([]string, 0, 3)

	addExe := AddExecutor{
		FlagSet:       flag.NewFlagSet(name, flag.ExitOnError),
		accountNames:  AccountNames(emptyArr),
		topLevelFlags: tlf,
	}

	addExe.StringVar(&addExe.resourceType, flagType, "", "specify the resource type to add (e.g. account, note)")
	addExe.StringVar(&addExe.toResourceType, flagTo, "", "specify the target resource type to add to (e.g. list, account, etc)")
	addExe.StringVar(&addExe.listID, flagListID, "", "the ID of the list to add to")
	addExe.Var(&addExe.accountNames, flagAccountName, "the name of the account to add to the resource")
	addExe.StringVar(&addExe.content, flagContent, "", "the content of the note")

	addExe.Usage = commandUsageFunc(name, summary, addExe.FlagSet)

	return &addExe
}

func (a *AddExecutor) Execute() error {
	if a.toResourceType == "" {
		return FlagNotSetError{flagText: flagTo}
	}

	funcMap := map[string]func(*client.Client) error{
		resourceList:    a.addToList,
		resourceAccount: a.addToAccount,
	}

	doFunc, ok := funcMap[a.toResourceType]
	if !ok {
		return UnsupportedTypeError{resourceType: a.toResourceType}
	}

	gtsClient, err := client.NewClientFromConfig(a.topLevelFlags.ConfigDir)
	if err != nil {
		return fmt.Errorf("unable to create the GoToSocial client; %w", err)
	}

	return doFunc(gtsClient)
}

func (a *AddExecutor) addToList(gtsClient *client.Client) error {
	funcMap := map[string]func(*client.Client) error{
		resourceAccount: a.addAccountsToList,
	}

	doFunc, ok := funcMap[a.resourceType]
	if !ok {
		return UnsupportedAddOperationError{
			ResourceType:      a.resourceType,
			AddToResourceType: a.toResourceType,
		}
	}

	return doFunc(gtsClient)
}

func (a *AddExecutor) addAccountsToList(gtsClient *client.Client) error {
	if a.listID == "" {
		return FlagNotSetError{flagText: flagListID}
	}

	if len(a.accountNames) == 0 {
		return NoAccountSpecifiedError{}
	}

	accountIDs := make([]string, len(a.accountNames))

	for ind := range a.accountNames {
		accountID, err := getTheirAccountID(gtsClient, a.accountNames[ind])
		if err != nil {
			return fmt.Errorf("unable to get the account ID for %s, %w", a.accountNames[ind], err)
		}

		accountIDs[ind] = accountID
	}

	if err := gtsClient.AddAccountsToList(a.listID, accountIDs); err != nil {
		return fmt.Errorf("unable to add the accounts to the list; %w", err)
	}

	fmt.Println("Successfully added the account(s) to the list.")

	return nil
}

func (a *AddExecutor) addToAccount(gtsClient *client.Client) error {
	funcMap := map[string]func(*client.Client) error{
		resourceNote: a.addNoteToAccount,
	}

	doFunc, ok := funcMap[a.resourceType]
	if !ok {
		return UnsupportedAddOperationError{
			ResourceType:      a.resourceType,
			AddToResourceType: a.toResourceType,
		}
	}

	return doFunc(gtsClient)
}

func (a *AddExecutor) addNoteToAccount(gtsClient *client.Client) error {
	if len(a.accountNames) != 1 {
		return fmt.Errorf("unexpected number of accounts specified; want 1, got %d", len(a.accountNames))
	}

	accountID, err := getAccountID(gtsClient, false, a.accountNames[0], a.topLevelFlags.ConfigDir)
	if err != nil {
		return fmt.Errorf("received an error while getting the account ID; %w", err)
	}

	if a.content == "" {
		return EmptyContentError{
			ResourceType: resourceNote,
			Hint:         "please use --" + flagContent,
		}
	}

	if err := gtsClient.SetPrivateNote(accountID, a.content); err != nil {
		return fmt.Errorf("unable to add the private note to the account; %w", err)
	}

	fmt.Println("Successfully added the private note to the account.")

	return nil
}
