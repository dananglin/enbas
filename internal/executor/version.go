package executor

import (
	"fmt"

	"codeflow.dananglin.me.uk/apollo/enbas/internal/cli"
	"codeflow.dananglin.me.uk/apollo/enbas/internal/command"
	"codeflow.dananglin.me.uk/apollo/enbas/internal/printer"
)

// versionFunc is the function for the version target which
// prints the application's build information.
func versionFunc(
	opts topLevelOpts,
	cmd command.Command,
) error {
	// Create the print settings
	printSettings := printer.NewSettings(opts.noColor, "", 0)

	switch cmd.Action {
	case cli.ActionShow:
		return versionShow(printSettings, cmd.FocusedTargetFlags)
	default:
		return unsupportedActionError{action: cmd.Action, target: cli.TargetVersion}
	}
}

func versionShow(
	printSettings printer.Settings,
	flags []string,
) error {
	var full bool

	// Parse the remaining flags.
	if err := cli.ParseVersionShowFlags(
		&full,
		flags,
	); err != nil {
		return err
	}

	if err := printer.PrintVersion(printSettings, full); err != nil {
		return fmt.Errorf("execution error: %w", err)
	}

	return nil
}
