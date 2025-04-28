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

func threadFunc(
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
		return threadShow(
			client,
			printSettings,
			cmd.RelatedTarget,
			cmd.RelatedTargetFlags,
		)
	default:
		return unsupportedActionError{action: cmd.Action, target: cli.TargetThread}
	}
}

func threadShow(
	client *rpc.Client,
	printSettings printer.Settings,
	relatedTarget string,
	relatedTargetFlags []string,
) error {
	switch relatedTarget {
	case cli.TargetStatus:
		return threadShowFromStatus(
			client,
			printSettings,
			relatedTargetFlags,
		)
	default:
		return unsupportedTargetToTargetError{
			action:        cli.ActionShow,
			focusedTarget: cli.TargetThread,
			preposition:   cli.TargetActionPreposition(cli.TargetThread, cli.ActionShow),
			relatedTarget: relatedTarget,
		}
	}
}

func threadShowFromStatus(
	client *rpc.Client,
	printSettings printer.Settings,
	flags []string,
) error {
	var statusID string

	// Parse the flags for the status target.
	if err := cli.ParseThreadShowFromStatusFlags(
		&statusID,
		flags,
	); err != nil {
		return err
	}

	if statusID == "" {
		return missingIDError{
			target: cli.TargetStatus,
			action: "view the thread from",
		}
	}

	var myAccountID string
	if err := client.Call("GTSClient.GetMyAccountID", gtsclient.NoRPCArgs{}, &myAccountID); err != nil {
		return fmt.Errorf("unable to get your account ID: %w", err)
	}

	var thread model.Thread
	if err := client.Call("GTSClient.GetThread", statusID, &thread); err != nil {
		return fmt.Errorf("error retrieving the thread: %w", err)
	}

	if err := client.Call("GTSClient.GetStatus", statusID, &thread.Context); err != nil {
		return fmt.Errorf("error retrieving the status in context: %w", err)
	}

	// Print the thread
	if err := printer.PrintThread(printSettings, thread, myAccountID); err != nil {
		return fmt.Errorf("error printing the thread: %w", err)
	}

	return nil
}
