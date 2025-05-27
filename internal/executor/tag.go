package executor

import (
	"fmt"
	"net/rpc"
	"strings"

	"codeflow.dananglin.me.uk/apollo/enbas/internal/cli"
	"codeflow.dananglin.me.uk/apollo/enbas/internal/command"
	"codeflow.dananglin.me.uk/apollo/enbas/internal/config"
	"codeflow.dananglin.me.uk/apollo/enbas/internal/gtsclient"
	"codeflow.dananglin.me.uk/apollo/enbas/internal/model"
	"codeflow.dananglin.me.uk/apollo/enbas/internal/printer"
	"codeflow.dananglin.me.uk/apollo/enbas/internal/server"
)

// tagFunc is the function for the tag target for interacting
// with a single hashtag.
func tagFunc(
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
	case cli.ActionFollow:
		return tagFollow(session.Client(), printSettings, cmd.FocusedTargetFlags)
	case cli.ActionUnfollow:
		return tagUnfollow(session.Client(), printSettings, cmd.FocusedTargetFlags)
	case cli.ActionShow:
		return tagShow(session.Client(), printSettings, cmd.FocusedTargetFlags)
	case cli.ActionFind:
		return tagFind(session.Client(), printSettings, cmd.FocusedTargetFlags)
	default:
		return unsupportedActionError{action: cmd.Action, target: cli.TargetTag}
	}
}

func tagFollow(
	client *rpc.Client,
	printSettings printer.Settings,
	flags []string,
) error {
	var tagName string

	// Parse the remaining flags.
	if err := cli.ParseTagFollowFlags(
		&tagName,
		flags,
	); err != nil {
		return err //nolint:wrapcheck
	}

	if tagName == "" {
		return missingValueError{
			valueType: "name",
			target:    cli.TargetTag,
			action:    cli.ActionFollow,
		}
	}

	tagName = strings.TrimLeft(tagName, "#")

	if err := client.Call("GTSClient.FollowTag", tagName, nil); err != nil {
		return fmt.Errorf("error following the tag: %w", err)
	}

	printer.PrintSuccess(printSettings, "You are now following '"+tagName+"'.")

	return nil
}

func tagUnfollow(
	client *rpc.Client,
	printSettings printer.Settings,
	flags []string,
) error {
	var tagName string

	// Parse the remaining flags.
	if err := cli.ParseTagUnfollowFlags(
		&tagName,
		flags,
	); err != nil {
		return err //nolint:wrapcheck
	}

	if tagName == "" {
		return missingValueError{
			valueType: "name",
			target:    cli.TargetTag,
			action:    cli.ActionUnfollow,
		}
	}

	tagName = strings.TrimLeft(tagName, "#")

	if err := client.Call("GTSClient.UnfollowTag", tagName, nil); err != nil {
		return fmt.Errorf("error unfollowing the tag: %w", err)
	}

	printer.PrintSuccess(printSettings, "Successfully unfollowed '"+tagName+"'.")

	return nil
}

func tagShow(
	client *rpc.Client,
	printSettings printer.Settings,
	flags []string,
) error {
	var tagName string

	// Parse the remaining flags.
	if err := cli.ParseTagShowFlags(
		&tagName,
		flags,
	); err != nil {
		return err //nolint:wrapcheck
	}

	if tagName == "" {
		return missingValueError{
			valueType: "name",
			target:    cli.TargetTag,
			action:    cli.ActionShow,
		}
	}

	tagName = strings.TrimLeft(tagName, "#")

	var tag model.Tag
	if err := client.Call("GTSClient.GetTag", tagName, &tag); err != nil {
		return fmt.Errorf("error retrieving the details of the tag: %w", err)
	}

	if err := printer.PrintTag(printSettings, tag); err != nil {
		return fmt.Errorf("error printing the tag details: %w", err)
	}

	return nil
}

func tagFind(
	client *rpc.Client,
	printSettings printer.Settings,
	flags []string,
) error {
	var (
		query string
		limit int
	)

	// Parse the remaining flags.
	if err := cli.ParseTagFindFlags(
		&query,
		&limit,
		flags,
	); err != nil {
		return err //nolint:wrapcheck
	}

	if query == "" {
		return missingSearchQueryError{}
	}

	var results model.TagList

	if err := client.Call(
		"GTSClient.SearchTags",
		gtsclient.SearchTagsArgs{
			Limit: limit,
			Query: query,
		},
		&results,
	); err != nil {
		return fmt.Errorf("error searching for tags: %w", err)
	}

	if err := printer.PrintTagList(printSettings, results); err != nil {
		return fmt.Errorf("error printing the search result: %w", err)
	}

	return nil
}
