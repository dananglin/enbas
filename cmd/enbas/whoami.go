package main

import (
	"flag"
	"fmt"

	"codeflow.dananglin.me.uk/apollo/enbas/internal/config"
)

type whoAmICommand struct {
	*flag.FlagSet
}

func newWhoAmICommand(name, summary string) *whoAmICommand {
	command := whoAmICommand{
		FlagSet: flag.NewFlagSet(name, flag.ExitOnError),
	}

	command.Usage = commandUsageFunc(name, summary, command.FlagSet)

	return &command
}

func (c *whoAmICommand) Execute() error {
	config, err := config.NewAuthenticationConfigFromFile()
	if err != nil {
		return fmt.Errorf("unable to load the credential config; %w", err)
	}

	fmt.Printf("You are %s\n", config.CurrentAccount)

	return nil
}
