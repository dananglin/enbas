package executor

import (
	"fmt"
	"net/rpc"

	"codeflow.dananglin.me.uk/apollo/enbas/internal/cli"
	"codeflow.dananglin.me.uk/apollo/enbas/internal/command"
	"codeflow.dananglin.me.uk/apollo/enbas/internal/config"
	"codeflow.dananglin.me.uk/apollo/enbas/internal/gtsclient"
	"codeflow.dananglin.me.uk/apollo/enbas/internal/printer"
	"codeflow.dananglin.me.uk/apollo/enbas/internal/server"
)

func noteFunc(
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
	case cli.ActionAdd:
		return noteAdd(
			client,
			printSettings,
			cmd.RelatedTarget,
			cmd.RelatedTargetFlags,
		)
	case cli.ActionRemove:
		return noteRemove(
			client,
			printSettings,
			cmd.RelatedTarget,
			cmd.RelatedTargetFlags,
		)
	default:
		return unsupportedActionError{action: cmd.Action, target: cli.TargetNote}
	}
}

func noteAdd(
	client *rpc.Client,
	printSettings printer.Settings,
	relatedTarget string,
	relatedTargetFlags []string,
) error {
	switch relatedTarget {
	case cli.TargetAccount:
		return noteAddToAccount(
			client,
			printSettings,
			relatedTargetFlags,
		)
	default:
		return unsupportedTargetToTargetError{
			action:        cli.ActionAdd,
			focusedTarget: cli.TargetNote,
			preposition:   cli.TargetActionPreposition(cli.TargetNote, cli.ActionAdd),
			relatedTarget: relatedTarget,
		}
	}
}

func noteAddToAccount(
	client *rpc.Client,
	printSettings printer.Settings,
	flags []string,
) error {
	var (
		accountName string
		content     string
	)

	// Parse the flags for the account target.
	if err := cli.ParseNoteAddToAccountFlags(
		&accountName,
		&content,
		flags,
	); err != nil {
		return err
	}

	if accountName == "" {
		return missingValueError{
			valueType: "name",
			target:    cli.TargetAccount,
			action:    "add the " + cli.TargetNote + " to",
		}
	}

	if content == "" {
		return missingValueError{
			valueType: "content",
			target:    cli.TargetNote,
			action:    "add to the " + cli.TargetAccount,
		}
	}

	var accountID string
	if err := client.Call("GTSClient.GetAccountID", accountName, &accountID); err != nil {
		return fmt.Errorf("error retrieving the account ID: %w", err)
	}

	if err := client.Call(
		"GTSClient.SetPrivateNote",
		gtsclient.SetPrivateNoteArgs{
			AccountID: accountID,
			Note:      content,
		},
		nil,
	); err != nil {
		return fmt.Errorf("unable to add the private note to the account: %w", err)
	}

	printer.PrintSuccess(printSettings, "Successfully added the private note to the account.")

	return nil
}

func noteRemove(
	client *rpc.Client,
	printSettings printer.Settings,
	relatedTarget string,
	relatedTargetFlags []string,
) error {
	switch relatedTarget {
	case cli.TargetAccount:
		return noteRemoveFromAccount(client, printSettings, relatedTargetFlags)
	default:
		return unsupportedTargetToTargetError{
			action:        cli.ActionRemove,
			focusedTarget: cli.TargetNote,
			preposition:   cli.TargetActionPreposition(cli.TargetNote, cli.ActionRemove),
			relatedTarget: relatedTarget,
		}
	}
}

func noteRemoveFromAccount(
	client *rpc.Client,
	printSettings printer.Settings,
	flags []string,
) error {
	var accountName string

	// Parse the flags for the account target.
	if err := cli.ParseNoteRemoveFromAccountFlags(
		&accountName,
		flags,
	); err != nil {
		return err
	}

	if accountName == "" {
		return missingValueError{
			valueType: "name",
			target:    cli.TargetAccount,
			action:    "remove the " + cli.TargetNote + " from",
		}
	}

	var accountID string
	if err := client.Call("GTSClient.GetAccountID", accountName, &accountID); err != nil {
		return fmt.Errorf("error retrieving the account ID: %w", err)
	}

	if err := client.Call(
		"GTSClient.SetPrivateNote",
		gtsclient.SetPrivateNoteArgs{
			AccountID: accountID,
			Note:      "",
		},
		nil,
	); err != nil {
		return fmt.Errorf("error removing the private private note from the account: %w", err)
	}

	printer.PrintSuccess(printSettings, "Successfully removed the private note from the account.")

	return nil
}
