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

func filtersFunc(
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
		return filtersShow(session.Client(), printSettings)
	default:
		return unsupportedActionError{action: cmd.Action, target: cli.TargetFilters}
	}
}

func filtersShow(
	client *rpc.Client,
	printSettings printer.Settings,
) error {
	var filters []model.FilterV2
	if err := client.Call(
		"GTSClient.GetAllFilters",
		gtsclient.NoRPCArgs{},
		&filters,
	); err != nil {
		return fmt.Errorf("error retrieving the list of filters: %w", err)
	}

	if len(filters) == 0 {
		printer.PrintInfo("You have no filters.\n")

		return nil
	}

	if err := printer.PrintFilters(printSettings, filters); err != nil {
		return fmt.Errorf("error printing the list of filters: %w", err)
	}

	return nil
}
