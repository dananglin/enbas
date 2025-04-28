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

// listsFunc is the function for the lists target for interacting
// with multiple lists.
func listsFunc(
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
		return listsShow(client, printSettings)
	default:
		return unsupportedActionError{action: cmd.Action, target: cli.TargetLists}
	}
}

func listsShow(
	client *rpc.Client,
	printSettings printer.Settings,
) error {
	var lists []model.List
	if err := client.Call(
		"GTSClient.GetAllLists",
		gtsclient.NoRPCArgs{},
		&lists,
	); err != nil {
		return fmt.Errorf("unable to retrieve the lists: %w", err)
	}

	if len(lists) == 0 {
		printer.PrintInfo("You have no lists.\n")

		return nil
	}

	if err := printer.PrintLists(printSettings, lists); err != nil {
		return fmt.Errorf("error printing the set of lists: %w", err)
	}

	return nil
}
