package main

import (
	"flag"
	"fmt"

	"codeflow.dananglin.me.uk/apollo/enbas/internal/client"
	"codeflow.dananglin.me.uk/apollo/enbas/internal/model"
	"codeflow.dananglin.me.uk/apollo/enbas/internal/utilities"
)

type showCommand struct {
	*flag.FlagSet
	myAccount               bool
	showAccountRelationship bool
	showUserPreferences     bool
	resourceType            string
	account                 string
	accountID               string
	statusID                string
	timelineCategory        string
	listID                  string
	tag                     string
	limit                   int
}

func newShowCommand(name, summary string) *showCommand {
	command := showCommand{
		FlagSet: flag.NewFlagSet(name, flag.ExitOnError),
	}

	command.BoolVar(&command.myAccount, myAccountFlag, false, "set to true to lookup your account")
	command.BoolVar(&command.showAccountRelationship, showAccountRelationshipFlag, false, "show your relationship to the specified account")
	command.BoolVar(&command.showUserPreferences, showUserPreferencesFlag, false, "show your preferences")
	command.StringVar(&command.resourceType, resourceTypeFlag, "", "specify the type of resource to display")
	command.StringVar(&command.account, accountFlag, "", "specify the account URI to lookup")
	command.StringVar(&command.accountID, accountIDFlag, "", "specify the account ID")
	command.StringVar(&command.statusID, statusIDFlag, "", "specify the ID of the status to display")
	command.StringVar(&command.timelineCategory, timelineCategoryFlag, "home", "specify the type of timeline to display (valid values are home, public, list and tag)")
	command.StringVar(&command.listID, listIDFlag, "", "specify the ID of the list to display")
	command.StringVar(&command.tag, tagFlag, "", "specify the name of the tag to use")
	command.IntVar(&command.limit, limitFlag, 20, "specify the limit of items to display")

	command.Usage = commandUsageFunc(name, summary, command.FlagSet)

	return &command
}

func (c *showCommand) Execute() error {
	if c.resourceType == "" {
		return flagNotSetError{flagText: resourceTypeFlag}
	}

	funcMap := map[string]func(*client.Client) error{
		instanceResource:  c.showInstance,
		accountResource:   c.showAccount,
		statusResource:    c.showStatus,
		timelineResource:  c.showTimeline,
		listResource:      c.showList,
		followersResource: c.showFollowers,
		followingResource: c.showFollowing,
	}

	doFunc, ok := funcMap[c.resourceType]
	if !ok {
		return unsupportedResourceTypeError{resourceType: c.resourceType}
	}

	gtsClient, err := client.NewClientFromConfig()
	if err != nil {
		return fmt.Errorf("unable to create the GoToSocial client; %w", err)
	}

	return doFunc(gtsClient)
}

func (c *showCommand) showInstance(gts *client.Client) error {
	instance, err := gts.GetInstance()
	if err != nil {
		return fmt.Errorf("unable to retrieve the instance details; %w", err)
	}

	fmt.Println(instance)

	return nil
}

func (c *showCommand) showAccount(gts *client.Client) error {
	var (
		account model.Account
		err     error
	)

	if c.myAccount {
		account, err = getMyAccount(gts)
		if err != nil {
			return fmt.Errorf("received an error while getting account details; %w", err)
		}
	} else {
		if c.account == "" {
			return flagNotSetError{flagText: accountFlag}
		}

		accountURI := c.account

		account, err = gts.GetAccount(accountURI)
		if err != nil {
			return fmt.Errorf("unable to retrieve the account details; %w", err)
		}
	}

	fmt.Println(account)

	if !c.myAccount && c.showAccountRelationship {
		relationship, err := gts.GetAccountRelationship(account.ID)
		if err != nil {
			return fmt.Errorf("unable to retrieve the relationship to this account; %w", err)
		}

		fmt.Println(relationship)
	}

	if c.myAccount && c.showUserPreferences {
		preferences, err := gts.GetUserPreferences()
		if err != nil {
			return fmt.Errorf("unable to retrieve the user preferences; %w", err)
		}

		fmt.Println(preferences)
	}

	return nil
}

func (c *showCommand) showStatus(gts *client.Client) error {
	if c.statusID == "" {
		return flagNotSetError{flagText: statusIDFlag}
	}

	status, err := gts.GetStatus(c.statusID)
	if err != nil {
		return fmt.Errorf("unable to retrieve the status; %w", err)
	}

	fmt.Println(status)

	return nil
}

func (c *showCommand) showTimeline(gts *client.Client) error {
	var (
		timeline model.Timeline
		err      error
	)

	switch c.timelineCategory {
	case "home":
		timeline, err = gts.GetHomeTimeline(c.limit)
	case "public":
		timeline, err = gts.GetPublicTimeline(c.limit)
	case "list":
		if c.listID == "" {
			return flagNotSetError{flagText: listIDFlag}
		}

		timeline, err = gts.GetListTimeline(c.listID, c.limit)
	case "tag":
		if c.tag == "" {
			return flagNotSetError{flagText: tagFlag}
		}

		timeline, err = gts.GetTagTimeline(c.tag, c.limit)
	default:
		return invalidTimelineCategoryError{category: c.timelineCategory}
	}

	if err != nil {
		return fmt.Errorf("unable to retrieve the %s timeline; %w", c.timelineCategory, err)
	}

	if len(timeline.Statuses) == 0 {
		fmt.Println("There are no statuses in this timeline.")

		return nil
	}

	fmt.Println(timeline)

	return nil
}

func (c *showCommand) showList(gts *client.Client) error {
	if c.listID == "" {
		return c.showLists(gts)
	}

	list, err := gts.GetList(c.listID)
	if err != nil {
		return fmt.Errorf("unable to retrieve the list; %w", err)
	}

	accounts, err := gts.GetAccountsFromList(c.listID, 0)
	if err != nil {
		return fmt.Errorf("unable to retrieve the accounts from the list; %w", err)
	}

	if len(accounts) > 0 {
		accountMap := make(map[string]string)
		for i := range accounts {
			accountMap[accounts[i].ID] = accounts[i].Username
		}

		list.Accounts = accountMap
	}

	fmt.Println(list)

	return nil
}

func (c *showCommand) showLists(gts *client.Client) error {
	lists, err := gts.GetAllLists()
	if err != nil {
		return fmt.Errorf("unable to retrieve the lists; %w", err)
	}

	if len(lists) == 0 {
		fmt.Println("You have no lists.")

		return nil
	}

	fmt.Println(utilities.HeaderFormat("LISTS"))
	fmt.Println(lists)

	return nil
}

func (c *showCommand) showFollowers(gts *client.Client) error {
	var accountID string

	if c.myAccount {
		account, err := getMyAccount(gts)
		if err != nil {
			return fmt.Errorf("received an error while getting account details; %w", err)
		}

		accountID = account.ID
	} else {
		if c.accountID == "" {
			return flagNotSetError{flagText: accountIDFlag}
		}

		accountID = c.accountID
	}

	followers, err := gts.GetFollowers(accountID, c.limit)
	if err != nil {
		return fmt.Errorf("unable to retrieve the list of followers; %w", err)
	}

	if len(followers) > 0 {
		fmt.Println(followers)
	} else {
		fmt.Println("There are no followers for this account or the list is hidden.")
	}

	return nil
}

func (c *showCommand) showFollowing(gts *client.Client) error {
	var accountID string

	if c.myAccount {
		account, err := getMyAccount(gts)
		if err != nil {
			return fmt.Errorf("received an error while getting account details; %w", err)
		}

		accountID = account.ID
	} else {
		if c.accountID == "" {
			return flagNotSetError{flagText: accountIDFlag}
		}

		accountID = c.accountID
	}

	following, err := gts.GetFollowing(accountID, c.limit)
	if err != nil {
		return fmt.Errorf("unable to retrieve the list of followed accounts; %w", err)
	}

	if len(following) > 0 {
		fmt.Println(following)
	} else {
		fmt.Println("This account is not following anyone or the list is hidden.")
	}

	return nil
}
