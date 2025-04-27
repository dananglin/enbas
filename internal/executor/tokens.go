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

// tokensFunc is the function for the tokens target for interacting
// with the user's list of tokens.
func tokensFunc(
	opts topLevelOpts,
	cmd command.Command,
) error {
	// Load the configuration from file.
	cfg, err := config.NewConfigFromFile(opts.configDir)
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
	client, err := server.Connect(cfg.Server, opts.configDir)
	if err != nil {
		return fmt.Errorf("error creating the client for the daemon process: %w", err)
	}
	defer client.Close()

	switch cmd.Action {
	case cli.ActionShow:
		return tokensShow(
			client,
			printSettings,
			cmd.FocusedTargetFlags,
		)
	default:
		return unsupportedActionError{action: cmd.Action, target: cli.TargetTokens}
	}
}

func tokensShow(
	client *rpc.Client,
	printSettings printer.Settings,
	flags []string,
) error {
	var limit int

	// Parse the remaining flags
	if err := cli.ParseTokensShowFlags(
		&limit,
		flags,
	); err != nil {
		return err
	}

	var list model.TokenList
	if err := client.Call(
		"GTSClient.GetTokens",
		limit,
		&list,
	); err != nil {
		return fmt.Errorf("error retrieving the list of tokens: %w", err)
	}

	if len(list.Tokens) > 0 {
		if err := printer.PrintTokenList(printSettings, list); err != nil {
			return fmt.Errorf("error printing the list of tokens: %w", err)
		}
	} else {
		printer.PrintInfo("You have no tokens.\n")
	}

	return nil
}
