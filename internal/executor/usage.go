package executor

import (
	"fmt"
	"strings"

	"codeflow.dananglin.me.uk/apollo/enbas/internal/cli"
	"codeflow.dananglin.me.uk/apollo/enbas/internal/command"
	"codeflow.dananglin.me.uk/apollo/enbas/internal/config"
	"codeflow.dananglin.me.uk/apollo/enbas/internal/printer"
)

// usageFunc is the function for the help target for printing the
// help documentation to the screen for the user.
func usageFunc(
	_ config.Config,
	printSettings printer.Settings,
	cmd command.Command,
) error {
	switch cmd.Action {
	case cli.ActionShow:
		return usageShow(printSettings, cmd.FocusedTargetFlags)
	default:
		return unsupportedActionError{action: cmd.Action, target: cli.TargetUsage}
	}
}

func usageShow(
	printSettings printer.Settings,
	flags []string,
) error {
	var (
		target    string
		operation string
	)

	// Parse the remaining flags.
	if err := cli.ParseUsageShowFlags(
		&target,
		&operation,
		flags,
	); err != nil {
		return err //nolint:wrapcheck
	}

	switch {
	case target != "":
		return usageShowTarget(printSettings, target)
	case operation != "":
		return usageShowOperation(printSettings, operation)
	default:
		return usageShowRoot(printSettings)
	}
}

func usageShowRoot(printSettings printer.Settings) error {
	if err := printer.PrintUsageRoot(
		printSettings,
		cli.TargetDescMap(),
		cli.TopLevelFlagsUsageMap(),
		cli.TargetStatus,
	); err != nil {
		return fmt.Errorf("error printing the usage documentation: %w", err)
	}

	return nil
}

func usageShowTarget(
	printSettings printer.Settings,
	target string,
) error {
	desc, ok := cli.TargetDesc(target)
	if !ok {
		return fmt.Errorf("usage error: %w", unrecognisedTargetError{target: target})
	}

	operations, ok := cli.GetUsageOperations(target)
	if !ok {
		return usageNoOpFoundError{target: target}
	}

	if err := printer.PrintUsageTarget(
		printSettings,
		target,
		desc,
		cli.TopLevelFlagsUsageMap(),
		operations,
	); err != nil {
		return fmt.Errorf("error printing the usage documentation: %w", err)
	}

	return nil
}

func usageShowOperation(
	printSettings printer.Settings,
	operation string,
) error {
	// Parse the operation parameter to get the name of the focused target.
	parsedOperation, err := command.Parse(strings.Split(operation, " "))
	if err != nil {
		return fmt.Errorf("unable to parse the operation: %w", err)
	}

	usageOperation, ok := cli.GetUsageOperation(
		parsedOperation.FocusedTarget,
		operation,
	)
	if !ok {
		return usageNoOpForTargetError{
			operation: operation,
			target:    parsedOperation.FocusedTarget,
		}
	}

	flagUsageMap := cli.FlagUsageMap()

	flagDescriptions := make([]string, len(usageOperation.Flags))

	for idx := range usageOperation.Flags {
		flagDescriptions[idx] = strings.ReplaceAll(
			strings.ReplaceAll(
				flagUsageMap[usageOperation.Flags[idx]],
				"{target}",
				parsedOperation.FocusedTarget,
			),
			"{action}",
			parsedOperation.Action,
		)
	}

	if err := printer.PrintUsageOperation(
		printSettings,
		cli.TopLevelFlagsUsageMap(),
		operation,
		usageOperation,
		flagDescriptions,
	); err != nil {
		return fmt.Errorf("error printing the usage documentation: %w", err)
	}

	return nil
}
