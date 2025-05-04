package executor

import (
	"fmt"
	"net/rpc"
	"time"

	"codeflow.dananglin.me.uk/apollo/enbas/internal/cli"
	"codeflow.dananglin.me.uk/apollo/enbas/internal/command"
	"codeflow.dananglin.me.uk/apollo/enbas/internal/config"
	internalFlag "codeflow.dananglin.me.uk/apollo/enbas/internal/flag"
	"codeflow.dananglin.me.uk/apollo/enbas/internal/gtsclient"
	"codeflow.dananglin.me.uk/apollo/enbas/internal/model"
	"codeflow.dananglin.me.uk/apollo/enbas/internal/printer"
	"codeflow.dananglin.me.uk/apollo/enbas/internal/server"
	"codeflow.dananglin.me.uk/apollo/enbas/internal/utilities"
)

// accountFunc is the function for the account target for
// interacting with a local or remote account.
func accountFunc(
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
		return accountShow(client, printSettings, cfg.Integrations.Browser, cmd.FocusedTargetFlags)
	case cli.ActionMute:
		return accountMute(client, printSettings, cmd.FocusedTargetFlags)
	case cli.ActionUnmute:
		return accountUnmute(client, printSettings, cmd.FocusedTargetFlags)
	case cli.ActionFollow:
		return accountFollow(client, printSettings, cmd.FocusedTargetFlags)
	case cli.ActionUnfollow:
		return accountUnfollow(client, printSettings, cmd.FocusedTargetFlags)
	case cli.ActionBlock:
		return accountBlock(client, printSettings, cmd.FocusedTargetFlags)
	case cli.ActionUnblock:
		return accountUnblock(client, printSettings, cmd.FocusedTargetFlags)
	case cli.ActionFind:
		return accountFind(client, printSettings, cmd.FocusedTargetFlags)
	default:
		return unsupportedActionError{action: cmd.Action, target: cli.TargetAccount}
	}
}

func accountShow(
	client *rpc.Client,
	printSettings printer.Settings,
	browser string,
	flags []string,
) error {
	var (
		accountName             string
		showInBrowser           bool
		excludeReblogs          bool
		excludeReplies          bool
		maxStatuses             int
		myAccount               bool
		onlyMedia               bool
		onlyPinned              bool
		onlyPublic              bool
		skipAccountRelationship bool
		skipUserPreferences     bool
		showStatuses            bool
	)

	// Parse the remaining flags.
	if err := cli.ParseAccountShowFlags(
		&accountName,
		&showInBrowser,
		&excludeReblogs,
		&excludeReplies,
		&maxStatuses,
		&myAccount,
		&onlyMedia,
		&onlyPinned,
		&onlyPublic,
		&skipAccountRelationship,
		&skipUserPreferences,
		&showStatuses,
		flags,
	); err != nil {
		return err
	}

	var account model.Account

	if myAccount {
		if err := client.Call("GTSClient.GetMyAccount", gtsclient.NoRPCArgs{}, &account); err != nil {
			return fmt.Errorf("unable to retrieve your account: %w", err)
		}
	} else {
		if accountName == "" {
			return missingAccountNameError{action: cli.ActionShow}
		}

		if err := client.Call("GTSClient.GetAccount", accountName, &account); err != nil {
			return fmt.Errorf("unable to get the account information: %w", err)
		}
	}

	if showInBrowser {
		if err := utilities.OpenLink(browser, account.URL); err != nil {
			return fmt.Errorf("unable to open link: %w", err)
		}

		return nil
	}

	var (
		relationship model.AccountRelationship
		preferences  model.Preferences
		statusList   model.StatusList
		myAccountID  string
	)

	relationship.Print = false
	preferences.Print = false
	statusList.Statuses = nil

	if !myAccount && !skipAccountRelationship {
		if err := client.Call("GTSClient.GetAccountRelationship", account.ID, &relationship); err != nil {
			return fmt.Errorf("unable to retrieve the relationship to this account: %w", err)
		}

		relationship.Print = true
	}

	if myAccount {
		myAccountID = account.ID
		if !skipUserPreferences {
			if err := client.Call("GTSClient.GetUserPreferences", gtsclient.NoRPCArgs{}, &preferences); err != nil {
				return fmt.Errorf("unable to retrieve the user preferences: %w", err)
			}

			preferences.Print = true
		}
	}

	if showStatuses {
		args := gtsclient.GetAccountStatusesArgs{
			AccountID:      account.ID,
			Limit:          maxStatuses,
			ExcludeReplies: excludeReplies,
			ExcludeReblogs: excludeReblogs,
			OnlyMedia:      onlyMedia,
			Pinned:         onlyPinned,
			OnlyPublic:     onlyPublic,
		}

		if err := client.Call("GTSClient.GetAccountStatuses", args, &statusList); err != nil {
			return fmt.Errorf("unable to retrieve the account's statuses: %w", err)
		}
	}

	if err := printer.PrintAccount(
		printSettings,
		account,
		relationship,
		preferences,
		statusList,
		myAccountID,
	); err != nil {
		return fmt.Errorf("error printing the account: %w", err)
	}

	return nil
}

func accountMute(
	client *rpc.Client,
	printSettings printer.Settings,
	flags []string,
) error {
	var (
		accountName       string
		duration          = internalFlag.NewTimeDurationValue(time.Duration(0))
		muteNotifications bool
	)

	// Parse the remaining flags
	if err := cli.ParseAccountMuteFlags(
		&accountName,
		&duration,
		&muteNotifications,
		flags,
	); err != nil {
		return err
	}

	if accountName == "" {
		return missingAccountNameError{action: cli.ActionMute}
	}

	var accountID string
	if err := client.Call("GTSClient.GetAccountID", accountName, &accountID); err != nil {
		return fmt.Errorf("received an error while getting the account ID: %w", err)
	}

	if err := client.Call(
		"GTSClient.MuteAccount",
		gtsclient.MuteAccountArgs{
			AccountID:     accountID,
			Notifications: muteNotifications,
			Duration:      int(duration.Value().Seconds()),
		},
		nil,
	); err != nil {
		return fmt.Errorf("error muting the account: %w", err)
	}

	printer.PrintSuccess(printSettings, "Successfully muted the account.")

	return nil
}

func accountFollow(
	client *rpc.Client,
	printSettings printer.Settings,
	flags []string,
) error {
	var (
		accountName string
		notify      bool
		showReblogs bool
	)

	// Parse the remaining flags
	if err := cli.ParseAccountFollowFlags(
		&accountName,
		&notify,
		&showReblogs,
		flags,
	); err != nil {
		return err
	}

	if accountName == "" {
		return missingAccountNameError{action: cli.ActionFollow}
	}

	var accountID string
	if err := client.Call("GTSClient.GetAccountID", accountName, &accountID); err != nil {
		return fmt.Errorf("received an error while getting the account ID: %w", err)
	}

	if err := client.Call(
		"GTSClient.FollowAccount",
		gtsclient.FollowAccountArgs{
			AccountID:   accountID,
			ShowReposts: showReblogs,
			Notify:      notify,
		},
		nil,
	); err != nil {
		return fmt.Errorf("error following the account: %w", err)
	}

	printer.PrintSuccess(printSettings, "Successfully sent the follow request.")

	return nil
}

func accountUnfollow(
	client *rpc.Client,
	printSettings printer.Settings,
	flags []string,
) error {
	var accountName string

	// Parse the remaining flags
	if err := cli.ParseAccountUnfollowFlags(
		&accountName,
		flags,
	); err != nil {
		return err
	}

	if accountName == "" {
		return missingAccountNameError{action: cli.ActionUnfollow}
	}

	var accountID string
	if err := client.Call("GTSClient.GetAccountID", accountName, &accountID); err != nil {
		return fmt.Errorf("received an error while getting the account ID: %w", err)
	}

	if err := client.Call("GTSClient.UnfollowAccount", accountID, nil); err != nil {
		return fmt.Errorf("unable to unfollow the account: %w", err)
	}

	printer.PrintSuccess(printSettings, "Successfully unfollowed the account.")

	return nil
}

func accountUnmute(
	client *rpc.Client,
	printSettings printer.Settings,
	flags []string,
) error {
	var accountName string

	// Parse the remaining flags
	if err := cli.ParseAccountUnmuteFlags(
		&accountName,
		flags,
	); err != nil {
		return err
	}

	if accountName == "" {
		return missingAccountNameError{action: cli.ActionUnmute}
	}

	var accountID string
	if err := client.Call("GTSClient.GetAccountID", accountName, &accountID); err != nil {
		return fmt.Errorf("received an error while getting the account ID: %w", err)
	}

	if err := client.Call("GTSClient.UnmuteAccount", accountID, nil); err != nil {
		return fmt.Errorf("error unmuting the account: %w", err)
	}

	printer.PrintSuccess(printSettings, "Successfully unmuted the account.")

	return nil
}

func accountBlock(
	client *rpc.Client,
	printSettings printer.Settings,
	flags []string,
) error {
	var accountName string

	// Parse the remaining flags
	if err := cli.ParseAccountBlockFlags(
		&accountName,
		flags,
	); err != nil {
		return err
	}

	if accountName == "" {
		return missingAccountNameError{action: cli.ActionBlock}
	}

	var accountID string
	if err := client.Call(
		"GTSClient.GetAccountID",
		accountName,
		&accountID,
	); err != nil {
		return fmt.Errorf("received an error while getting the account ID: %w", err)
	}

	if err := client.Call(
		"GTSClient.BlockAccount",
		accountID,
		nil,
	); err != nil {
		return fmt.Errorf("unable to block the account: %w", err)
	}

	printer.PrintSuccess(printSettings, "Successfully blocked the account.")

	return nil
}

func accountUnblock(
	client *rpc.Client,
	printSettings printer.Settings,
	flags []string,
) error {
	var accountName string

	// Parse the remaining flags
	if err := cli.ParseAccountUnblockFlags(
		&accountName,
		flags,
	); err != nil {
		return err
	}

	if accountName == "" {
		return missingAccountNameError{action: cli.ActionUnblock}
	}

	var accountID string
	if err := client.Call(
		"GTSClient.GetAccountID",
		accountName,
		&accountID,
	); err != nil {
		return fmt.Errorf("received an error while getting the account ID: %w", err)
	}

	if err := client.Call(
		"GTSClient.UnblockAccount",
		accountID,
		nil,
	); err != nil {
		return fmt.Errorf("unable to unblock the account: %w", err)
	}

	printer.PrintSuccess(printSettings, "Successfully unblocked the account.")

	return nil
}

func accountFind(
	client *rpc.Client,
	printSettings printer.Settings,
	flags []string,
) error {
	var (
		query                string
		limit                int
		restrictToFollowings bool
		resolve              bool
	)

	// Parse the remaining flags
	if err := cli.ParseAccountFindFlags(
		&query,
		&limit,
		&restrictToFollowings,
		&resolve,
		flags,
	); err != nil {
		return err
	}

	if query == "" {
		return missingSearchQueryError{}
	}

	var results model.AccountList

	if err := client.Call(
		"GTSClient.SearchAccounts",
		gtsclient.SearchAccountsArgs{
			Limit:     limit,
			Query:     query,
			Resolve:   resolve,
			Following: restrictToFollowings,
		},
		&results,
	); err != nil {
		return fmt.Errorf("error searching for accounts: %w", err)
	}

	if err := printer.PrintAccountList(printSettings, results); err != nil {
		return fmt.Errorf("error printing the search result: %w", err)
	}

	return nil
}
