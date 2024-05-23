package executor

import (
	"flag"
	"fmt"

	"codeflow.dananglin.me.uk/apollo/enbas/internal/config"
)

type SwitchExecutor struct {
	*flag.FlagSet

	topLevelFlags  TopLevelFlags
	toResourceType string
	accountName    string
}

func NewSwitchExecutor(tlf TopLevelFlags, name, summary string) *SwitchExecutor {
	switchExe := SwitchExecutor{
		FlagSet:       flag.NewFlagSet(name, flag.ExitOnError),
		topLevelFlags: tlf,
	}

	switchExe.StringVar(&switchExe.toResourceType, flagTo, "", "the account to switch to")
	switchExe.StringVar(&switchExe.accountName, flagAccountName, "", "the name of the account to switch to")

	switchExe.Usage = commandUsageFunc(name, summary, switchExe.FlagSet)

	return &switchExe
}

func (s *SwitchExecutor) Execute() error {
	funcMap := map[string]func() error{
		resourceAccount: s.switchToAccount,
	}

	doFunc, ok := funcMap[s.toResourceType]
	if !ok {
		return UnsupportedTypeError{resourceType: s.toResourceType}
	}

	return doFunc()
}

func (s *SwitchExecutor) switchToAccount() error {
	if s.accountName == "" {
		return NoAccountSpecifiedError{}
	}

	if err := config.UpdateCurrentAccount(s.accountName, s.topLevelFlags.ConfigDir); err != nil {
		return fmt.Errorf("unable to switch account to the account; %w", err)
	}

	fmt.Printf("The current account is now set to %q.\n", s.accountName)

	return nil
}
