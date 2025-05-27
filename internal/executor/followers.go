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

func followersFunc(
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
		return followersShow(
			session.Client(),
			printSettings,
			cmd.RelatedTarget,
			cmd.RelatedTargetFlags,
		)
	default:
		return unsupportedActionError{action: cmd.Action, target: cli.TargetFollowers}
	}
}

func followersShow(
	client *rpc.Client,
	printSettings printer.Settings,
	relatedTarget string,
	relatedTargetFlags []string,
) error {
	switch relatedTarget {
	case cli.TargetAccount:
		return followersShowFromAccount(
			client,
			printSettings,
			relatedTargetFlags,
		)
	default:
		return unsupportedTargetToTargetError{
			action:        cli.ActionShow,
			focusedTarget: cli.TargetFollowers,
			preposition:   cli.TargetActionPreposition(cli.TargetFollowers, cli.ActionShow),
			relatedTarget: relatedTarget,
		}
	}
}

func followersShowFromAccount(
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
	if err := cli.ParseFollowersShowFromAccountFlags(
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
			return fmt.Errorf("error retrieving your account ID: %w", err)
		}
	} else {
		if accountName == "" {
			return missingValueError{
				valueType: "name",
				target:    cli.TargetAccount,
				action:    "show the list of followers from",
			}
		}

		if err := client.Call(
			"GTSClient.GetAccountID",
			accountName,
			&accountID,
		); err != nil {
			return fmt.Errorf("error retrieving the account ID: %w", err)
		}
	}

	var followers model.AccountList
	if err := client.Call(
		"GTSClient.GetFollowers",
		gtsclient.GetFollowersArgs{
			AccountID: accountID,
			Limit:     limit,
		},
		&followers,
	); err != nil {
		return fmt.Errorf("error retrieving the list of followers: %w", err)
	}

	if len(followers.Accounts) > 0 {
		if err := printer.PrintAccountList(printSettings, followers); err != nil {
			return fmt.Errorf("error printing the list of followers: %w", err)
		}
	} else {
		printer.PrintInfo("There are no followers for this account (or the list is hidden).\n")
	}

	return nil
}
