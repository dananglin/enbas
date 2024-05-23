package executor

import (
	"flag"
	"fmt"

	"codeflow.dananglin.me.uk/apollo/enbas/internal/config"
)

type WhoAmIExecutor struct {
	*flag.FlagSet

	topLevelFlags TopLevelFlags
}

func NewWhoAmIExecutor(tlf TopLevelFlags, name, summary string) *WhoAmIExecutor {
	whoExe := WhoAmIExecutor{
		FlagSet:       flag.NewFlagSet(name, flag.ExitOnError),
		topLevelFlags: tlf,
	}

	whoExe.Usage = commandUsageFunc(name, summary, whoExe.FlagSet)

	return &whoExe
}

func (c *WhoAmIExecutor) Execute() error {
	config, err := config.NewCredentialsConfigFromFile(c.topLevelFlags.ConfigDir)
	if err != nil {
		return fmt.Errorf("unable to load the credential config; %w", err)
	}

	fmt.Printf("You are logged in as %q.\n", config.CurrentAccount)

	return nil
}
