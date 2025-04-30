package executor

import (
	"fmt"

	"codeflow.dananglin.me.uk/apollo/enbas/internal/cli"
	"codeflow.dananglin.me.uk/apollo/enbas/internal/command"
	"codeflow.dananglin.me.uk/apollo/enbas/internal/config"
	"codeflow.dananglin.me.uk/apollo/enbas/internal/printer"
)

// aliasesFunc is the function for the 'aliases' target for interacting
// with the user's list of alises.
func aliasesFunc(
	cfg config.Config,
	printSettings printer.Settings,
	cmd command.Command,
) error {
	if cfg.IsZero() {
		return zeroConfigurationError{path: cfg.Path}
	}

	switch cmd.Action {
	case cli.ActionShow:
		return aliasesShow(
			cfg.Aliases,
			printSettings,
		)
	default:
		return unsupportedActionError{action: cmd.Action, target: cli.TargetAliases}
	}
}

func aliasesShow(
	aliases map[string]string,
	printSettings printer.Settings,
) error {
	if len(aliases) > 0 {
		if err := printer.PrintAliases(printSettings, aliases); err != nil {
			return fmt.Errorf("error printing the list of aliases: %w", err)
		}
	} else {
		printer.PrintInfo("You have no aliases.\n")
	}

	return nil
}
