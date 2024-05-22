package main

import (
	"flag"
	"fmt"

	"codeflow.dananglin.me.uk/apollo/enbas/internal/config"
)

type whoAmICommand struct {
	*flag.FlagSet

	topLevelFlags topLevelFlags
}

func newWhoAmICommand(tlf topLevelFlags, name, summary string) *whoAmICommand {
	command := whoAmICommand{
		FlagSet: flag.NewFlagSet(name, flag.ExitOnError),
		topLevelFlags: tlf,
	}

	command.Usage = commandUsageFunc(name, summary, command.FlagSet)

	return &command
}

func (c *whoAmICommand) Execute() error {
	config, err := config.NewCredentialsConfigFromFile(c.topLevelFlags.configDir)
	if err != nil {
		return fmt.Errorf("unable to load the credential config; %w", err)
	}

	fmt.Printf("You are %s\n", config.CurrentAccount)

	return nil
}
