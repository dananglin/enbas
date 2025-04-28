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

// notificationFunc is the function for the notification target for
// interacting with a single notification.
func notificationFunc(
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
		return notificationShow(client, printSettings, cmd.FocusedTargetFlags)
	default:
		return unsupportedActionError{action: cmd.Action, target: cli.TargetNotification}
	}
}

func notificationShow(
	client *rpc.Client,
	printSettings printer.Settings,
	flags []string,
) error {
	var notificationID string

	// Parse the remaining flags.
	if err := cli.ParseNotificationShowFlags(
		&notificationID,
		flags,
	); err != nil {
		return err
	}

	if notificationID == "" {
		return missingIDError{
			target: cli.TargetNotification,
			action: cli.ActionShow,
		}
	}

	var myAccountID string
	if err := client.Call("GTSClient.GetMyAccountID", gtsclient.NoRPCArgs{}, &myAccountID); err != nil {
		return fmt.Errorf("unable to get your account ID: %w", err)
	}

	var notification model.Notification
	if err := client.Call(
		"GTSClient.GetNotification",
		notificationID,
		&notification,
	); err != nil {
		return fmt.Errorf("error retrieving the notification: %w", err)
	}

	if err := printer.PrintNotification(printSettings, notification, myAccountID); err != nil {
		return fmt.Errorf("error printing the notification: %w", err)
	}

	return nil
}
