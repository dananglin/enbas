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

func mutedAccountsFunc(
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
		return mutedAccountsShow(session.Client(), printSettings, cmd.FocusedTargetFlags)
	default:
		return unsupportedActionError{action: cmd.Action, target: cli.TargetMutedAccounts}
	}
}

func mutedAccountsShow(
	client *rpc.Client,
	printSettings printer.Settings,
	flags []string,
) error {
	var limit int

	// Parse the remaining flags.
	if err := cli.ParseMutedAccountsShowFlags(
		&limit,
		flags,
	); err != nil {
		return err
	}

	var muted model.AccountList
	if err := client.Call(
		"GTSClient.GetMutedAccounts",
		limit,
		&muted,
	); err != nil {
		return fmt.Errorf("error retrieving the list of muted accounts: %w", err)
	}

	if len(muted.Accounts) > 0 {
		if err := printer.PrintAccountList(printSettings, muted); err != nil {
			return fmt.Errorf("error printing the list of muted accounts: %w", err)
		}
	} else {
		printer.PrintInfo("You have not muted any accounts.\n")
	}

	return nil
}
