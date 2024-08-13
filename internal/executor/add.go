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
	if a.votes.Empty() {
		return errors.New("please use --" + flagVote + " to make a choice in this poll")
	}

	poll, err := gtsClient.GetPoll(a.pollID)
	if err != nil {
		return fmt.Errorf("unable to retrieve the poll: %w", err)
	}

	if poll.Expired {
		return PollClosedError{}
	}

	if !poll.Multiple && !a.votes.ExpectedLength(1) {
		return MultipleChoiceError{}
	}

	if err := gtsClient.VoteInPoll(a.pollID, []int(a.votes)); err != nil {
		return fmt.Errorf("unable to add your vote(s) to the poll: %w", err)
	}

	a.printer.PrintSuccess("Successfully added your vote(s) to the poll.")

	return nil
}
