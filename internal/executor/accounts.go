package executor

import (
	"fmt"
	"net/rpc"

	"codeflow.dananglin.me.uk/apollo/enbas/internal/cli"
	"codeflow.dananglin.me.uk/apollo/enbas/internal/command"
	"codeflow.dananglin.me.uk/apollo/enbas/internal/config"
	internalFlag "codeflow.dananglin.me.uk/apollo/enbas/internal/flag"
	"codeflow.dananglin.me.uk/apollo/enbas/internal/gtsclient"
	"codeflow.dananglin.me.uk/apollo/enbas/internal/model"
	"codeflow.dananglin.me.uk/apollo/enbas/internal/printer"
	"codeflow.dananglin.me.uk/apollo/enbas/internal/server"
)

// accountsFunc is the function for the account target for
// interacting with multiple accounts.
func accountsFunc(
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
		return accountsAdd(
			client,
			printSettings,
			cmd.RelatedTarget,
			cmd.RelatedTargetFlags,
		)
	case cli.ActionRemove:
		return accountsRemove(
			client,
			printSettings,
			cmd.RelatedTarget,
			cmd.RelatedTargetFlags,
		)
	default:
		return unsupportedActionError{action: cmd.Action, target: cli.TargetAccounts}
	}
}

func accountsAdd(
	client *rpc.Client,
	printSettings printer.Settings,
	relatedTarget string,
	relatedTargetFlags []string,
) error {
	switch relatedTarget {
	case cli.TargetList:
		return accountsAddToList(
			client,
			printSettings,
			relatedTargetFlags,
		)
	default:
		return unsupportedTargetToTargetError{
			action:        cli.ActionAdd,
			focusedTarget: cli.TargetAccounts,
			preposition:   cli.TargetActionPreposition(cli.TargetAccounts, cli.ActionAdd),
			relatedTarget: relatedTarget,
		}
	}
}

func accountsAddToList(
	client *rpc.Client,
	printSettings printer.Settings,
	flags []string,
) error {
	var (
		listID       string
		accountNames = internalFlag.NewStringSliceValue()
	)

	// Parse the remaining flags
	if err := cli.ParseAccountsAddToListFlags(
		&listID,
		&accountNames,
		flags,
	); err != nil {
		return err
	}

	if listID == "" {
		return missingIDError{
			target: cli.TargetList,
			action: "add the accounts to",
		}
	}

	if accountNames.Empty() {
		return zeroAccountNamesError{
			action: "add to the list",
		}
	}

	var accounts []model.Account
	if err := client.Call(
		"GTSClient.GetMultipleAccounts",
		accountNames,
		&accounts,
	); err != nil {
		return fmt.Errorf("error retrieving the accounts: %w", err)
	}

	accountIDs := make([]string, len(accounts))
	for idx := range accounts {
		var relationship model.AccountRelationship
		if err := client.Call(
			"GTSClient.GetAccountRelationship",
			accounts[idx].ID,
			&relationship,
		); err != nil {
			return fmt.Errorf("error retrieving your relationship to %s: %w", accounts[idx].Acct, err)
		}

		if !relationship.Following {
			return notFollowingError{account: accounts[idx].Acct}
		}

		accountIDs[idx] = accounts[idx].ID
	}

	if err := client.Call(
		"GTSClient.AddAccountsToList",
		gtsclient.AddAccountsToListArgs{
			ListID:     listID,
			AccountIDs: accountIDs,
		},
		nil,
	); err != nil {
		return fmt.Errorf("error adding the accounts to the list: %w", err)
	}

	printer.PrintSuccess(printSettings, "Successfully added the account(s) to the list.")

	return nil
}

func accountsRemove(
	client *rpc.Client,
	printSettings printer.Settings,
	relatedTarget string,
	relatedTargetFlags []string,
) error {
	switch relatedTarget {
	case cli.TargetList:
		return accountsRemoveFromList(
			client,
			printSettings,
			relatedTargetFlags,
		)
	default:
		return unsupportedTargetToTargetError{
			action:        cli.ActionRemove,
			focusedTarget: cli.TargetAccounts,
			preposition:   cli.TargetActionPreposition(cli.TargetAccounts, cli.ActionRemove),
			relatedTarget: relatedTarget,
		}
	}
}

func accountsRemoveFromList(
	client *rpc.Client,
	printSettings printer.Settings,
	flags []string,
) error {
	var (
		listID       string
		accountNames = internalFlag.NewStringSliceValue()
	)

	// Parse the remaining flags.
	if err := cli.ParseAccountsRemoveFromListFlags(
		&listID,
		&accountNames,
		flags,
	); err != nil {
		return err
	}

	if listID == "" {
		return missingIDError{
			target: cli.TargetList,
			action: "add the accounts to",
		}
	}

	if accountNames.Empty() {
		return zeroAccountNamesError{
			action: "add to the list",
		}
	}

	var accountIDs []string
	if err := client.Call(
		"GTSClient.GetMultipleAccountIDs",
		accountNames,
		&accountIDs,
	); err != nil {
		return fmt.Errorf("error retrieving the account IDs: %w", err)
	}

	if err := client.Call(
		"GTSClient.RemoveAccountsFromList",
		gtsclient.RemoveAccountsFromListArgs{
			ListID:     listID,
			AccountIDs: accountIDs,
		},
		nil,
	); err != nil {
		return fmt.Errorf("error removing the accounts from the list: %w", err)
	}

	printer.PrintSuccess(printSettings, "Successfully removed the account(s) from the list.")

	return nil
}
