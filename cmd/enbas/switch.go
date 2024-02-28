package main

import (
	"flag"
	"fmt"

	"codeflow.dananglin.me.uk/apollo/enbas/internal/config"
)

type switchCommand struct {
	*flag.FlagSet
	toAccount string
}

func newSwitchCommand(name, summary string) *switchCommand {
	command := switchCommand{
		FlagSet:   flag.NewFlagSet(name, flag.ExitOnError),
		toAccount: "",
	}

	command.StringVar(&command.toAccount, toAccountFlag, "", "the account to switch to")

	command.Usage = commandUsageFunc(name, summary, command.FlagSet)

	return &command
}

func (c *switchCommand) Execute() error {
	if c.toAccount == "" {
		return flagNotSetError{flagText: toAccountFlag}
	}

	if err := config.UpdateCurrentAccount(c.toAccount); err != nil {
		return fmt.Errorf("unable to switch accounts; %w", err)
	}

	fmt.Printf("The current account is now set to %q.\n", c.toAccount)

	return nil
}
