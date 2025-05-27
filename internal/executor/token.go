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
	cfg config.Config,
	printSettings printer.Settings,
	cmd command.Command,
) error {
	if cfg.IsZero() {
		return zeroConfigurationError{path: cfg.Path}
	}

	// Create the session to interact with the GoToSocial instance.
	session, err := server.StartSession(cfg.Server, cfg.Path)
	if err != nil {
		return fmt.Errorf("error creating the client for the daemon process: %w", err)
	}
	defer server.EndSession(session)

	switch cmd.Action {
	case cli.ActionShow:
		return tokenShow(
			session.Client(),
			printSettings,
			cmd.FocusedTargetFlags,
		)
	case cli.ActionInvalidate:
		return tokenInvalidate(
			session.Client(),
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
