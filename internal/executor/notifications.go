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

// notificationsFunc is the function for the notification target for
// interacting with multiple notifications.
func notificationsFunc(
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
	case cli.ActionClear:
		return notificationsClear(session.Client(), printSettings)
	case cli.ActionShow:
		return notificationsShow(session.Client(), printSettings, cmd.FocusedTargetFlags)
	default:
		return unsupportedActionError{action: cmd.Action, target: cli.TargetNotifications}
	}
}

func notificationsClear(
	client *rpc.Client,
	printSettings printer.Settings,
) error {
	if err := client.Call(
		"GTSClient.DeleteNotifications",
		gtsclient.NoRPCArgs{},
		nil,
	); err != nil {
		return fmt.Errorf("error deleting the notifications: %w", err)
	}

	printer.PrintSuccess(printSettings, "You have successfully cleared your notifications.")

	return nil
}

func notificationsShow(
	client *rpc.Client,
	printSettings printer.Settings,
	flags []string,
) error {
	var (
		limit                   int
		excludeNotificationType internalFlag.MultiEnumValue
		includeNotificationType internalFlag.MultiEnumValue
	)

	// Parse the remaining flags.
	if err := cli.ParseNotificationsShowFlags(
		&limit,
		&excludeNotificationType,
		&includeNotificationType,
		flags,
	); err != nil {
		return err
	}

	var myAccountID string
	if err := client.Call("GTSClient.GetMyAccountID", gtsclient.NoRPCArgs{}, &myAccountID); err != nil {
		return fmt.Errorf("unable to get your account ID: %w", err)
	}

	var notificationList []model.Notification
	if err := client.Call(
		"GTSClient.GetNotificationList",
		gtsclient.GetNotificationListArgs{
			Limit:        limit,
			ExcludeTypes: excludeNotificationType.Values(),
			IncludeTypes: includeNotificationType.Values(),
		},
		&notificationList,
	); err != nil {
		return fmt.Errorf(
			"error getting the list of notifications: %w",
			err,
		)
	}

	if len(notificationList) > 0 {
		if err := printer.PrintNotificationList(
			printSettings,
			notificationList,
			myAccountID,
		); err != nil {
			return fmt.Errorf("error printing the list of notifications: %w", err)
		}
	} else {
		printer.PrintInfo("You have no notifications.\n")
	}

	return nil
}
