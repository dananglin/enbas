package executor

import (
	"fmt"
	"net/rpc"
	"slices"

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
	case cli.ActionClear:
		return notificationsClear(client, printSettings)
	case cli.ActionShow:
		return notificationsShow(client, printSettings, cmd.FocusedTargetFlags)
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
		excludeNotificationType internalFlag.StringSliceValue
		includeNotificationType internalFlag.StringSliceValue
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

	for _, exclude := range slices.All(excludeNotificationType) {
		if _, err := model.ParseNotificationType(exclude); err != nil {
			return err
		}
	}

	for _, include := range slices.All(includeNotificationType) {
		if _, err := model.ParseNotificationType(include); err != nil {
			return err
		}
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
			IncludeTypes: []string(includeNotificationType),
			ExcludeTypes: []string(excludeNotificationType),
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
