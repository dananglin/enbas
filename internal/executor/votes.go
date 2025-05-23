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

func votesFunc(
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
	case cli.ActionAdd:
		return votesAdd(
			client,
			printSettings,
			cmd.RelatedTarget,
			cmd.RelatedTargetFlags,
		)
	default:
		return unsupportedActionError{action: cmd.Action, target: cli.TargetVotes}
	}
}

func votesAdd(
	client *rpc.Client,
	printSettings printer.Settings,
	relatedTarget string,
	relatedTargetFlags []string,
) error {
	switch relatedTarget {
	case cli.TargetStatus:
		return votesAddToStatus(
			client,
			printSettings,
			relatedTargetFlags,
		)
	default:
		return unsupportedTargetToTargetError{
			action:        cli.ActionAdd,
			focusedTarget: cli.TargetVotes,
			preposition:   cli.TargetActionPreposition(cli.TargetVotes, cli.ActionAdd),
			relatedTarget: relatedTarget,
		}
	}
}

func votesAddToStatus(
	client *rpc.Client,
	printSettings printer.Settings,
	flags []string,
) error {
	var (
		statusID string
		votes    = internalFlag.NewMultiIntValue()
	)

	// Parse the remaining flags.
	if err := cli.ParseVotesAddToStatusFlags(
		&statusID,
		&votes,
		flags,
	); err != nil {
		return err
	}

	if statusID == "" {
		return missingIDError{
			target: cli.TargetStatus,
			action: "add the votes to",
		}
	}

	if votes.Empty() {
		return zeroValuesError{
			valueType: "poll options",
			action:    "vote for",
		}
	}

	var status model.Status
	if err := client.Call(
		"GTSClient.GetStatus",
		statusID,
		&status,
	); err != nil {
		return fmt.Errorf("unable to get the status: %w", err)
	}

	if status.Poll.ID == "" {
		return pollMissingError{}
	}

	if status.Poll.Expired {
		return pollClosedError{}
	}

	if !status.Poll.Multiple && !votes.ExpectedLength(1) {
		return pollNoMultipleChoiceError{}
	}

	var myAccountID string
	if err := client.Call(
		"GTSClient.GetMyAccountID",
		gtsclient.NoRPCArgs{},
		&myAccountID,
	); err != nil {
		return fmt.Errorf("unable to get your account ID: %w", err)
	}

	if status.Account.ID == myAccountID {
		return voteInOwnPollError{}
	}

	pollID := status.Poll.ID

	if err := client.Call(
		"GTSClient.VoteInPoll",
		gtsclient.VoteInPollArgs{
			PollID:  pollID,
			Choices: votes.Values(),
		},
		nil,
	); err != nil {
		return fmt.Errorf("unable to add your vote(s) to the poll: %w", err)
	}

	printer.PrintSuccess(printSettings, "Successfully added your vote(s) to the poll.")

	return nil
}
