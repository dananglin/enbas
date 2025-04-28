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

// blockedAccountsFunc is the function for the blocked-accounts target
// for interacting with blocked accounts.
func blockedAccountsFunc(
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
		return blockedAccountsShow(client, printSettings, cmd.FocusedTargetFlags)
	default:
		return unsupportedActionError{action: cmd.Action, target: cli.TargetBlockedAccounts}
	}
}

func blockedAccountsShow(
	client *rpc.Client,
	printSettings printer.Settings,
	flags []string,
) error {
	var limit int

	// Parse the remaining flags.
	if err := cli.ParseBlockedAccountsShowFlags(
		&limit,
		flags,
	); err != nil {
		return err
	}

	var blocked model.AccountList
	if err := client.Call(
		"GTSClient.GetBlockedAccounts",
		limit,
		&blocked,
	); err != nil {
		return fmt.Errorf("error retrieving the list of blocked accounts: %w", err)
	}

	if len(blocked.Accounts) > 0 {
		if err := printer.PrintAccountList(printSettings, blocked); err != nil {
			return fmt.Errorf("error printing the list of blocked accounts: %w", err)
		}
	} else {
		printer.PrintInfo("You have no blocked accounts.\n")
	}

	return nil
}
