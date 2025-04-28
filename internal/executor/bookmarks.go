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

// bookmarksFunc is the function for the bookmarks target for interacting
// with the user's bookmarks.
func bookmarksFunc(
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
		return bookmarksShow(
			client,
			printSettings,
			cmd.FocusedTargetFlags,
		)
	default:
		return unsupportedActionError{action: cmd.Action, target: cli.TargetBookmarks}
	}
}

func bookmarksShow(
	client *rpc.Client,
	printSettings printer.Settings,
	flags []string,
) error {
	var limit int

	// Parse the remaining flags.
	if err := cli.ParseBookmarksShowFlags(
		&limit,
		flags,
	); err != nil {
		return err
	}

	var bookmarks model.StatusList
	if err := client.Call("GTSClient.GetBookmarks", limit, &bookmarks); err != nil {
		return fmt.Errorf("error retrieving the list of your bookmarks: %w", err)
	}

	if len(bookmarks.Statuses) == 0 {
		printer.PrintInfo("You have no bookmarks.\n")

		return nil
	}

	var myAccountID string
	if err := client.Call("GTSClient.GetMyAccountID", gtsclient.NoRPCArgs{}, &myAccountID); err != nil {
		return fmt.Errorf("unable to get your account ID: %w", err)
	}

	if err := printer.PrintStatusList(printSettings, bookmarks, myAccountID); err != nil {
		return fmt.Errorf("error printing the list of your bookmarks: %w", err)
	}

	return nil
}
