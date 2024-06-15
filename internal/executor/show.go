// SPDX-FileCopyrightText: 2024 Dan Anglin <d.n.i.anglin@gmail.com>
//
// SPDX-License-Identifier: GPL-3.0-or-later

package executor

import (
	"flag"
	"fmt"

	"codeflow.dananglin.me.uk/apollo/enbas/internal/client"
	"codeflow.dananglin.me.uk/apollo/enbas/internal/model"
	"codeflow.dananglin.me.uk/apollo/enbas/internal/utilities"
)

type ShowExecutor struct {
	*flag.FlagSet
	topLevelFlags           TopLevelFlags
	myAccount               bool
	skipAccountRelationship bool
	showUserPreferences     bool
	showInBrowser           bool
	resourceType            string
	accountName             string
	statusID                string
	timelineCategory        string
	listID                  string
	tag                     string
	pollID                  string
	limit                   int
}

func NewShowExecutor(tlf TopLevelFlags, name, summary string) *ShowExecutor {
	showExe := ShowExecutor{
		FlagSet:       flag.NewFlagSet(name, flag.ExitOnError),
		topLevelFlags: tlf,
	}

	showExe.BoolVar(&showExe.myAccount, flagMyAccount, false, "Set to true to lookup your account")
	showExe.BoolVar(&showExe.skipAccountRelationship, flagSkipRelationship, false, "Set to true to skip showing your relationship to the specified account")
	showExe.BoolVar(&showExe.showUserPreferences, flagShowPreferences, false, "Show your preferences")
	showExe.BoolVar(&showExe.showInBrowser, flagBrowser, false, "Set to true to view in the browser")
	showExe.StringVar(&showExe.resourceType, flagType, "", "Specify the type of resource to display")
	showExe.StringVar(&showExe.accountName, flagAccountName, "", "Specify the account name in full (username@domain)")
	showExe.StringVar(&showExe.statusID, flagStatusID, "", "Specify the ID of the status to display")
	showExe.StringVar(&showExe.timelineCategory, flagTimelineCategory, model.TimelineCategoryHome, "Specify the timeline category to view")
	showExe.StringVar(&showExe.listID, flagListID, "", "Specify the ID of the list to display")
	showExe.StringVar(&showExe.tag, flagTag, "", "Specify the name of the tag to use")
	showExe.StringVar(&showExe.pollID, flagPollID, "", "Specify the ID of the poll to display")
	showExe.IntVar(&showExe.limit, flagLimit, 20, "Specify the limit of items to display")

	showExe.Usage = commandUsageFunc(name, summary, showExe.FlagSet)

	return &showExe
}

func (s *ShowExecutor) Execute() error {
	if s.resourceType == "" {
		return FlagNotSetError{flagText: flagType}
	}

	funcMap := map[string]func(*client.Client) error{
		resourceInstance:      s.showInstance,
		resourceAccount:       s.showAccount,
		resourceStatus:        s.showStatus,
		resourceTimeline:      s.showTimeline,
		resourceList:          s.showList,
		resourceFollowers:     s.showFollowers,
		resourceFollowing:     s.showFollowing,
		resourceBlocked:       s.showBlocked,
		resourceBookmarks:     s.showBookmarks,
		resourceLiked:         s.showLiked,
		resourceStarred:       s.showLiked,
		resourceFollowRequest: s.showFollowRequests,
		resourcePoll:          s.showPoll,
	}

	doFunc, ok := funcMap[s.resourceType]
	if !ok {
		return UnsupportedTypeError{resourceType: s.resourceType}
	}

	gtsClient, err := client.NewClientFromConfig(s.topLevelFlags.ConfigDir)
	if err != nil {
		return fmt.Errorf("unable to create the GoToSocial client: %w", err)
	}

	return doFunc(gtsClient)
}

func (s *ShowExecutor) showInstance(gtsClient *client.Client) error {
	instance, err := gtsClient.GetInstance()
	if err != nil {
		return fmt.Errorf("unable to retrieve the instance details: %w", err)
	}

	utilities.Display(instance, *s.topLevelFlags.NoColor, s.topLevelFlags.Pager)

	return nil
}

func (s *ShowExecutor) showAccount(gtsClient *client.Client) error {
	var (
		account model.Account
		err     error
	)

	if s.myAccount {
		account, err = getMyAccount(gtsClient, s.topLevelFlags.ConfigDir)
		if err != nil {
			return fmt.Errorf("received an error while getting the account details: %w", err)
		}
	} else {
		if s.accountName == "" {
			return FlagNotSetError{flagText: flagAccountName}
		}

		account, err = getAccount(gtsClient, s.accountName)
		if err != nil {
			return fmt.Errorf("received an error while getting the account details: %w", err)
		}
	}

	if s.showInBrowser {
		utilities.OpenLink(account.URL)

		return nil
	}

	utilities.Display(account, *s.topLevelFlags.NoColor, s.topLevelFlags.Pager)

	if !s.myAccount && !s.skipAccountRelationship {
		relationship, err := gtsClient.GetAccountRelationship(account.ID)
		if err != nil {
			return fmt.Errorf("unable to retrieve the relationship to this account: %w", err)
		}

		utilities.Display(relationship, *s.topLevelFlags.NoColor, s.topLevelFlags.Pager)
	}

	if s.myAccount && s.showUserPreferences {
		preferences, err := gtsClient.GetUserPreferences()
		if err != nil {
			return fmt.Errorf("unable to retrieve the user preferences: %w", err)
		}

		utilities.Display(preferences, *s.topLevelFlags.NoColor, s.topLevelFlags.Pager)
	}

	return nil
}

func (s *ShowExecutor) showStatus(gtsClient *client.Client) error {
	if s.statusID == "" {
		return FlagNotSetError{flagText: flagStatusID}
	}

	status, err := gtsClient.GetStatus(s.statusID)
	if err != nil {
		return fmt.Errorf("unable to retrieve the status: %w", err)
	}

	if s.showInBrowser {
		utilities.OpenLink(status.URL)

		return nil
	}

	utilities.Display(status, *s.topLevelFlags.NoColor, s.topLevelFlags.Pager)

	return nil
}

func (s *ShowExecutor) showTimeline(gtsClient *client.Client) error {
	var (
		timeline model.StatusList
		err      error
	)

	switch s.timelineCategory {
	case model.TimelineCategoryHome:
		timeline, err = gtsClient.GetHomeTimeline(s.limit)
	case model.TimelineCategoryPublic:
		timeline, err = gtsClient.GetPublicTimeline(s.limit)
	case model.TimelineCategoryList:
		if s.listID == "" {
			return FlagNotSetError{flagText: flagListID}
		}

		var list model.List

		list, err = gtsClient.GetList(s.listID)
		if err != nil {
			return fmt.Errorf("unable to retrieve the list: %w", err)
		}

		timeline, err = gtsClient.GetListTimeline(list.ID, list.Title, s.limit)
	case model.TimelineCategoryTag:
		if s.tag == "" {
			return FlagNotSetError{flagText: flagTag}
		}

		timeline, err = gtsClient.GetTagTimeline(s.tag, s.limit)
	default:
		return model.InvalidTimelineCategoryError{Value: s.timelineCategory}
	}

	if err != nil {
		return fmt.Errorf("unable to retrieve the %s timeline: %w", s.timelineCategory, err)
	}

	if len(timeline.Statuses) == 0 {
		fmt.Println("There are no statuses in this timeline.")

		return nil
	}

	utilities.Display(timeline, *s.topLevelFlags.NoColor, s.topLevelFlags.Pager)

	return nil
}

func (s *ShowExecutor) showList(gtsClient *client.Client) error {
	if s.listID == "" {
		return s.showLists(gtsClient)
	}

	list, err := gtsClient.GetList(s.listID)
	if err != nil {
		return fmt.Errorf("unable to retrieve the list: %w", err)
	}

	accounts, err := gtsClient.GetAccountsFromList(s.listID, 0)
	if err != nil {
		return fmt.Errorf("unable to retrieve the accounts from the list: %w", err)
	}

	if len(accounts) > 0 {
		accountMap := make(map[string]string)
		for i := range accounts {
			accountMap[accounts[i].Acct] = accounts[i].Username
		}

		list.Accounts = accountMap
	}

	utilities.Display(list, *s.topLevelFlags.NoColor, s.topLevelFlags.Pager)

	return nil
}

func (s *ShowExecutor) showLists(gtsClient *client.Client) error {
	lists, err := gtsClient.GetAllLists()
	if err != nil {
		return fmt.Errorf("unable to retrieve the lists: %w", err)
	}

	if len(lists) == 0 {
		fmt.Println("You have no lists.")

		return nil
	}

	utilities.Display(lists, *s.topLevelFlags.NoColor, s.topLevelFlags.Pager)

	return nil
}

func (s *ShowExecutor) showFollowers(gtsClient *client.Client) error {
	accountID, err := getAccountID(gtsClient, s.myAccount, s.accountName, s.topLevelFlags.ConfigDir)
	if err != nil {
		return fmt.Errorf("received an error while getting the account ID: %w", err)
	}

	followers, err := gtsClient.GetFollowers(accountID, s.limit)
	if err != nil {
		return fmt.Errorf("unable to retrieve the list of followers: %w", err)
	}

	if len(followers.Accounts) > 0 {
		utilities.Display(followers, *s.topLevelFlags.NoColor, s.topLevelFlags.Pager)
	} else {
		fmt.Println("There are no followers for this account or the list is hidden.")
	}

	return nil
}

func (s *ShowExecutor) showFollowing(gtsClient *client.Client) error {
	accountID, err := getAccountID(gtsClient, s.myAccount, s.accountName, s.topLevelFlags.ConfigDir)
	if err != nil {
		return fmt.Errorf("received an error while getting the account ID: %w", err)
	}

	following, err := gtsClient.GetFollowing(accountID, s.limit)
	if err != nil {
		return fmt.Errorf("unable to retrieve the list of followed accounts: %w", err)
	}

	if len(following.Accounts) > 0 {
		utilities.Display(following, *s.topLevelFlags.NoColor, s.topLevelFlags.Pager)
	} else {
		fmt.Println("This account is not following anyone or the list is hidden.")
	}

	return nil
}

func (s *ShowExecutor) showBlocked(gtsClient *client.Client) error {
	blocked, err := gtsClient.GetBlockedAccounts(s.limit)
	if err != nil {
		return fmt.Errorf("unable to retrieve the list of blocked accounts: %w", err)
	}

	if len(blocked.Accounts) > 0 {
		utilities.Display(blocked, *s.topLevelFlags.NoColor, s.topLevelFlags.Pager)
	} else {
		fmt.Println("You have no blocked accounts.")
	}

	return nil
}

func (s *ShowExecutor) showBookmarks(gtsClient *client.Client) error {
	bookmarks, err := gtsClient.GetBookmarks(s.limit)
	if err != nil {
		return fmt.Errorf("unable to retrieve the list of bookmarks: %w", err)
	}

	if len(bookmarks.Statuses) > 0 {
		utilities.Display(bookmarks, *s.topLevelFlags.NoColor, s.topLevelFlags.Pager)
	} else {
		fmt.Println("You have no bookmarks.")
	}

	return nil
}

func (s *ShowExecutor) showLiked(gtsClient *client.Client) error {
	liked, err := gtsClient.GetLikedStatuses(s.limit, s.resourceType)
	if err != nil {
		return fmt.Errorf("unable to retrieve the list of your %s statuses: %w", s.resourceType, err)
	}

	if len(liked.Statuses) > 0 {
		utilities.Display(liked, *s.topLevelFlags.NoColor, s.topLevelFlags.Pager)
	} else {
		fmt.Printf("You have no %s statuses.\n", s.resourceType)
	}

	return nil
}

func (s *ShowExecutor) showFollowRequests(gtsClient *client.Client) error {
	accounts, err := gtsClient.GetFollowRequests(s.limit)
	if err != nil {
		return fmt.Errorf("unable to retrieve the list of follow requests: %w", err)
	}

	if len(accounts.Accounts) > 0 {
		utilities.Display(accounts, *s.topLevelFlags.NoColor, s.topLevelFlags.Pager)
	} else {
		fmt.Println("You have no follow requests.")
	}

	return nil
}

func (s *ShowExecutor) showPoll(gtsClient *client.Client) error {
	if s.pollID == "" {
		return FlagNotSetError{flagText: flagPollID}
	}

	poll, err := gtsClient.GetPoll(s.pollID)
	if err != nil {
		return fmt.Errorf("unable to retrieve the poll: %w", err)
	}

	utilities.Display(poll, *s.topLevelFlags.NoColor, s.topLevelFlags.Pager)

	return nil
}
