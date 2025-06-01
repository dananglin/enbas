package executor

import (
	"os"

	"codeflow.dananglin.me.uk/apollo/enbas/internal/cli"
	"codeflow.dananglin.me.uk/apollo/enbas/internal/command"
	"codeflow.dananglin.me.uk/apollo/enbas/internal/config"
	internalFlag "codeflow.dananglin.me.uk/apollo/enbas/internal/flag"
	"codeflow.dananglin.me.uk/apollo/enbas/internal/printer"
)

// targetFunc is a type of function that performs the relevant operation
// on a specific target.
type targetFunc func(cfg config.Config, printSettings printer.Settings, cmd command.Command) error

func Execute() error {
	var (
		noColorFlag internalFlag.BoolValue
		noColor     bool
		configPath  string
	)

	// Initialise the print settings.
	printSettings := printer.NewSettings(
		true,
		"",
		0,
	)

	// Parse the top level flags.
	flagset := cli.NewTopLevelFlagset(&configPath, &noColorFlag)
	if err := flagset.Parse(os.Args[1:]); err != nil {
		printer.PrintFailure(
			printSettings,
			"error parsing the top-level flags: "+err.Error()+".",
		)

		return err //nolint:wrapcheck
	}

	// Get the user's 'no color' settings.
	if noColorFlag.IsSet() {
		noColor = noColorFlag.Value()
	} else if os.Getenv("NO_COLOR") != "" {
		noColor = true
	}

	// Load the configuration if the configuration file
	// is present.
	var cfg config.Config

	calculatedConfigPath, err := config.Path(configPath)
	if err != nil {
		printer.PrintFailure(
			printSettings,
			"error calculating the path to the configuration file: "+err.Error()+".",
		)

		return err //nolint:wrapcheck
	}

	cfgFileExists, err := config.FileExists(calculatedConfigPath)
	if err != nil {
		printer.PrintFailure(
			printSettings,
			"error checking if the configuration file is present: "+err.Error()+".",
		)

		return err //nolint:wrapcheck
	}

	if cfgFileExists {
		cfg, err = config.NewConfigFromFile(calculatedConfigPath)
		if err != nil {
			printer.PrintFailure(
				printSettings,
				"error loading the configuration: "+err.Error()+".",
			)

			return err //nolint:wrapcheck
		}
	} else {
		cfg.Path = calculatedConfigPath
	}

	if !cfg.IsZero() {
		// Update the print settings if the configuration was
		// successfully loaded from file.
		printSettings = printer.NewSettings(
			noColor,
			cfg.Integrations.Pager,
			cfg.LineWrapMaxWidth,
		)
	} else {
		// Otherwise update the print settings by only adjusting the
		// 'no color' setting.
		printSettings = printer.NewSettings(
			noColor,
			"",
			0,
		)
	}

	// Parse the command
	var cmd command.Command

	if flagset.NArg() == 0 {
		cmd = command.UsageCommand()
	} else {
		var err error

		// Parse the alias if it's used.
		args, err := command.ExtractArgsFromAlias(flagset.Args(), cfg.Aliases)
		if err != nil {
			printer.PrintFailure(
				printSettings,
				"error parsing the alias from the command: "+err.Error()+".",
			)

			return err //nolint:wrapcheck
		}

		// Parse the final command.
		cmd, err = command.Parse(args)
		if err != nil {
			printer.PrintFailure(
				printSettings,
				"error parsing the action and its arguments: "+err.Error()+".",
			)

			return err //nolint:wrapcheck
		}
	}

	if err := cmd.Validate(); err != nil {
		printer.PrintFailure(
			printSettings,
			"invalid command: "+err.Error()+".",
		)

		return err //nolint:wrapcheck
	}

	funcMap := targetFuncMap()

	targetFunc, ok := funcMap[cmd.FocusedTarget]
	if !ok {
		err := unrecognisedTargetError{target: cmd.FocusedTarget}

		printer.PrintFailure(
			printSettings,
			err.Error()+".",
		)

		return err
	}

	if err := targetFunc(cfg, printSettings, cmd); err != nil {
		printer.PrintFailure(
			printSettings,
			err.Error()+".",
		)

		return err
	}

	return nil
}
