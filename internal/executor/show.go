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
	limit                   int
}

func NewShowExecutor(tlf TopLevelFlags, name, summary string) *ShowExecutor {
	command := ShowExecutor{
		FlagSet:       flag.NewFlagSet(name, flag.ExitOnError),
		topLevelFlags: tlf,
	}

	command.BoolVar(&command.myAccount, flagMyAccount, false, "set to true to lookup your account")
	command.BoolVar(&command.skipAccountRelationship, flagSkipRelationship, false, "set to true to skip showing your relationship to the specified account")
	command.BoolVar(&command.showUserPreferences, flagShowPreferences, false, "show your preferences")
	command.BoolVar(&command.showInBrowser, flagBrowser, false, "set to true to view in the browser")
	command.StringVar(&command.resourceType, flagType, "", "specify the type of resource to display")
	command.StringVar(&command.accountName, flagAccountName, "", "specify the account name in full (username@domain)")
	command.StringVar(&command.statusID, flagStatusID, "", "specify the ID of the status to display")
	command.StringVar(&command.timelineCategory, flagTimelineCategory, "home", "specify the type of timeline to display (valid values are home, public, list and tag)")
	command.StringVar(&command.listID, flagListID, "", "specify the ID of the list to display")
	command.StringVar(&command.tag, flagTag, "", "specify the name of the tag to use")
	command.IntVar(&command.limit, flagLimit, 20, "specify the limit of items to display")

	command.Usage = commandUsageFunc(name, summary, command.FlagSet)

	return &command
}

func (c *ShowExecutor) Execute() error {
	if c.resourceType == "" {
		return FlagNotSetError{flagText: flagType}
	}

	funcMap := map[string]func(*client.Client) error{
		resourceInstance:  c.showInstance,
		resourceAccount:   c.showAccount,
		resourceStatus:    c.showStatus,
		resourceTimeline:  c.showTimeline,
		resourceList:      c.showList,
		resourceFollowers: c.showFollowers,
		resourceFollowing: c.showFollowing,
		resourceBlocked:   c.showBlocked,
	}

	doFunc, ok := funcMap[c.resourceType]
	if !ok {
		return UnsupportedTypeError{resourceType: c.resourceType}
	}

	gtsClient, err := client.NewClientFromConfig(c.topLevelFlags.ConfigDir)
	if err != nil {
		return fmt.Errorf("unable to create the GoToSocial client; %w", err)
	}

	return doFunc(gtsClient)
}

func (c *ShowExecutor) showInstance(gtsClient *client.Client) error {
	instance, err := gtsClient.GetInstance()
	if err != nil {
		return fmt.Errorf("unable to retrieve the instance details; %w", err)
	}

	utilities.Display(instance, *c.topLevelFlags.NoColor)

	return nil
}

func (c *ShowExecutor) showAccount(gtsClient *client.Client) error {
	var (
		account model.Account
		err     error
	)

	if c.myAccount {
		account, err = getMyAccount(gtsClient, c.topLevelFlags.ConfigDir)
		if err != nil {
			return fmt.Errorf("received an error while getting the account details; %w", err)
		}
	} else {
		if c.accountName == "" {
			return FlagNotSetError{flagText: flagAccountName}
		}

		account, err = getAccount(gtsClient, c.accountName)
		if err != nil {
			return fmt.Errorf("received an error while getting the account details; %w", err)
		}
	}

	if c.showInBrowser {
		utilities.OpenLink(account.URL)

		return nil
	}

	utilities.Display(account, *c.topLevelFlags.NoColor)

	if !c.myAccount && !c.skipAccountRelationship {
		relationship, err := gtsClient.GetAccountRelationship(account.ID)
		if err != nil {
			return fmt.Errorf("unable to retrieve the relationship to this account; %w", err)
		}

		utilities.Display(relationship, *c.topLevelFlags.NoColor)
	}

	if c.myAccount && c.showUserPreferences {
		preferences, err := gtsClient.GetUserPreferences()
		if err != nil {
			return fmt.Errorf("unable to retrieve the user preferences; %w", err)
		}

		utilities.Display(preferences, *c.topLevelFlags.NoColor)
	}

	return nil
}

func (c *ShowExecutor) showStatus(gtsClient *client.Client) error {
	if c.statusID == "" {
		return FlagNotSetError{flagText: flagStatusID}
	}

	status, err := gtsClient.GetStatus(c.statusID)
	if err != nil {
		return fmt.Errorf("unable to retrieve the status; %w", err)
	}

	if c.showInBrowser {
		utilities.OpenLink(status.URL)

		return nil
	}

	utilities.Display(status, *c.topLevelFlags.NoColor)

	return nil
}

func (c *ShowExecutor) showTimeline(gtsClient *client.Client) error {
	var (
		timeline model.Timeline
		err      error
	)

	switch c.timelineCategory {
	case "home":
		timeline, err = gtsClient.GetHomeTimeline(c.limit)
	case "public":
		timeline, err = gtsClient.GetPublicTimeline(c.limit)
	case "list":
		if c.listID == "" {
			return FlagNotSetError{flagText: flagListID}
		}

		timeline, err = gtsClient.GetListTimeline(c.listID, c.limit)
	case "tag":
		if c.tag == "" {
			return FlagNotSetError{flagText: flagTag}
		}

		timeline, err = gtsClient.GetTagTimeline(c.tag, c.limit)
	default:
		return InvalidTimelineCategoryError{category: c.timelineCategory}
	}

	if err != nil {
		return fmt.Errorf("unable to retrieve the %s timeline; %w", c.timelineCategory, err)
	}

	if len(timeline.Statuses) == 0 {
		fmt.Println("There are no statuses in this timeline.")

		return nil
	}

	utilities.Display(timeline, *c.topLevelFlags.NoColor)

	return nil
}

func (c *ShowExecutor) showList(gtsClient *client.Client) error {
	if c.listID == "" {
		return c.showLists(gtsClient)
	}

	list, err := gtsClient.GetList(c.listID)
	if err != nil {
		return fmt.Errorf("unable to retrieve the list; %w", err)
	}

	accounts, err := gtsClient.GetAccountsFromList(c.listID, 0)
	if err != nil {
		return fmt.Errorf("unable to retrieve the accounts from the list; %w", err)
	}

	if len(accounts) > 0 {
		accountMap := make(map[string]string)
		for i := range accounts {
			accountMap[accounts[i].Acct] = accounts[i].Username
		}

		list.Accounts = accountMap
	}

	utilities.Display(list, *c.topLevelFlags.NoColor)

	return nil
}

func (c *ShowExecutor) showLists(gtsClient *client.Client) error {
	lists, err := gtsClient.GetAllLists()
	if err != nil {
		return fmt.Errorf("unable to retrieve the lists; %w", err)
	}

	if len(lists) == 0 {
		fmt.Println("You have no lists.")

		return nil
	}

	utilities.Display(lists, *c.topLevelFlags.NoColor)

	return nil
}

func (c *ShowExecutor) showFollowers(gtsClient *client.Client) error {
	accountID, err := getAccountID(gtsClient, c.myAccount, c.accountName, c.topLevelFlags.ConfigDir)
	if err != nil {
		return fmt.Errorf("received an error while getting the account ID; %w", err)
	}

	followers, err := gtsClient.GetFollowers(accountID, c.limit)
	if err != nil {
		return fmt.Errorf("unable to retrieve the list of followers; %w", err)
	}

	if len(followers.Accounts) > 0 {
		utilities.Display(followers, *c.topLevelFlags.NoColor)
	} else {
		fmt.Println("There are no followers for this account or the list is hidden.")
	}

	return nil
}

func (c *ShowExecutor) showFollowing(gtsClient *client.Client) error {
	accountID, err := getAccountID(gtsClient, c.myAccount, c.accountName, c.topLevelFlags.ConfigDir)
	if err != nil {
		return fmt.Errorf("received an error while getting the account ID; %w", err)
	}

	following, err := gtsClient.GetFollowing(accountID, c.limit)
	if err != nil {
		return fmt.Errorf("unable to retrieve the list of followed accounts; %w", err)
	}

	if len(following.Accounts) > 0 {
		utilities.Display(following, *c.topLevelFlags.NoColor)
	} else {
		fmt.Println("This account is not following anyone or the list is hidden.")
	}

	return nil
}

func (c *ShowExecutor) showBlocked(gtsClient *client.Client) error {
	blocked, err := gtsClient.GetBlockedAccounts(c.limit)
	if err != nil {
		return fmt.Errorf("unable to retrieve the list of blocked accounts; %w", err)
	}

	if len(blocked.Accounts) > 0 {
		utilities.Display(blocked, *c.topLevelFlags.NoColor)
	} else {
		fmt.Println("You have no blocked accounts.")
	}

	return nil
}
