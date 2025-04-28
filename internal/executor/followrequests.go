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

// followRequestsFunc is the function for the follow-requests target for
// interacting with the user's list of follow requests.
func followRequestsFunc(
	cfg config.Config,
	printSettings printer.Settings,
	cmd command.Command,
) error {
	if cfg.IsZero() {
		return zeroConfigurationError{path: cfg.Path}
	}

	// Create the client to interact with the GoToSocial instance.
	client, err := server.Connect(cfg.Server, cfg.Path)
	if err != nil {
		return fmt.Errorf("error creating the client for the daemon process: %w", err)
	}
	defer client.Close()

	switch cmd.Action {
	case cli.ActionShow:
		return followRequestsShow(
			client,
			printSettings,
			cmd.FocusedTargetFlags,
		)
	default:
		return unsupportedActionError{action: cmd.Action, target: cli.TargetFollowRequests}
	}
}

func followRequestsShow(
	client *rpc.Client,
	printSettings printer.Settings,
	flags []string,
) error {
	var limit int

	if err := cli.ParseFollowRequestsShowFlags(
		&limit,
		flags,
	); err != nil {
		return err
	}

	var requests model.AccountList
	if err := client.Call("GTSClient.GetFollowRequests", limit, &requests); err != nil {
		return fmt.Errorf("unable to retrieve the list of follow requests: %w", err)
	}

	if len(requests.Accounts) > 0 {
		if err := printer.PrintAccountList(printSettings, requests); err != nil {
			return fmt.Errorf("error printing the list of follow requests: %w", err)
		}
	} else {
		printer.PrintInfo("You have no follow requests.\n")
	}

	return nil
}
