package main

import (
	"flag"
	"fmt"

	"codeflow.dananglin.me.uk/apollo/enbas/internal/client"
	"codeflow.dananglin.me.uk/apollo/enbas/internal/config"
	"codeflow.dananglin.me.uk/apollo/enbas/internal/model"
	"codeflow.dananglin.me.uk/apollo/enbas/internal/utilities"
)

type showCommand struct {
	*flag.FlagSet
	myAccount        bool
	resourceType     string
	account          string
	statusID         string
	timelineCategory string
	listID           string
	tag              string
	timelineLimit    int
}

func newShowCommand(name, summary string) *showCommand {
	command := showCommand{
		FlagSet: flag.NewFlagSet(name, flag.ExitOnError),
	}

	command.BoolVar(&command.myAccount, myAccountFlag, false, "set to true to lookup your account")
	command.StringVar(&command.resourceType, resourceTypeFlag, "", "specify the type of resource to display")
	command.StringVar(&command.account, accountFlag, "", "specify the account URI to lookup")
	command.StringVar(&command.statusID, statusIDFlag, "", "specify the ID of the status to display")
	command.StringVar(&command.timelineCategory, timelineCategoryFlag, "home", "specify the type of timeline to display (valid values are home, public, list and tag)")
	command.StringVar(&command.listID, listIDFlag, "", "specify the ID of the list to display")
	command.StringVar(&command.tag, tagFlag, "", "specify the name of the tag to use")
	command.IntVar(&command.timelineLimit, timelineLimitFlag, 5, "specify the number of statuses to display")

	command.Usage = commandUsageFunc(name, summary, command.FlagSet)

	return &command
}

func (c *showCommand) Execute() error {
	if c.resourceType == "" {
		return flagNotSetError{flagText: resourceTypeFlag}
	}

	funcMap := map[string]func(*client.Client) error{
		instanceResource: c.showInstance,
		accountResource:  c.showAccount,
		statusResource:   c.showStatus,
		timelineResource: c.showTimeline,
		listResource:     c.showList,
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
	var accountURI string

	if c.myAccount {
		authConfig, err := config.NewAuthenticationConfigFromFile()
		if err != nil {
			return fmt.Errorf("unable to retrieve the authentication configuration; %w", err)
		}

		accountURI = authConfig.CurrentAccount
	} else {
		if c.account == "" {
			return flagNotSetError{flagText: accountFlag}
		}

		accountURI = c.account
	}

	account, err := gts.GetAccount(accountURI)
	if err != nil {
		return fmt.Errorf("unable to retrieve the account details; %w", err)
	}

	fmt.Println(account)

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
		timeline, err = gts.GetHomeTimeline(c.timelineLimit)
	case "public":
		timeline, err = gts.GetPublicTimeline(c.timelineLimit)
	case "list":
		if c.listID == "" {
			return flagNotSetError{flagText: listIDFlag}
		}

		timeline, err = gts.GetListTimeline(c.listID, c.timelineLimit)
	case "tag":
		if c.tag == "" {
			return flagNotSetError{flagText: tagFlag}
		}

		timeline, err = gts.GetTagTimeline(c.tag, c.timelineLimit)
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
