package executor

import (
	"flag"
	"fmt"

	"codeflow.dananglin.me.uk/apollo/enbas/internal/client"
	"codeflow.dananglin.me.uk/apollo/enbas/internal/config"
	"codeflow.dananglin.me.uk/apollo/enbas/internal/printer"
)

type FollowOrUnfollowExecutor struct {
	*flag.FlagSet

	printer      *printer.Printer
	config       *config.Config
	resourceType string
	accountName  string
	showReposts  bool
	notify       bool
	action       string
}

func NewFollowOrUnfollowExecutor(printer *printer.Printer, config *config.Config, name, summary string) *FollowOrUnfollowExecutor {
	command := FollowOrUnfollowExecutor{
		FlagSet: flag.NewFlagSet(name, flag.ExitOnError),

		printer: printer,
		config:  config,
		action:  name,
	}

	command.StringVar(&command.resourceType, flagType, "", "Specify the type of resource to follow")
	command.StringVar(&command.accountName, flagAccountName, "", "Specify the account name in full (username@domain)")
	command.BoolVar(&command.showReposts, flagShowReposts, true, "Show reposts from the account you want to follow")
	command.BoolVar(&command.notify, flagNotify, false, "Get notifications when the account you want to follow posts a status")

	command.Usage = commandUsageFunc(name, summary, command.FlagSet)

	return &command
}

func (f *FollowOrUnfollowExecutor) Execute() error {
	funcMap := map[string]func(*client.Client) error{
		resourceAccount: f.followOrUnfollowAccount,
	}

	doFunc, ok := funcMap[f.resourceType]
	if !ok {
		return UnsupportedTypeError{resourceType: f.resourceType}
	}

	gtsClient, err := client.NewClientFromFile(f.config.CredentialsFile)
	if err != nil {
		return fmt.Errorf("unable to create the GoToSocial client: %w", err)
	}

	return doFunc(gtsClient)
}

func (f *FollowOrUnfollowExecutor) followOrUnfollowAccount(gtsClient *client.Client) error {
	accountID, err := getAccountID(gtsClient, false, f.accountName, f.config.CredentialsFile)
	if err != nil {
		return fmt.Errorf("received an error while getting the account ID: %w", err)
	}

	switch f.action {
	case CommandFollow:
		return f.followAccount(gtsClient, accountID)
	case CommandUnfollow:
		return f.unfollowAccount(gtsClient, accountID)
	default:
		return nil
	}
}

func (f *FollowOrUnfollowExecutor) followAccount(gtsClient *client.Client, accountID string) error {
	form := client.FollowAccountForm{
		AccountID:   accountID,
		ShowReposts: f.showReposts,
		Notify:      f.notify,
	}

	if err := gtsClient.FollowAccount(form); err != nil {
		return fmt.Errorf("unable to follow the account: %w", err)
	}

	f.printer.PrintSuccess("Successfully sent the follow request.")

	return nil
}

func (f *FollowOrUnfollowExecutor) unfollowAccount(gtsClient *client.Client, accountID string) error {
	if err := gtsClient.UnfollowAccount(accountID); err != nil {
		return fmt.Errorf("unable to unfollow the account: %w", err)
	}

	f.printer.PrintSuccess("Successfully unfollowed the account.")

	return nil
}
