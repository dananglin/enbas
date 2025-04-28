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

// favouritesFunc is the function for the 'favourites' target
// for interacting with favourite statuses.
func favouritesFunc(
	opts topLevelOpts,
	cmd command.Command,
) error {
	// Load the configuration from file.
	cfg, err := config.NewConfigFromFile(opts.configPath)
	if err != nil {
		return fmt.Errorf("unable to load configuration: %w", err)
	}

	// Create the print settings.
	printSettings := printer.NewSettings(
		opts.noColor,
		cfg.Integrations.Pager,
		cfg.LineWrapMaxWidth,
	)

	// Create the client to interact with the GoToSocial instance.
	client, err := server.Connect(cfg.Server, opts.configPath)
	if err != nil {
		return fmt.Errorf("error creating the client for the daemon process: %w", err)
	}
	defer client.Close()

	switch cmd.Action {
	case cli.ActionShow:
		return favouritesShow(client, printSettings, cmd.FocusedTargetFlags)
	default:
		return unsupportedActionError{action: cmd.Action, target: cli.TargetFavourites}
	}
}

func favouritesShow(
	client *rpc.Client,
	printSettings printer.Settings,
	flags []string,
) error {
	var limit int

	// Parse the remaining flags.
	if err := cli.ParseFavouritesShowFlags(
		&limit,
		flags,
	); err != nil {
		return err
	}

	var favourites model.StatusList
	if err := client.Call(
		"GTSClient.GetFavourites",
		limit,
		&favourites,
	); err != nil {
		return fmt.Errorf("error retrieving the list of your favourite statuses: %w", err)
	}

	if len(favourites.Statuses) > 0 {
		var myAccountID string
		if err := client.Call(
			"GTSClient.GetMyAccountID",
			gtsclient.NoRPCArgs{},
			&myAccountID,
		); err != nil {
			return fmt.Errorf("unable to get your account ID: %w", err)
		}

		if err := printer.PrintStatusList(printSettings, favourites, myAccountID); err != nil {
			return fmt.Errorf("error printing the list of your favourite statuses: %w", err)
		}
	} else {
		printer.PrintInfo("You have no favourite statuses.\n")
	}

	return nil
}
