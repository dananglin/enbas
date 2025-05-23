package executor

import (
	"fmt"
	"net/rpc"
	"time"

	"codeflow.dananglin.me.uk/apollo/enbas/internal/cli"
	"codeflow.dananglin.me.uk/apollo/enbas/internal/command"
	"codeflow.dananglin.me.uk/apollo/enbas/internal/config"
	internalFlag "codeflow.dananglin.me.uk/apollo/enbas/internal/flag"
	"codeflow.dananglin.me.uk/apollo/enbas/internal/gtsclient"
	"codeflow.dananglin.me.uk/apollo/enbas/internal/model"
	"codeflow.dananglin.me.uk/apollo/enbas/internal/printer"
	"codeflow.dananglin.me.uk/apollo/enbas/internal/server"
)

// filterFunc is the function for the filter target for interacting
// with a single filter.
func filterFunc(
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
		return filterShow(client, printSettings, cmd.FocusedTargetFlags)
	case cli.ActionCreate:
		return filterCreate(client, printSettings, cmd.FocusedTargetFlags)
	case cli.ActionEdit:
		return filterEdit(client, printSettings, cmd.FocusedTargetFlags)
	case cli.ActionDelete:
		return filterDelete(client, printSettings, cmd.FocusedTargetFlags)
	default:
		return unsupportedActionError{action: cmd.Action, target: cli.TargetFilter}
	}
}

func filterShow(
	client *rpc.Client,
	printSettings printer.Settings,
	flags []string,
) error {
	var filterID string

	// Parse the remaining flags.
	if err := cli.ParseFilterShowFlags(
		&filterID,
		flags,
	); err != nil {
		return err //nolint:wrapcheck
	}

	if filterID == "" {
		return missingIDError{
			target: cli.TargetFilter,
			action: cli.ActionShow,
		}
	}

	var filter model.FilterV2

	if err := client.Call(
		"GTSClient.GetFilter",
		filterID,
		&filter,
	); err != nil {
		return fmt.Errorf("error retrieving the filter: %w", err)
	}

	if err := printer.PrintFilter(printSettings, filter); err != nil {
		return fmt.Errorf("error printing the filter: %w", err)
	}

	return nil
}

func filterCreate(
	client *rpc.Client,
	printSettings printer.Settings,
	flags []string,
) error {
	var (
		title     string
		contexts  internalFlag.MultiEnumValue
		expiresIn = internalFlag.NewTimeDurationValue(time.Duration(0))
		action    internalFlag.EnumValue
	)

	// Parse the remaining flags.
	if err := cli.ParseFilterCreateFlags(
		&title,
		&contexts,
		&expiresIn,
		&action,
		flags,
	); err != nil {
		return err //nolint:wrapcheck
	}

	if title == "" {
		return missingValueError{
			valueType: "title",
			target:    cli.TargetFilter,
			action:    cli.ActionCreate,
		}
	}

	if contexts.Empty() {
		return zeroValuesError{
			valueType: "filter context",
			action:    cli.ActionCreate + " the " + cli.TargetFilter + " with",
		}
	}

	var filter model.FilterV2
	if err := client.Call(
		"GTSClient.CreateFilter",
		gtsclient.CreateFilterArgs{
			Title:        title,
			FilterAction: action.Value(),
			Context:      contexts.Values(),
			ExpiresIn:    expiresIn.Value(),
		},
		&filter,
	); err != nil {
		return fmt.Errorf("error creating the filter: %w", err)
	}

	printer.PrintSuccess(printSettings, "Successfully created the filter with ID: "+filter.ID)

	return nil
}

func filterEdit(
	client *rpc.Client,
	printSettings printer.Settings,
	flags []string,
) error {
	var (
		filterID  string
		title     string
		contexts  internalFlag.MultiEnumValue
		expiresIn = internalFlag.NewTimeDurationValue(time.Duration(0))
		action    internalFlag.EnumValue
	)

	// Parse the remaining flags.
	if err := cli.ParseFilterEditFlags(
		&filterID,
		&title,
		&contexts,
		&expiresIn,
		&action,
		flags,
	); err != nil {
		return err //nolint:wrapcheck
	}

	if filterID == "" {
		return missingIDError{
			target: cli.TargetFilter,
			action: cli.ActionEdit,
		}
	}

	var filter model.FilterV2
	if err := client.Call(
		"GTSClient.GetFilter",
		filterID,
		&filter,
	); err != nil {
		return fmt.Errorf("error retrieving the existing filter: %w", err)
	}

	var editArgs gtsclient.EditFilterArgs
	editArgs.FilterID = filterID

	if title != "" {
		editArgs.Title = title
	} else {
		editArgs.Title = filter.Title
	}

	if action.Value() != "" {
		editArgs.FilterAction = action.Value()
	} else {
		editArgs.FilterAction = filter.Action
	}

	if !contexts.Empty() {
		editArgs.Context = contexts.Values()
	} else {
		editArgs.Context = filter.Context
	}

	switch {
	case expiresIn.IsSet():
		// Use the expiry set by the user.
		editArgs.ExpiresIn = expiresIn.Value()
	case filter.ExpiresAt.IsZero():
		// If the filter does not expire, make sure it stays that way.
		editArgs.ExpiresIn = time.Duration(0)
	default:
		// Calculate the remaining time before the filter's expiry.
		editArgs.ExpiresIn = time.Until(filter.ExpiresAt)
	}

	if err := client.Call(
		"GTSClient.EditFilter",
		editArgs,
		nil,
	); err != nil {
		return fmt.Errorf("error editing the filter: %w", err)
	}

	// The filter value from GTSClient.EditFilter does not include the list of
	// keywords or filtered statuses so we shall get the filter again to print
	// the results to the user.
	if err := client.Call(
		"GTSClient.GetFilter",
		filterID,
		&filter,
	); err != nil {
		return fmt.Errorf("error retrieving the existing filter: %w", err)
	}

	printer.PrintSuccess(printSettings, "Successfully edited the filter.")

	if err := printer.PrintFilter(printSettings, filter); err != nil {
		return fmt.Errorf("error printing the updated filter: %w", err)
	}

	return nil
}

func filterDelete(
	client *rpc.Client,
	printSettings printer.Settings,
	flags []string,
) error {
	var filterID string

	// Parse the remaining flags.
	if err := cli.ParseFilterDeleteFlags(&filterID, flags); err != nil {
		return err //nolint:wrapcheck
	}

	if filterID == "" {
		return missingIDError{
			target: cli.TargetFilter,
			action: cli.ActionDelete,
		}
	}

	if err := client.Call("GTSClient.DeleteFilter", filterID, nil); err != nil {
		return fmt.Errorf("error deleting the filter: %w", err)
	}

	printer.PrintSuccess(printSettings, "The filter was successfully deleted.")

	return nil
}
