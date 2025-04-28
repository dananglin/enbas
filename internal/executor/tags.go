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

func tagsFunc(
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
		return tagsShow(client, printSettings, cmd.FocusedTargetFlags)
	default:
		return unsupportedActionError{action: cmd.Action, target: cli.TargetTags}
	}
}

func tagsShow(
	client *rpc.Client,
	printSettings printer.Settings,
	flags []string,
) error {
	var limit int

	// Parse the remaining flags.
	if err := cli.ParseTagsShowFlags(
		&limit,
		flags,
	); err != nil {
		return err
	}

	var list model.TagList
	if err := client.Call(
		"GTSClient.GetFollowedTags",
		limit,
		&list,
	); err != nil {
		return fmt.Errorf("error retrieving the list of followed tags: %w", err)
	}

	if len(list.Tags) > 0 {
		if err := printer.PrintTagList(printSettings, list); err != nil {
			return fmt.Errorf("error printing the list of followed tags: %w", err)
		}
	} else {
		printer.PrintInfo("This account is not following any tags.\n")
	}

	return nil
}
