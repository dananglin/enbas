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
		return mutedAccountsShow(client, printSettings, cmd.FocusedTargetFlags)
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
