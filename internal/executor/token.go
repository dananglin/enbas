package executor

import (
	"fmt"
	"net/rpc"

	"codeflow.dananglin.me.uk/apollo/enbas/internal/cli"
	"codeflow.dananglin.me.uk/apollo/enbas/internal/command"
	"codeflow.dananglin.me.uk/apollo/enbas/internal/config"
	"codeflow.dananglin.me.uk/apollo/enbas/internal/model"
	"codeflow.dananglin.me.uk/apollo/enbas/internal/printer"
	"codeflow.dananglin.me.uk/apollo/enbas/internal/server"
)

func tokenFunc(
	opts topLevelOpts,
	cmd command.Command,
) error {
	// Load the configuration from file.
	cfg, err := config.NewConfigFromFile(opts.configPath)
	if err != nil {
		return fmt.Errorf("unable to load configuration: %w", err)
	}

	// Create the print settings.
	printSettings := printer.NewSettings(
		opts.noColor,
		cfg.Integrations.Pager,
		cfg.LineWrapMaxWidth,
	)

	// Create the client to interact with the GoToSocial instance.
	client, err := server.Connect(cfg.Server, opts.configPath)
	if err != nil {
		return fmt.Errorf("error creating the client for the daemon process: %w", err)
	}
	defer client.Close()

	switch cmd.Action {
	case cli.ActionShow:
		return tokenShow(
			client,
			printSettings,
			cmd.FocusedTargetFlags,
		)
	case cli.ActionInvalidate:
		return tokenInvalidate(
			client,
			printSettings,
			cmd.FocusedTargetFlags,
		)
	default:
		return unsupportedActionError{action: cmd.Action, target: cli.TargetToken}
	}
}

func tokenShow(
	client *rpc.Client,
	printSettings printer.Settings,
	flags []string,
) error {
	var tokenID string

	// Parse the remaining flags.
	if err := cli.ParseTokenShowFlags(
		&tokenID,
		flags,
	); err != nil {
		return err
	}

	var token model.Token
	if err := client.Call(
		"GTSClient.GetToken",
		tokenID,
		&token,
	); err != nil {
		return fmt.Errorf("error retrieving the details of the token: %w", err)
	}

	if err := printer.PrintToken(
		printSettings,
		token,
	); err != nil {
		return fmt.Errorf("error printing the token: %w", err)
	}

	return nil
}

func tokenInvalidate(
	client *rpc.Client,
	printSettings printer.Settings,
	flags []string,
) error {
	var tokenID string

	// Parse the remaining flags.
	if err := cli.ParseTokenInvalidateFlags(
		&tokenID,
		flags,
	); err != nil {
		return err
	}

	if err := client.Call(
		"GTSClient.InvalidateToken",
		tokenID,
		nil,
	); err != nil {
		return fmt.Errorf("error invalidating the token: %w", err)
	}

	printer.PrintSuccess(
		printSettings,
		"The token was successfully invalidated.",
	)

	return nil
}
