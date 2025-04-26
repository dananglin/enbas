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

func timelineFunc(
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
		return timelineShow(client, printSettings, cmd.FocusedTargetFlags)
	default:
		return unsupportedActionError{action: cmd.Action, target: cli.TargetTimeline}
	}
}

func timelineShow(
	client *rpc.Client,
	printSettings printer.Settings,
	flags []string,
) error {
	var (
		err              error
		limit            int
		listID           string
		tagName          string
		timeline         model.StatusList
		timelineCategory string
	)

	// Parse the remaining flags.
	if err := cli.ParseTimelineShowFlags(
		&limit,
		&listID,
		&tagName,
		&timelineCategory,
		flags,
	); err != nil {
		return err
	}

	switch timelineCategory {
	case model.TimelineCategoryHome:
		err = client.Call("GTSClient.GetHomeTimeline", limit, &timeline)
	case model.TimelineCategoryPublic:
		err = client.Call("GTSClient.GetPublicTimeline", limit, &timeline)
	case model.TimelineCategoryList:
		if listID == "" {
			return missingIDError{
				target: cli.TargetList,
				action: "show in the timeline",
			}
		}

		var list model.List

		if err := client.Call("GTSClient.GetList", listID, &list); err != nil {
			return fmt.Errorf("unable to retrieve the list: %w", err)
		}

		err = client.Call(
			"GTSClient.GetListTimeline",
			gtsclient.GetListTimelineArgs{
				ListID: list.ID,
				Title:  list.Title,
				Limit:  limit,
			},
			&timeline,
		)
	case model.TimelineCategoryTag:
		if tagName == "" {
			return missingTagNameError{
				action: "view the timeline in",
			}
		}

		err = client.Call(
			"GTSClient.GetTagTimeline",
			gtsclient.GetTagTimelineArgs{
				TagName: tagName,
				Limit:   limit,
			},
			&timeline,
		)
	default:
		return model.InvalidTimelineCategoryError{Value: timelineCategory}
	}

	if err != nil {
		return fmt.Errorf("error retrieving the timeline: %w", err)
	}

	if len(timeline.Statuses) == 0 {
		printer.PrintInfo("There are no statuses in this timeline.\n")

		return nil
	}

	var myAccountID string
	if err := client.Call("GTSClient.GetMyAccountID", gtsclient.NoRPCArgs{}, &myAccountID); err != nil {
		return fmt.Errorf("unable to get your account ID: %w", err)
	}

	if err := printer.PrintStatusList(printSettings, timeline, myAccountID); err != nil {
		return fmt.Errorf("error printing the timeline: %w", err)
	}

	return nil
}
