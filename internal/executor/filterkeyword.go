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

func filterKeywordFunc(
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
		return filterKeywordAdd(
			client,
			printSettings,
			cmd.RelatedTarget,
			cmd.RelatedTargetFlags,
		)
	case cli.ActionDelete:
		return filterKeywordDelete(
			client,
			printSettings,
			cmd.FocusedTargetFlags,
		)
	case cli.ActionEdit:
		return filterKeywordEdit(
			client,
			printSettings,
			cmd.FocusedTargetFlags,
		)
	case cli.ActionShow:
		return filterKeywordShow(
			client,
			printSettings,
			cmd.FocusedTargetFlags,
		)
	default:
		return unsupportedActionError{action: cmd.Action, target: cli.TargetFilterKeyword}
	}
}

func filterKeywordAdd(
	client *rpc.Client,
	printSettings printer.Settings,
	relatedTarget string,
	relatedTargetFlags []string,
) error {
	switch relatedTarget {
	case cli.TargetFilter:
		return filterKeywordAddToFilter(
			client,
			printSettings,
			relatedTargetFlags,
		)
	default:
		return unsupportedTargetToTargetError{
			action:        cli.ActionAdd,
			focusedTarget: cli.TargetFilterKeyword,
			preposition: cli.TargetActionPreposition(
				cli.TargetFilterKeyword,
				cli.ActionAdd,
			),
			relatedTarget: relatedTarget,
		}
	}
}

func filterKeywordAddToFilter(
	client *rpc.Client,
	printSettings printer.Settings,
	flags []string,
) error {
	var (
		filterID        string
		keyword         string
		filterWholeWord bool
	)

	// Parse the remaining flags
	if err := cli.ParseFilterKeywordAddToFilterFlags(
		&filterID,
		&keyword,
		&filterWholeWord,
		flags,
	); err != nil {
		return err //nolint:wrapcheck
	}

	if filterID == "" {
		return missingIDError{
			target: cli.TargetFilter,
			action: "add the filter-keyword to",
		}
	}

	if keyword == "" {
		return missingValueError{
			valueType: "keyword",
			target:    cli.TargetFilterKeyword,
			action:    "add to the filter",
		}
	}

	var filterKeyword model.FilterKeyword

	if err := client.Call(
		"GTSClient.AddFilterKeywordToFilter",
		gtsclient.AddFilterKeywordToFilterArgs{
			FilterID:  filterID,
			Keyword:   keyword,
			WholeWord: filterWholeWord,
		},
		&filterKeyword,
	); err != nil {
		return fmt.Errorf("error adding the filter-keyword to the filter: %w", err)
	}

	printer.PrintSuccess(
		printSettings,
		"Successfully added the filter-keyword (ID: "+filterKeyword.ID+") to the filter",
	)

	return nil
}

func filterKeywordDelete(
	client *rpc.Client,
	printSettings printer.Settings,
	flags []string,
) error {
	var filterKeywordID string

	// Parse the remaining flags.
	if err := cli.ParseFilterKeywordDeleteFlags(
		&filterKeywordID,
		flags,
	); err != nil {
		return err //nolint:wrapcheck
	}

	if filterKeywordID == "" {
		return missingIDError{
			target: cli.TargetFilterKeyword,
			action: cli.ActionDelete,
		}
	}

	if err := client.Call(
		"GTSClient.DeleteFilterKeyword",
		filterKeywordID,
		nil,
	); err != nil {
		return fmt.Errorf("error deleting the filter-keyword: %w", err)
	}

	printer.PrintSuccess(printSettings, "The filter-keyword was successfully deleted.")

	return nil
}

func filterKeywordEdit(
	client *rpc.Client,
	printSettings printer.Settings,
	flags []string,
) error {
	var (
		filterKeywordID string
		keyword         string
		wholeWord       internalFlag.BoolValue
	)

	// Parse the remaining flags.
	if err := cli.ParseFilterKeywordEditFlags(
		&filterKeywordID,
		&keyword,
		&wholeWord,
		flags,
	); err != nil {
		return err //nolint:wrapcheck
	}

	if filterKeywordID == "" {
		return missingIDError{
			target: cli.TargetFilterKeyword,
			action: cli.ActionEdit,
		}
	}

	var filterKeyword model.FilterKeyword
	if err := client.Call(
		"GTSClient.GetFilterKeyword",
		filterKeywordID,
		&filterKeyword,
	); err != nil {
		return fmt.Errorf("error retrieving the filter-keyword: %w", err)
	}

	editArgs := gtsclient.UpdateFilterKeywordArgs{}
	editArgs.FilterKeywordID = filterKeywordID

	if keyword != "" {
		editArgs.Keyword = keyword
	} else {
		editArgs.Keyword = filterKeyword.Keyword
	}

	if wholeWord.IsSet() {
		editArgs.WholeWord = wholeWord.Value()
	} else {
		editArgs.WholeWord = filterKeyword.WholeWord
	}

	if err := client.Call(
		"GTSClient.UpdateFilterKeyword",
		editArgs,
		&filterKeyword,
	); err != nil {
		return fmt.Errorf("error updating the filter-keyword: %w", err)
	}

	printer.PrintSuccess(printSettings, "Successfully edited the filter-keyword.")

	if err := printer.PrintFilterKeyword(printSettings, filterKeyword); err != nil {
		return fmt.Errorf("error printing the filter-keyword: %w", err)
	}

	return nil
}

func filterKeywordShow(
	client *rpc.Client,
	printSettings printer.Settings,
	flags []string,
) error {
	var filterKeywordID string

	// Parse the remaining flags.
	if err := cli.ParseFilterKeywordShowFlags(
		&filterKeywordID,
		flags,
	); err != nil {
		return err //nolint:wrapcheck
	}

	if filterKeywordID == "" {
		return missingIDError{
			target: cli.TargetFilterKeyword,
			action: cli.ActionShow,
		}
	}

	var filterKeyword model.FilterKeyword
	if err := client.Call(
		"GTSClient.GetFilterKeyword",
		filterKeywordID,
		&filterKeyword,
	); err != nil {
		return fmt.Errorf("error retrieving the filter-keyword: %w", err)
	}

	if err := printer.PrintFilterKeyword(printSettings, filterKeyword); err != nil {
		return fmt.Errorf("error printing the filter-keyword: %w", err)
	}

	return nil
}
