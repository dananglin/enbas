package executor

import (
	"fmt"
	"net/rpc"

	"codeflow.dananglin.me.uk/apollo/enbas/internal/gtsclient"
	"codeflow.dananglin.me.uk/apollo/enbas/internal/model"
	"codeflow.dananglin.me.uk/apollo/enbas/internal/printer"
	"codeflow.dananglin.me.uk/apollo/enbas/internal/server"
)

func (a *AddExecutor) Execute() error {
	if a.toResourceType == "" {
		return FlagNotSetError{flagText: flagTo}
	}

	funcMap := map[string]func(*rpc.Client) error{
		resourceList:      a.addToList,
		resourceAccount:   a.addToAccount,
		resourceBookmarks: a.addToBookmarks,
		resourceStatus:    a.addToStatus,
	}

	doFunc, ok := funcMap[a.toResourceType]
	if !ok {
		return UnsupportedTypeError{resourceType: a.toResourceType}
	}

	client, err := server.Connect(a.config.Server, a.configDir)
	if err != nil {
		return fmt.Errorf("error creating the client for the daemon process: %w", err)
	}
	defer client.Close()

	return doFunc(client)
}

func (a *AddExecutor) addToList(client *rpc.Client) error {
	if a.listID == "" {
		return MissingIDError{
			resource: resourceList,
			action:   "add to",
		}
	}

	funcMap := map[string]func(*rpc.Client) error{
		resourceAccount: a.addAccountsToList,
	}

	doFunc, ok := funcMap[a.resourceType]
	if !ok {
		return UnsupportedAddOperationError{
			resourceType:      a.resourceType,
			addToResourceType: a.toResourceType,
		}
	}

	return doFunc(client)
}

func (a *AddExecutor) addAccountsToList(client *rpc.Client) error {
	if a.accountNames.Empty() {
		return NoAccountSpecifiedError{}
	}

	accounts, err := getOtherAccounts(client, a.accountNames)
	if err != nil {
		return fmt.Errorf("unable to get the accounts: %w", err)
	}

	accountIDs := make([]string, len(accounts))

	for ind := range accounts {
		var relationship model.AccountRelationship
		if err := client.Call("GTSClient.GetAccountRelationship", accounts[ind].ID, &relationship); err != nil {
			return fmt.Errorf("unable to get your relationship to %s: %w", accounts[ind].Acct, err)
		}

		if !relationship.Following {
			return NotFollowingError{account: accounts[ind].Acct}
		}

		accountIDs[ind] = accounts[ind].ID
	}

	if err := client.Call(
		"GTSClient.AddAccountsToList",
		gtsclient.AddAccountsToListArgs{
			ListID:     a.listID,
			AccountIDs: accountIDs,
		},
		nil,
	); err != nil {
		return fmt.Errorf("unable to add the accounts to the list: %w", err)
	}

	printer.PrintSuccess(a.printSettings, "Successfully added the account(s) to the list.")

	return nil
}

func (a *AddExecutor) addToAccount(client *rpc.Client) error {
	funcMap := map[string]func(*rpc.Client) error{
		resourceNote: a.addNoteToAccount,
	}

	doFunc, ok := funcMap[a.resourceType]
	if !ok {
		return UnsupportedAddOperationError{
			resourceType:      a.resourceType,
			addToResourceType: a.toResourceType,
		}
	}

	return doFunc(client)
}

func (a *AddExecutor) addNoteToAccount(client *rpc.Client) error {
	accountID, err := getAccountID(client, false, a.accountNames)
	if err != nil {
		return fmt.Errorf("received an error while getting the account ID: %w", err)
	}

	if a.content == "" {
		return Error{"please add content to the note you want to add"}
	}

	if err := client.Call(
		"GTSClient.SetPrivateNote",
		gtsclient.SetPrivateNoteArgs{
			AccountID: accountID,
			Note:      a.content,
		},
		nil,
	); err != nil {
		return fmt.Errorf("unable to add the private note to the account: %w", err)
	}

	printer.PrintSuccess(a.printSettings, "Successfully added the private note to the account.")

	return nil
}

func (a *AddExecutor) addToBookmarks(client *rpc.Client) error {
	funcMap := map[string]func(*rpc.Client) error{
		resourceStatus: a.addStatusToBookmarks,
	}

	doFunc, ok := funcMap[a.resourceType]
	if !ok {
		return UnsupportedAddOperationError{
			resourceType:      a.resourceType,
			addToResourceType: a.toResourceType,
		}
	}

	return doFunc(client)
}

func (a *AddExecutor) addStatusToBookmarks(client *rpc.Client) error {
	if a.statusID == "" {
		return MissingIDError{
			resource: resourceStatus,
			action:   "add to your bookmarks",
		}
	}

	if err := client.Call("GTSClient.AddStatusToBookmarks", a.statusID, nil); err != nil {
		return fmt.Errorf("unable to add the status to your bookmarks: %w", err)
	}

	printer.PrintSuccess(a.printSettings, "Successfully added the status to your bookmarks.")

	return nil
}

func (a *AddExecutor) addToStatus(client *rpc.Client) error {
	if a.statusID == "" {
		return MissingIDError{
			resource: resourceStatus,
			action:   "add to",
		}
	}

	funcMap := map[string]func(*rpc.Client) error{
		resourceStar:  a.addStarToStatus,
		resourceLike:  a.addStarToStatus,
		resourceBoost: a.addBoostToStatus,
		resourceVote:  a.addVoteToStatus,
	}

	doFunc, ok := funcMap[a.resourceType]
	if !ok {
		return UnsupportedAddOperationError{
			resourceType:      a.resourceType,
			addToResourceType: a.toResourceType,
		}
	}

	return doFunc(client)
}

func (a *AddExecutor) addStarToStatus(client *rpc.Client) error {
	if err := client.Call("GTSClient.LikeStatus", a.statusID, nil); err != nil {
		return fmt.Errorf("error adding the %s to the status: %w", a.resourceType, err)
	}

	printer.PrintSuccess(a.printSettings, "Successfully added a "+a.resourceType+" to the status.")

	return nil
}

func (a *AddExecutor) addBoostToStatus(client *rpc.Client) error {
	if err := client.Call("GTSClient.ReblogStatus", a.statusID, nil); err != nil {
		return fmt.Errorf("unable to add the boost to the status: %w", err)
	}

	printer.PrintSuccess(a.printSettings, "Successfully added the boost to the status.")

	return nil
}

func (a *AddExecutor) addVoteToStatus(client *rpc.Client) error {
	if a.votes.Empty() {
		return Error{"please add your vote(s) to the poll"}
	}

	var status model.Status
	if err := client.Call("GTSClient.GetStatus", a.statusID, &status); err != nil {
		return fmt.Errorf("unable to get the status: %w", err)
	}

	if status.Poll == nil {
		return Error{"this status does not have a poll"}
	}

	if status.Poll.Expired {
		return Error{"this poll is closed"}
	}

	if !status.Poll.Multiple && !a.votes.ExpectedLength(1) {
		return Error{"this poll does not allow multiple choices"}
	}

	myAccountID, err := getAccountID(client, true, nil)
	if err != nil {
		return fmt.Errorf("unable to get your account ID: %w", err)
	}

	if status.Account.ID == myAccountID {
		return Error{"you cannot vote in your own poll"}
	}

	pollID := status.Poll.ID

	if err := client.Call(
		"GTSClient.VoteInPoll",
		gtsclient.VoteInPollArgs{
			PollID:  pollID,
			Choices: a.votes,
		},
		nil,
	); err != nil {
		return fmt.Errorf("unable to add your vote(s) to the poll: %w", err)
	}

	printer.PrintSuccess(a.printSettings, "Successfully added your vote(s) to the poll.")

	return nil
}
