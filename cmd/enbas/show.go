package main

import (
	"errors"
	"flag"
	"fmt"

	"codeflow.dananglin.me.uk/apollo/enbas/internal/client"
	"codeflow.dananglin.me.uk/apollo/enbas/internal/config"
)

type showCommand struct {
	*flag.FlagSet
	myAccount  bool
	targetType string
	account    string
	statusID   string
}

func newShowCommand(name, summary string) *showCommand {
	command := showCommand{
		FlagSet:    flag.NewFlagSet(name, flag.ExitOnError),
		myAccount:  false,
		targetType: "",
		account:    "",
		statusID:   "",
	}

	command.BoolVar(&command.myAccount, "my-account", false, "set to true to lookup your account")
	command.StringVar(&command.targetType, "type", "", "specify the type of resource to display")
	command.StringVar(&command.account, "account", "", "specify the account URI to lookup")
	command.StringVar(&command.statusID, "status-id", "", "specify the ID of the status to display")
	command.Usage = commandUsageFunc(name, summary, command.FlagSet)

	return &command
}

func (c *showCommand) Execute() error {
	gtsClient, err := client.NewClientFromConfig()
	if err != nil {
		return fmt.Errorf("unable to create the GoToSocial client; %w", err)
	}

	funcMap := map[string]func(*client.Client) error{
		"instance": c.showInstance,
		"account":  c.showAccount,
		"status":   c.showStatus,
	}

	doFunc, ok := funcMap[c.targetType]
	if !ok {
		return fmt.Errorf("unsupported type %q", c.targetType)
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
			return errors.New("the account flag is not set")
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
		return errors.New("the status-id flag is not set")
	}

	status, err := gts.GetStatus(c.statusID)
	if err != nil {
		return fmt.Errorf("unable to retrieve the status; %w", err)
	}

	fmt.Println(status)

	return nil
}
