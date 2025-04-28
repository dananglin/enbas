package executor

import (
	"fmt"

	"codeflow.dananglin.me.uk/apollo/enbas/internal/cli"
	"codeflow.dananglin.me.uk/apollo/enbas/internal/command"
	"codeflow.dananglin.me.uk/apollo/enbas/internal/config"
	"codeflow.dananglin.me.uk/apollo/enbas/internal/gtsclient"
	"codeflow.dananglin.me.uk/apollo/enbas/internal/printer"
	"codeflow.dananglin.me.uk/apollo/enbas/internal/server"
)

func serverFunc(
	opts topLevelOpts,
	cmd command.Command,
) error {
	// Load the configuration from file.
	cfg, err := config.NewConfigFromFile(opts.configPath)
	if err != nil {
		return fmt.Errorf("unable to load configuration: %w", err)
	}

	// Create the printer for the executor.
	printSettings := printer.NewSettings(
		opts.noColor,
		"",
		cfg.LineWrapMaxWidth,
	)

	switch cmd.Action {
	case cli.ActionStart:
		return serverStart(cfg, printSettings, cmd.FocusedTargetFlags)
	default:
		return unsupportedActionError{action: cmd.Action, target: cli.TargetServer}
	}
}

func serverStart(
	cfg config.Config,
	printSettings printer.Settings,
	flags []string,
) error {
	var withoutIdleTimeout bool

	// Parse the remaining flags.
	if err := cli.ParseServerStartFlags(
		&withoutIdleTimeout,
		flags,
	); err != nil {
		return err
	}

	gtsClient, err := gtsclient.NewGTSClient(cfg)
	if err != nil {
		return fmt.Errorf("unable to create the GoToSocial client: %w", err)
	}

	if err := server.Run(
		printSettings,
		gtsClient,
		cfg.Server.SocketPath,
		withoutIdleTimeout,
		cfg.Server.IdleTimeout,
	); err != nil {
		return fmt.Errorf("error running Enbas in server mode: %w", err)
	}

	return nil
}
