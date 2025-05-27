package executor

import (
	"fmt"
	"net/rpc"

	"codeflow.dananglin.me.uk/apollo/enbas/internal/cli"
	"codeflow.dananglin.me.uk/apollo/enbas/internal/command"
	"codeflow.dananglin.me.uk/apollo/enbas/internal/config"
	"codeflow.dananglin.me.uk/apollo/enbas/internal/gtsclient"
	"codeflow.dananglin.me.uk/apollo/enbas/internal/model"
	"codeflow.dananglin.me.uk/apollo/enbas/internal/printer"
	"codeflow.dananglin.me.uk/apollo/enbas/internal/server"
)

func followingsFunc(
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
		return followingsShow(
			session.Client(),
			printSettings,
			cmd.RelatedTarget,
			cmd.RelatedTargetFlags,
		)
	default:
		return unsupportedActionError{action: cmd.Action, target: cli.TargetFollowings}
	}
}

func followingsShow(
	client *rpc.Client,
	printSettings printer.Settings,
	relatedTarget string,
	relatedTargetFlags []string,
) error {
	switch relatedTarget {
	case cli.TargetAccount:
		return followingsShowFromAccount(
			client,
			printSettings,
			relatedTargetFlags,
		)
	default:
		return unsupportedTargetToTargetError{
			action:        cli.ActionShow,
			focusedTarget: cli.TargetFollowings,
			preposition:   cli.TargetActionPreposition(cli.TargetFollowings, cli.ActionShow),
			relatedTarget: relatedTarget,
		}
	}
}

func followingsShowFromAccount(
	client *rpc.Client,
	printSettings printer.Settings,
	flags []string,
) error {
	var (
		accountName string
		limit       int
		myAccount   bool
	)

	// Parse the remaining flags
	if err := cli.ParseFollowingsShowFromAccountFlags(
		&accountName,
		&limit,
		&myAccount,
		flags,
	); err != nil {
		return err
	}

	var accountID string

	if myAccount {
		if err := client.Call(
			"GTSClient.GetMyAccountID",
			gtsclient.NoRPCArgs{},
			&accountID,
		); err != nil {
			return fmt.Errorf("error getting your account ID: %w", err)
		}
	} else {
		if err := client.Call(
			"GTSClient.GetAccountID",
			accountName,
			&accountID,
		); err != nil {
			return fmt.Errorf("error retrieving the account ID: %w", err)
		}
	}

	var followings model.AccountList
	if err := client.Call(
		"GTSClient.GetFollowing",
		gtsclient.GetFollowingsArgs{
			AccountID: accountID,
			Limit:     limit,
		},
		&followings,
	); err != nil {
		return fmt.Errorf("error retrieving the list of followings: %w", err)
	}

	if len(followings.Accounts) > 0 {
		if err := printer.PrintAccountList(printSettings, followings); err != nil {
			return fmt.Errorf("error printing the list of followings: %w", err)
		}
	} else {
		printer.PrintInfo("This account is not following anyone or the list is hidden.\n")
	}

	return nil
}
