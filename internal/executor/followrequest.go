package executor

import (
	"fmt"
	"net/rpc"

	"codeflow.dananglin.me.uk/apollo/enbas/internal/cli"
	"codeflow.dananglin.me.uk/apollo/enbas/internal/command"
	"codeflow.dananglin.me.uk/apollo/enbas/internal/config"
	"codeflow.dananglin.me.uk/apollo/enbas/internal/printer"
	"codeflow.dananglin.me.uk/apollo/enbas/internal/server"
)

// followRequestFunc is the function for the follow-request target for interacting
// with the user's follow requests.
func followRequestFunc(
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
	case cli.ActionAccept:
		return followRequestAccept(client, printSettings, cmd.FocusedTargetFlags)
	case cli.ActionReject:
		return followRequestReject(client, printSettings, cmd.FocusedTargetFlags)
	default:
		return unsupportedActionError{action: cmd.Action, target: cli.TargetFollowRequest}
	}
}

func followRequestAccept(
	client *rpc.Client,
	printSettings printer.Settings,
	flags []string,
) error {
	var accountName string

	// Parse the remaining flags.
	if err := cli.ParseFollowRequestAcceptFlags(
		&accountName,
		flags,
	); err != nil {
		return err
	}

	if accountName == "" {
		return missingAccountNameError{action: cli.ActionAccept}
	}

	var accountID string
	if err := client.Call(
		"GTSClient.GetAccountID",
		accountName,
		&accountID,
	); err != nil {
		return fmt.Errorf("received an error while getting the account ID: %w", err)
	}

	if err := client.Call(
		"GTSClient.AcceptFollowRequest",
		accountID,
		nil,
	); err != nil {
		return fmt.Errorf("unable to accept the follow request: %w", err)
	}

	printer.PrintSuccess(printSettings, "Successfully accepted the follow request from "+accountName+".")

	return nil
}

func followRequestReject(
	client *rpc.Client,
	printSettings printer.Settings,
	flags []string,
) error {
	var (
		accountName string
		accountID   string
	)

	// Parse the remaining flags.
	if err := cli.ParseFollowRequestRejectFlags(
		&accountName,
		flags,
	); err != nil {
		return err
	}

	if err := client.Call("GTSClient.GetAccountID", accountName, &accountID); err != nil {
		return fmt.Errorf("received an error while getting the account ID: %w", err)
	}

	if err := client.Call("GTSClient.RejectFollowRequest", accountID, nil); err != nil {
		return fmt.Errorf("unable to reject the follow request: %w", err)
	}

	printer.PrintSuccess(printSettings, "Successfully rejected the follow request from "+accountName+".")

	return nil
}
