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

func filterStatusFunc(
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
		return filterStatusAdd(
			client,
			printSettings,
			cmd.RelatedTarget,
			cmd.RelatedTargetFlags,
		)
	case cli.ActionDelete:
		return filterStatusDelete(
			client,
			printSettings,
			cmd.FocusedTargetFlags,
		)
	case cli.ActionShow:
		return filterStatusShow(
			client,
			printSettings,
			cmd.FocusedTargetFlags,
		)
	default:
		return unsupportedActionError{action: cmd.Action, target: cli.TargetFilterStatus}
	}
}

func filterStatusAdd(
	client *rpc.Client,
	printSettings printer.Settings,
	relatedTarget string,
	relatedTargetFlags []string,
) error {
	switch relatedTarget {
	case cli.TargetFilter:
		return filterStatusAddToFilter(
			client,
			printSettings,
			relatedTargetFlags,
		)
	default:
		return unsupportedTargetToTargetError{
			action:        cli.ActionAdd,
			focusedTarget: cli.TargetFilterStatus,
			preposition: cli.TargetActionPreposition(
				cli.TargetFilterStatus,
				cli.ActionAdd,
			),
			relatedTarget: relatedTarget,
		}
	}
}

func filterStatusAddToFilter(
	client *rpc.Client,
	printSettings printer.Settings,
	flags []string,
) error {
	var (
		filterID string
		statusID string
	)

	// Parse the remaining flags
	if err := cli.ParseFilterStatusAddToFilterFlags(
		&filterID,
		&statusID,
		flags,
	); err != nil {
		return err //nolint:wrapcheck
	}

	if filterID == "" {
		return missingIDError{
			target: cli.TargetFilter,
			action: "add the filter-status to",
		}
	}

	if statusID == "" {
		return missingIDError{
			target: cli.TargetStatus,
			action: "filter",
		}
	}

	var filterStatus model.FilterStatus

	if err := client.Call(
		"GTSClient.AddFilterStatusToFilter",
		gtsclient.AddFilterStatusToFilterArgs{
			FilterID: filterID,
			StatusID: statusID,
		},
		&filterStatus,
	); err != nil {
		return fmt.Errorf("error adding the filter-status to the filter: %w", err)
	}

	printer.PrintSuccess(
		printSettings,
		"Successfully added the filter-status (ID: "+filterStatus.ID+") to the filter",
	)

	return nil
}

func filterStatusDelete(
	client *rpc.Client,
	printSettings printer.Settings,
	flags []string,
) error {
	var filterStatusID string

	// Parse the remaining flags.
	if err := cli.ParseFilterStatusDeleteFlags(
		&filterStatusID,
		flags,
	); err != nil {
		return err //nolint:wrapcheck
	}

	if filterStatusID == "" {
		return missingIDError{
			target: cli.TargetFilterStatus,
			action: cli.ActionDelete,
		}
	}

	if err := client.Call(
		"GTSClient.DeleteFilterStatus",
		filterStatusID,
		nil,
	); err != nil {
		return fmt.Errorf("error deleting the filter-status: %w", err)
	}

	printer.PrintSuccess(printSettings, "The filter-keyword was successfully deleted.")

	return nil
}

func filterStatusShow(
	client *rpc.Client,
	printSettings printer.Settings,
	flags []string,
) error {
	var filterStatusID string

	if err := cli.ParseFilterStatusShowFlags(
		&filterStatusID,
		flags,
	); err != nil {
		return err //nolint:wrapcheck
	}

	if filterStatusID == "" {
		return missingIDError{
			target: cli.TargetFilterStatus,
			action: cli.ActionShow,
		}
	}

	var filterStatus model.FilterStatus
	if err := client.Call(
		"GTSClient.GetFilterStatus",
		filterStatusID,
		&filterStatus,
	); err != nil {
		return fmt.Errorf("error retrieving the filter-status: %w", err)
	}

	if err := printer.PrintFilterStatus(printSettings, filterStatus); err != nil {
		return fmt.Errorf("error printing the filter-status: %w", err)
	}

	return nil
}
