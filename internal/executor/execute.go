package executor

import (
	"os"

	"codeflow.dananglin.me.uk/apollo/enbas/internal/cli"
	"codeflow.dananglin.me.uk/apollo/enbas/internal/command"
	internalFlag "codeflow.dananglin.me.uk/apollo/enbas/internal/flag"
	"codeflow.dananglin.me.uk/apollo/enbas/internal/printer"
)

type topLevelOpts struct {
	configPath string
	noColor    bool
}

// targetFunc is a type of function that performs the relevant operation
// on a specific target.
type targetFunc func(opts topLevelOpts, cmd command.Command) error

func Execute() error {
	var (
		opts        topLevelOpts
		noColorFlag internalFlag.BoolPtrValue
	)

	errorPrintSettings := printer.NewSettings(opts.noColor, "", 0)

	flagset := cli.NewTopLevelFlagset(&opts.configPath, &noColorFlag)
	if err := flagset.Parse(os.Args[1:]); err != nil {
		printer.PrintFailure(
			errorPrintSettings,
			"error parsing the top-level flags: "+err.Error()+".",
		)

		return err
	}

	if noColorFlag.Value != nil {
		opts.noColor = *noColorFlag.Value
	} else if os.Getenv("NO_COLOR") != "" {
		opts.noColor = true
	}

	var cmd command.Command

	if flagset.NArg() == 0 {
		cmd = command.HelpCommand()
	} else {
		var err error
		cmd, err = command.Parse(flagset.Args())
		if err != nil {
			printer.PrintFailure(
				errorPrintSettings,
				"error parsing the action and its arguments: "+err.Error()+".",
			)

			return err
		}
	}

	if err := cmd.Validate(); err != nil {
		printer.PrintFailure(
			errorPrintSettings,
			"invalid command: "+err.Error()+".",
		)
		return err
	}

	funcMap := targetFuncMap()

	targetFunc, ok := funcMap[cmd.FocusedTarget]
	if !ok {
		err := unrecognisedTargetError{target: cmd.FocusedTarget}

		printer.PrintFailure(
			errorPrintSettings,
			err.Error()+".",
		)

		return err
	}

	if err := targetFunc(opts, cmd); err != nil {
		printer.PrintFailure(
			errorPrintSettings,
			err.Error()+".",
		)

		return err
	}

	return nil
}
