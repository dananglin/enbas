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

func instanceFunc(
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
		return instanceShow(client, printSettings)
	default:
		return unsupportedActionError{action: cmd.Action, target: cli.TargetInstance}
	}
}

func instanceShow(
	client *rpc.Client,
	printSettings printer.Settings,
) error {
	var instance model.InstanceV2
	if err := client.Call("GTSClient.GetInstance", gtsclient.NoRPCArgs{}, &instance); err != nil {
		return fmt.Errorf("unable to retrieve the instance details: %w", err)
	}

	if err := printer.PrintInstance(printSettings, instance); err != nil {
		return fmt.Errorf("error printing the instance details: %w", err)
	}

	return nil
}
