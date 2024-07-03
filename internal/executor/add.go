// SPDX-FileCopyrightText: 2024 Dan Anglin <d.n.i.anglin@gmail.com>
//
// SPDX-License-Identifier: GPL-3.0-or-later

package executor

import (
	"errors"
	"flag"
	"fmt"

	"codeflow.dananglin.me.uk/apollo/enbas/internal/client"
	"codeflow.dananglin.me.uk/apollo/enbas/internal/config"
	"codeflow.dananglin.me.uk/apollo/enbas/internal/printer"
)

type AddExecutor struct {
	*flag.FlagSet

	printer        *printer.Printer
	config         *config.Config
	resourceType   string
	toResourceType string
	listID         string
	statusID       string
	pollID         string
	choices        MultiIntFlagValue
	accountNames   MultiStringFlagValue
	content        string
}

func NewAddExecutor(printer *printer.Printer, config *config.Config, name, summary string) *AddExecutor {
	emptyArr := make([]string, 0, 3)

	addExe := AddExecutor{
		FlagSet: flag.NewFlagSet(name, flag.ExitOnError),

		printer:      printer,
		config:       config,
		accountNames: MultiStringFlagValue(emptyArr),
	}

	addExe.StringVar(&addExe.resourceType, flagType, "", "Specify the resource type to add (e.g. account, note)")
	addExe.StringVar(&addExe.toResourceType, flagTo, "", "Specify the target resource type to add to (e.g. list, account, etc)")
	addExe.StringVar(&addExe.listID, flagListID, "", "The ID of the list")
	addExe.StringVar(&addExe.statusID, flagStatusID, "", "The ID of the status")
	addExe.StringVar(&addExe.content, flagContent, "", "The content of the resource")
	addExe.StringVar(&addExe.pollID, flagPollID, "", "The ID of the poll")
	addExe.Var(&addExe.accountNames, flagAccountName, "The name of the account")
	addExe.Var(&addExe.choices, flagVote, "Add a vote to an option in a poll")

	addExe.Usage = commandUsageFunc(name, summary, addExe.FlagSet)

	return &addExe
}

func (a *AddExecutor) Execute() error {
	if a.toResourceType == "" {
		return FlagNotSetError{flagText: flagTo}
	}

	funcMap := map[string]func(*client.Client) error{
		resourceList:      a.addToList,
		resourceAccount:   a.addToAccount,
		resourceBookmarks: a.addToBookmarks,
		resourceStatus:    a.addToStatus,
		resourcePoll:      a.addToPoll,
	}

	doFunc, ok := funcMap[a.toResourceType]
	if !ok {
		return UnsupportedTypeError{resourceType: a.toResourceType}
	}

	gtsClient, err := client.NewClientFromFile(a.config.CredentialsFile)
	if err != nil {
		return fmt.Errorf("unable to create the GoToSocial client: %w", err)
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
			return fmt.Errorf("unable to get the account ID for %s: %w", a.accountNames[ind], err)
		}

		relationship, err := gtsClient.GetAccountRelationship(accountID)
		if err != nil {
			return fmt.Errorf("unable to get your relationship to %s: %w", a.accountNames[ind], err)
		}

		if !relationship.Following {
			return NotFollowingError{Account: a.accountNames[ind]}
		}

		accountIDs[ind] = accountID
	}

	if err := gtsClient.AddAccountsToList(a.listID, accountIDs); err != nil {
		return fmt.Errorf("unable to add the accounts to the list: %w", err)
	}

	a.printer.PrintSuccess("Successfully added the account(s) to the list.")

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
		return fmt.Errorf("unexpected number of accounts specified: want 1, got %d", len(a.accountNames))
	}

	accountID, err := getAccountID(gtsClient, false, a.accountNames[0], a.config.CredentialsFile)
	if err != nil {
		return fmt.Errorf("received an error while getting the account ID: %w", err)
	}

	if a.content == "" {
		return EmptyContentError{
			ResourceType: resourceNote,
			Hint:         "please use --" + flagContent,
		}
	}

	if err := gtsClient.SetPrivateNote(accountID, a.content); err != nil {
		return fmt.Errorf("unable to add the private note to the account: %w", err)
	}

	a.printer.PrintSuccess("Successfully added the private note to the account.")

	return nil
}

func (a *AddExecutor) addToBookmarks(gtsClient *client.Client) error {
	funcMap := map[string]func(*client.Client) error{
		resourceStatus: a.addStatusToBookmarks,
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

func (a *AddExecutor) addStatusToBookmarks(gtsClient *client.Client) error {
	if a.statusID == "" {
		return FlagNotSetError{flagText: flagStatusID}
	}

	if err := gtsClient.AddStatusToBookmarks(a.statusID); err != nil {
		return fmt.Errorf("unable to add the status to your bookmarks: %w", err)
	}

	a.printer.PrintSuccess("Successfully added the status to your bookmarks.")

	return nil
}

func (a *AddExecutor) addToStatus(gtsClient *client.Client) error {
	if a.statusID == "" {
		return FlagNotSetError{flagText: flagStatusID}
	}

	funcMap := map[string]func(*client.Client) error{
		resourceStar:  a.addStarToStatus,
		resourceLike:  a.addStarToStatus,
		resourceBoost: a.addBoostToStatus,
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

func (a *AddExecutor) addStarToStatus(gtsClient *client.Client) error {
	if err := gtsClient.LikeStatus(a.statusID); err != nil {
		return fmt.Errorf("unable to add the %s to the status: %w", a.resourceType, err)
	}

	a.printer.PrintSuccess("Successfully added a " + a.resourceType + " to the status.")

	return nil
}

func (a *AddExecutor) addBoostToStatus(gtsClient *client.Client) error {
	if err := gtsClient.ReblogStatus(a.statusID); err != nil {
		return fmt.Errorf("unable to add the boost to the status: %w", err)
	}

	a.printer.PrintSuccess("Successfully added the boost to the status.")

	return nil
}

func (a *AddExecutor) addToPoll(gtsClient *client.Client) error {
	if a.pollID == "" {
		return FlagNotSetError{flagText: flagPollID}
	}

	funcMap := map[string]func(*client.Client) error{
		resourceVote: a.addVoteToPoll,
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

func (a *AddExecutor) addVoteToPoll(gtsClient *client.Client) error {
	if len(a.choices) == 0 {
		return errors.New("please use --" + flagVote + " to make a choice in this poll")
	}

	poll, err := gtsClient.GetPoll(a.pollID)
	if err != nil {
		return fmt.Errorf("unable to retrieve the poll: %w", err)
	}

	if poll.Expired {
		return PollClosedError{}
	}

	if !poll.Multiple && len(a.choices) > 1 {
		return MultipleChoiceError{}
	}

	if err := gtsClient.VoteInPoll(a.pollID, []int(a.choices)); err != nil {
		return fmt.Errorf("unable to add your vote(s) to the poll: %w", err)
	}

	a.printer.PrintSuccess("Successfully added your vote(s) to the poll.")

	return nil
}
