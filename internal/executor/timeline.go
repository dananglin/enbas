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

func timelineFunc(
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
		return timelineShow(session.Client(), printSettings, cmd.FocusedTargetFlags)
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
		err      error
		limit    int
		listID   string
		tagName  string
		category internalFlag.EnumValue
	)

	// Parse the remaining flags.
	if err := cli.ParseTimelineShowFlags(
		&limit,
		&listID,
		&tagName,
		&category,
		flags,
	); err != nil {
		return err
	}

	var timeline model.StatusList

	switch category.Value() {
	case "home":
		err = client.Call("GTSClient.GetHomeTimeline", limit, &timeline)
	case "public":
		err = client.Call("GTSClient.GetPublicTimeline", limit, &timeline)
	case "list":
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
	case "tag":
		if tagName == "" {
			return missingValueError{
				valueType: "name",
				target:    cli.TargetTag,
				action:    "view the timeline in",
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
		return invalidTimelineCategoryError{category: category.Value()}
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
