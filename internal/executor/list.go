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

// listFunc is the function for the list target for interacting
// with a single list.
func listFunc(
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
	case cli.ActionCreate:
		return listCreate(client, printSettings, cmd.FocusedTargetFlags)
	case cli.ActionDelete:
		return listDelete(client, printSettings, cmd.FocusedTargetFlags)
	case cli.ActionEdit:
		return listEdit(client, printSettings, cmd.FocusedTargetFlags)
	case cli.ActionShow:
		return listShow(client, printSettings, cmd.FocusedTargetFlags)
	default:
		return unsupportedActionError{action: cmd.Action, target: cli.TargetList}
	}
}

func listShow(
	client *rpc.Client,
	printSettings printer.Settings,
	flags []string,
) error {
	var listID string

	// Parse the remaining flags.
	if err := cli.ParseListShowFlags(
		&listID,
		flags,
	); err != nil {
		return err //nolint:wrapcheck
	}

	if listID == "" {
		return missingIDError{
			target: cli.TargetList,
			action: cli.ActionShow,
		}
	}

	var list model.List

	if err := client.Call("GTSClient.GetList", listID, &list); err != nil {
		return fmt.Errorf("unable to retrieve the list: %w", err)
	}

	acctMap, err := getAccountsFromList(client, listID)
	if err != nil {
		return err
	}

	if len(acctMap) > 0 {
		list.Accounts = acctMap
	}

	if err := printer.PrintList(printSettings, list); err != nil {
		return fmt.Errorf("error printing the list: %w", err)
	}

	return nil
}

func listCreate(
	client *rpc.Client,
	printSettings printer.Settings,
	flags []string,
) error {
	var (
		exclusive     bool
		repliesPolicy internalFlag.EnumValue
		title         string
	)

	// Parse the remaining flags.
	if err := cli.ParseListCreateFlags(
		&exclusive,
		&repliesPolicy,
		&title,
		flags,
	); err != nil {
		return err //nolint:wrapcheck
	}

	var list model.List
	if err := client.Call(
		"GTSClient.CreateList",
		gtsclient.CreateListArgs{
			Title:         title,
			RepliesPolicy: repliesPolicy.Value(),
			Exclusive:     exclusive,
		},
		&list,
	); err != nil {
		return fmt.Errorf("unable to create the list: %w", err)
	}

	printer.PrintSuccess(printSettings, "Successfully created the following list:")

	if err := printer.PrintList(printSettings, list); err != nil {
		return fmt.Errorf("error printing the list: %w", err)
	}

	return nil
}

func listEdit(
	client *rpc.Client,
	printSettings printer.Settings,
	flags []string,
) error {
	var (
		listID        string
		exclusive     internalFlag.BoolValue
		repliesPolicy internalFlag.EnumValue
		title         string
	)

	// Parse the remaining flags.
	if err := cli.ParseListEditFlags(
		&listID,
		&exclusive,
		&repliesPolicy,
		&title,
		flags,
	); err != nil {
		return err //nolint:wrapcheck
	}

	if listID == "" {
		return missingIDError{
			target: cli.TargetList,
			action: cli.ActionEdit,
		}
	}

	var listToUpdate model.List
	if err := client.Call("GTSClient.GetList", listID, &listToUpdate); err != nil {
		return fmt.Errorf("unable to get the list: %w", err)
	}

	if title != "" {
		listToUpdate.Title = title
	}

	if repliesPolicy.Value() != "" {
		listToUpdate.RepliesPolicy = repliesPolicy.Value()
	}

	if exclusive.IsSet() {
		listToUpdate.Exclusive = exclusive.Value()
	}

	var updatedList model.List
	if err := client.Call("GTSClient.UpdateList", listToUpdate, &updatedList); err != nil {
		return fmt.Errorf("error updating the list: %w", err)
	}

	acctMap, err := getAccountsFromList(client, updatedList.ID)
	if err != nil {
		return err
	}

	if len(acctMap) > 0 {
		updatedList.Accounts = acctMap
	}

	printer.PrintSuccess(printSettings, "Successfully edited the list.")

	if err := printer.PrintList(printSettings, updatedList); err != nil {
		return fmt.Errorf("error printing the list: %w", err)
	}

	return nil
}

func listDelete(
	client *rpc.Client,
	printSettings printer.Settings,
	flags []string,
) error {
	var listID string

	// Parse the remaining flags.
	if err := cli.ParseListDeleteFlags(
		&listID,
		flags,
	); err != nil {
		return err //nolint:wrapcheck
	}

	if listID == "" {
		return missingIDError{
			target: cli.TargetList,
			action: cli.ActionDelete,
		}
	}

	if err := client.Call("GTSClient.DeleteList", listID, nil); err != nil {
		return fmt.Errorf("unable to delete the list: %w", err)
	}

	printer.PrintSuccess(printSettings, "The list was successfully deleted.")

	return nil
}
