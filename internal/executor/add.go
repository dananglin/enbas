package executor

import (
	"errors"
	"fmt"

	"codeflow.dananglin.me.uk/apollo/enbas/internal/client"
)

func (a *AddExecutor) Execute() error {
	if a.toResourceType == "" {
		return FlagNotSetError{flagText: flagTo}
	}

	funcMap := map[string]func(*client.Client) error{
		resourceList:      a.addToList,
		resourceAccount:   a.addToAccount,
		resourceBookmarks: a.addToBookmarks,
		resourceStatus:    a.addToStatus,
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

	if a.accountNames.Empty() {
		return NoAccountSpecifiedError{}
	}

	accounts, err := getOtherAccounts(gtsClient, a.accountNames)
	if err != nil {
		return fmt.Errorf("unable to get the accounts: %w", err)
	}

	accountIDs := make([]string, len(accounts))

	for ind := range accounts {
		relationship, err := gtsClient.GetAccountRelationship(accounts[ind].ID)
		if err != nil {
			return fmt.Errorf("unable to get your relationship to %s: %w", accounts[ind].Acct, err)
		}

		if !relationship.Following {
			return NotFollowingError{Account: accounts[ind].Acct}
		}

		accountIDs[ind] = accounts[ind].ID
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
	accountID, err := getAccountID(gtsClient, false, a.accountNames)
	if err != nil {
		return fmt.Errorf("received an error while getting the account ID: %w", err)
	}

	if a.content == "" {
		return errors.New("please add content to the status that you want to create")
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
		resourceVote:  a.addVoteToStatus,
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

func (a *AddExecutor) addVoteToStatus(gtsClient *client.Client) error {
	if a.votes.Empty() {
		return NoVotesError{}
	}

	status, err := gtsClient.GetStatus(a.statusID)
	if err != nil {
		return fmt.Errorf("unable to get the status: %w", err)
	}

	if status.Poll == nil {
		return NoPollInStatusError{}
	}

	if status.Poll.Expired {
		return PollClosedError{}
	}

	if !status.Poll.Multiple && !a.votes.ExpectedLength(1) {
		return MultipleChoiceError{}
	}

	myAccountID, err := getAccountID(gtsClient, true, nil)
	if err != nil {
		return fmt.Errorf("unable to get your account ID: %w", err)
	}

	if status.Account.ID == myAccountID {
		return PollOwnerVoteError{}
	}

	pollID := status.Poll.ID

	if err := gtsClient.VoteInPoll(pollID, []int(a.votes)); err != nil {
		return fmt.Errorf("unable to add your vote(s) to the poll: %w", err)
	}

	a.printer.PrintSuccess("Successfully added your vote(s) to the poll.")

	return nil
}
