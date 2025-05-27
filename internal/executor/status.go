package executor

import (
	"fmt"
	"net/rpc"
	"path/filepath"
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

// statusFunc is the function for the status target for interacting
// with the user's statuses.
func statusFunc(
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
	case cli.ActionCreate:
		return statusCreate(
			session.Client(),
			printSettings,
			cmd.FocusedTargetFlags,
		)
	case cli.ActionAdd:
		return statusAdd(
			session.Client(),
			printSettings,
			cmd.RelatedTarget,
			cmd.RelatedTargetFlags,
		)
	case cli.ActionRemove:
		return statusRemove(
			session.Client(),
			printSettings,
			cmd.RelatedTarget,
			cmd.RelatedTargetFlags,
		)
	case cli.ActionMute:
		return statusMute(
			session.Client(),
			printSettings,
			cmd.FocusedTargetFlags,
		)
	case cli.ActionUnmute:
		return statusUnmute(
			session.Client(),
			printSettings,
			cmd.FocusedTargetFlags,
		)
	case cli.ActionFavourite:
		return statusFavourite(
			session.Client(),
			printSettings,
			cmd.FocusedTargetFlags,
		)
	case cli.ActionUnfavourite:
		return statusUnfavourite(
			session.Client(),
			printSettings,
			cmd.FocusedTargetFlags,
		)
	case cli.ActionReblog:
		return statusReblog(
			session.Client(),
			printSettings,
			cmd.FocusedTargetFlags,
		)
	case cli.ActionUnreblog:
		return statusUnreblog(
			session.Client(),
			printSettings,
			cmd.FocusedTargetFlags,
		)
	case cli.ActionShow:
		return statusShow(
			session.Client(),
			printSettings,
			cfg.Integrations.Browser,
			cmd.FocusedTargetFlags,
		)
	case cli.ActionDelete:
		return statusDelete(
			session.Client(),
			printSettings,
			cfg.CacheDirectory,
			cmd.FocusedTargetFlags,
		)
	case cli.ActionFind:
		return statusFind(
			session.Client(),
			printSettings,
			cmd.FocusedTargetFlags,
		)
	default:
		return unsupportedActionError{action: cmd.Action, target: cli.TargetStatus}
	}
}

func statusCreate(
	client *rpc.Client,
	printSettings printer.Settings,
	flags []string,
) error {
	var (
		addPoll                   bool
		attachmentIDs             = internalFlag.NewMultiStringValue()
		content                   string
		contentType               internalFlag.EnumValue
		inReplyTo                 string
		language                  string
		localOnly                 bool
		mediaDescriptions         = internalFlag.NewMultiStringValue()
		mediaFiles                = internalFlag.NewMultiStringValue()
		mediaFocusValues          = internalFlag.NewMultiStringValue()
		notBoostable              bool
		notLikeable               bool
		notReplyable              bool
		pollAllowsMultipleChoices bool
		pollExpiresIn             = internalFlag.NewTimeDurationValue(24 * time.Hour)
		pollHidesVoteCounts       bool
		pollOptions               = internalFlag.NewMultiStringValue()
		sensitive                 internalFlag.BoolValue
		summary                   string
		visibility                internalFlag.EnumValue
	)

	// Parse the remaining flags.
	if err := cli.ParseStatusCreateFlags(
		&addPoll,
		&attachmentIDs,
		&content,
		&contentType,
		&inReplyTo,
		&language,
		&localOnly,
		&mediaDescriptions,
		&mediaFiles,
		&mediaFocusValues,
		&notBoostable,
		&notLikeable,
		&notReplyable,
		&pollAllowsMultipleChoices,
		&pollExpiresIn,
		&pollHidesVoteCounts,
		&pollOptions,
		&sensitive,
		&summary,
		&visibility,
		flags,
	); err != nil {
		return err //nolint:wrapcheck
	}

	// Return an error if there's no status body and no media attachments.
	if content == "" && (attachmentIDs.Length()+mediaFiles.Length() == 0) {
		return noContentOrMediaError{}
	}

	// Return an error if a poll is to be created with media attachments.
	if addPoll && (attachmentIDs.Length()+mediaFiles.Length() > 0) {
		return statusHasPollAndMediaError{}
	}

	allAttachmentIDs := attachmentIDs.Values()

	if !mediaFiles.Empty() {
		var (
			descriptionsPresent = false
			focusValuesPresent  = false
		)

		if !mediaDescriptions.Empty() {
			if !mediaDescriptions.ExpectedLength(mediaFiles.Length()) {
				return mismatchedMediaFlagsError{
					kind: "descriptions",
					want: mediaFiles.Length(),
					got:  mediaDescriptions.Length(),
				}
			}

			descriptionsPresent = true
		}

		if !mediaFocusValues.Empty() {
			if !mediaFocusValues.ExpectedLength(mediaFiles.Length()) {
				return mismatchedMediaFlagsError{
					kind: "focus values",
					want: mediaFiles.Length(),
					got:  mediaFocusValues.Length(),
				}
			}

			focusValuesPresent = true
		}

		for idx := range mediaFiles.Length() {
			var (
				mediaFile   string
				description string
				focus       string
				attachment  model.MediaAttachment
				err         error
			)

			mediaFile = mediaFiles.Values()[idx]

			if descriptionsPresent {
				description, err = utilities.ReadContents(mediaDescriptions.Values()[idx])
				if err != nil {
					return fmt.Errorf(
						"error reading the contents from %s: %w",
						mediaDescriptions.Values()[idx],
						err,
					)
				}
			}

			if focusValuesPresent {
				focus = mediaFocusValues.Values()[idx]
			}

			if err = client.Call(
				"GTSClient.CreateMediaAttachment",
				gtsclient.CreateMediaAttachmentArgs{
					Path:        mediaFile,
					Description: description,
					Focus:       focus,
				},
				&attachment,
			); err != nil {
				return fmt.Errorf("error creating the media attachment for %s: %w", mediaFile, err)
			}

			printer.PrintSuccess(
				printSettings,
				"Successfully created the media attachment with ID: "+attachment.ID,
			)

			allAttachmentIDs = append(allAttachmentIDs, attachment.ID)
		}
	}

	content, err := utilities.ReadContents(content)
	if err != nil {
		return fmt.Errorf("unable to read the content for the status: %w", err)
	}

	var preferences model.Preferences
	if err := client.Call(
		"GTSClient.GetUserPreferences",
		gtsclient.NoRPCArgs{},
		&preferences,
	); err != nil {
		printer.PrintInfo("WARNING: Unable to get your posting preferences: " + err.Error() + ".\n")
	}

	if language == "" {
		language = preferences.PostingDefaultLanguage
	}

	statusVisibility := visibility.Value()
	if statusVisibility == "" {
		statusVisibility = preferences.PostingDefaultVisibility
	}

	var statusSensitive bool
	if sensitive.IsSet() {
		statusSensitive = sensitive.Value()
	} else {
		statusSensitive = preferences.PostingDefaultSensitive
	}

	form := gtsclient.CreateStatusForm{
		Content:       content,
		ContentType:   "text/" + contentType.Value(),
		Language:      language,
		SpoilerText:   summary,
		Boostable:     !notBoostable,
		LocalOnly:     localOnly,
		InReplyTo:     inReplyTo,
		Likeable:      !notLikeable,
		Replyable:     !notReplyable,
		Sensitive:     statusSensitive,
		Visibility:    statusVisibility,
		Poll:          nil,
		AttachmentIDs: nil,
	}

	if len(allAttachmentIDs) > 0 {
		form.AttachmentIDs = allAttachmentIDs
	}

	if addPoll {
		if pollOptions.Length() == 0 {
			return noPollOptionsError{}
		}

		poll := gtsclient.CreateStatusPollForm{
			Options:    pollOptions.Values(),
			Multiple:   pollAllowsMultipleChoices,
			HideTotals: pollHidesVoteCounts,
			ExpiresIn:  int(pollExpiresIn.Value().Seconds()),
		}

		form.Poll = &poll
	}

	var status model.Status
	if err := client.Call(
		"GTSClient.CreateStatus",
		form,
		&status,
	); err != nil {
		return fmt.Errorf("error creating the status: %w", err)
	}

	printer.PrintSuccess(printSettings, "Successfully created the status with ID: "+status.ID)

	return nil
}

func statusAdd(
	client *rpc.Client,
	printSettings printer.Settings,
	relatedTarget string,
	flags []string,
) error {
	switch relatedTarget {
	case cli.TargetBookmarks:
		return statusAddToBookmarks(
			client,
			printSettings,
			flags,
		)
	default:
		return unsupportedTargetToTargetError{
			action:        cli.ActionAdd,
			focusedTarget: cli.TargetStatus,
			preposition:   cli.TargetActionPreposition(cli.TargetStatus, cli.ActionAdd),
			relatedTarget: relatedTarget,
		}
	}
}

func statusAddToBookmarks(
	client *rpc.Client,
	printSettings printer.Settings,
	flags []string,
) error {
	var statusID string

	// Parse the remaining flags.
	if err := cli.ParseStatusAddToBookmarksFlags(
		&statusID,
		flags,
	); err != nil {
		return err
	}

	if statusID == "" {
		return missingIDError{
			target: cli.TargetStatus,
			action: "add to your bookmarks",
		}
	}

	if err := client.Call(
		"GTSClient.AddStatusToBookmarks",
		statusID,
		nil,
	); err != nil {
		return fmt.Errorf("unable to add the status to your bookmarks: %w", err)
	}

	printer.PrintSuccess(printSettings, "Successfully added the status to your bookmarks.")

	return nil
}

func statusRemove(
	client *rpc.Client,
	printSettings printer.Settings,
	relatedTarget string,
	flags []string,
) error {
	switch relatedTarget {
	case cli.TargetBookmarks:
		return statusRemoveFromBookmarks(
			client,
			printSettings,
			flags,
		)
	default:
		return unsupportedTargetToTargetError{
			action:        cli.ActionRemove,
			focusedTarget: cli.TargetStatus,
			preposition:   cli.TargetActionPreposition(cli.TargetStatus, cli.ActionRemove),
			relatedTarget: relatedTarget,
		}
	}
}

func statusRemoveFromBookmarks(
	client *rpc.Client,
	printSettings printer.Settings,
	flags []string,
) error {
	var statusID string

	// Parse the remaining flags.
	if err := cli.ParseStatusRemoveFromBookmarksFlags(
		&statusID,
		flags,
	); err != nil {
		return err
	}

	if statusID == "" {
		return missingIDError{
			target: cli.TargetStatus,
			action: cli.ActionRemove,
		}
	}

	if err := client.Call("GTSClient.RemoveStatusFromBookmarks", statusID, nil); err != nil {
		return fmt.Errorf("error removing the status from your bookmarks: %w", err)
	}

	printer.PrintSuccess(printSettings, "Successfully removed the status from your bookmarks.")

	return nil
}

func statusMute(
	client *rpc.Client,
	printSettings printer.Settings,
	flags []string,
) error {
	var statusID string

	// Parse the remaining flags.
	if err := cli.ParseStatusMuteFlags(
		&statusID,
		flags,
	); err != nil {
		return err
	}

	if statusID == "" {
		return missingIDError{
			target: cli.TargetStatus,
			action: cli.ActionMute,
		}
	}

	var status model.Status
	if err := client.Call("GTSClient.GetStatus", statusID, &status); err != nil {
		return fmt.Errorf("unable to retrieve the status: %w", err)
	}

	var myAccountID string
	if err := client.Call("GTSClient.GetMyAccountID", gtsclient.NoRPCArgs{}, &myAccountID); err != nil {
		return fmt.Errorf("unable to get your account ID: %w", err)
	}

	canMute := false

	if status.Account.ID == myAccountID {
		canMute = true
	} else {
		for _, mentioned := range status.Mentions {
			if mentioned.ID == myAccountID {
				canMute = true

				break
			}
		}
	}

	if !canMute {
		return forbiddenActionOnStatusError{action: cli.ActionMute, includeNotMentioned: true}
	}

	if err := client.Call("GTSClient.MuteStatus", statusID, nil); err != nil {
		return fmt.Errorf("error muting the status: %w", err)
	}

	printer.PrintSuccess(printSettings, "Successfully muted the status.")

	return nil
}

func statusUnmute(
	client *rpc.Client,
	printSettings printer.Settings,
	flags []string,
) error {
	var statusID string

	// Parse the remaining flags.
	if err := cli.ParseStatusUnmuteFlags(
		&statusID,
		flags,
	); err != nil {
		return err
	}

	if statusID == "" {
		return missingIDError{
			target: cli.TargetStatus,
			action: cli.ActionUnmute,
		}
	}

	var status model.Status
	if err := client.Call(
		"GTSClient.GetStatus",
		statusID,
		&status,
	); err != nil {
		return fmt.Errorf("unable to retrieve the status: %w", err)
	}

	var myAccountID string
	if err := client.Call(
		"GTSClient.GetMyAccountID",
		gtsclient.NoRPCArgs{},
		&myAccountID,
	); err != nil {
		return fmt.Errorf("unable to get your account ID: %w", err)
	}

	canUnmute := false

	if status.Account.ID == myAccountID {
		canUnmute = true
	} else {
		for _, mention := range status.Mentions {
			if mention.ID == myAccountID {
				canUnmute = true

				break
			}
		}
	}

	if !canUnmute {
		return forbiddenActionOnStatusError{action: cli.ActionUnmute, includeNotMentioned: true}
	}

	if err := client.Call(
		"GTSClient.UnmuteStatus",
		statusID,
		nil,
	); err != nil {
		return fmt.Errorf("error unmuting the status: %w", err)
	}

	printer.PrintSuccess(printSettings, "Successfully unmuted the status.")

	return nil
}

func statusFavourite(
	client *rpc.Client,
	printSettings printer.Settings,
	flags []string,
) error {
	var statusID string

	// Parse the remaining flags
	if err := cli.ParseStatusFavouriteFlags(
		&statusID,
		flags,
	); err != nil {
		return err
	}

	if statusID == "" {
		return missingIDError{
			target: cli.TargetStatus,
			action: cli.ActionFavourite,
		}
	}

	if err := client.Call(
		"GTSClient.LikeStatus",
		statusID,
		nil,
	); err != nil {
		return fmt.Errorf("error favouriting the status: %w", err)
	}

	printer.PrintSuccess(
		printSettings,
		"Successfully favourited the status.",
	)

	return nil
}

func statusUnfavourite(
	client *rpc.Client,
	printSettings printer.Settings,
	flags []string,
) error {
	var statusID string

	// Parse the remaining flags
	if err := cli.ParseStatusUnfavouriteFlags(
		&statusID,
		flags,
	); err != nil {
		return err
	}

	if statusID == "" {
		return missingIDError{
			target: cli.TargetStatus,
			action: cli.ActionUnfavourite,
		}
	}

	if err := client.Call(
		"GTSClient.UnlikeStatus",
		statusID,
		nil,
	); err != nil {
		return fmt.Errorf("error unfavouriting status: %w", err)
	}

	printer.PrintSuccess(
		printSettings,
		"Successfully unfavourited the status.",
	)

	return nil
}

func statusReblog(
	client *rpc.Client,
	printSettings printer.Settings,
	flags []string,
) error {
	var statusID string

	// Parse the remaining flags
	if err := cli.ParseStatusReblogFlags(
		&statusID,
		flags,
	); err != nil {
		return err
	}

	if err := client.Call(
		"GTSClient.ReblogStatus",
		statusID,
		nil,
	); err != nil {
		return fmt.Errorf("unable to add the boost to the status: %w", err)
	}

	printer.PrintSuccess(printSettings, "Successfully reblogged the status.")

	return nil
}

func statusUnreblog(
	client *rpc.Client,
	printSettings printer.Settings,
	flags []string,
) error {
	var statusID string

	// Parse the remaining flags
	if err := cli.ParseStatusUnreblogFlags(
		&statusID,
		flags,
	); err != nil {
		return err
	}

	if err := client.Call(
		"GTSClient.UnreblogStatus",
		statusID,
		nil,
	); err != nil {
		return fmt.Errorf("unable to remove the boost from the status: %w", err)
	}

	printer.PrintSuccess(printSettings, "Successfully unreblogged the status.")

	return nil
}

func statusShow(
	client *rpc.Client,
	printSettings printer.Settings,
	browser string,
	flags []string,
) error {
	var (
		statusID          string
		showInBrowser     bool
		showWhoFavourited bool
		showWhoReblogged  bool
		status            model.Status
		rebloggedBy       model.AccountList
		favouritedBy      model.AccountList
	)

	if err := cli.ParseStatusShowFlags(
		&statusID,
		&showInBrowser,
		&showWhoFavourited,
		&showWhoReblogged,
		flags,
	); err != nil {
		return err
	}

	if statusID == "" {
		return missingIDError{
			target: cli.TargetStatus,
			action: cli.ActionShow,
		}
	}

	if err := client.Call(
		"GTSClient.GetStatus",
		statusID,
		&status,
	); err != nil {
		return fmt.Errorf("error retrieving the status: %w", err)
	}

	if showInBrowser {
		if err := utilities.OpenLink(
			browser,
			status.URL,
		); err != nil {
			return fmt.Errorf("unable to open link: %w", err)
		}

		return nil
	}

	rebloggedBy.Accounts = nil
	favouritedBy.Accounts = nil

	if showWhoReblogged {
		if err := client.Call(
			"GTSClient.GetAccountsWhoRebloggedStatus",
			statusID,
			&rebloggedBy,
		); err != nil {
			return fmt.Errorf(
				"error retrieving the list of accounts that reblogged the status: %w",
				err,
			)
		}
	}

	if showWhoFavourited {
		if err := client.Call(
			"GTSClient.GetAccountsWhoLikedStatus",
			statusID,
			&favouritedBy,
		); err != nil {
			return fmt.Errorf(
				"error retrieving the list of accounts that liked the status: %w",
				err,
			)
		}
	}

	var myAccountID string
	if err := client.Call(
		"GTSClient.GetMyAccountID",
		gtsclient.NoRPCArgs{},
		&myAccountID,
	); err != nil {
		return fmt.Errorf("unable to get your account ID: %w", err)
	}

	if err := printer.PrintStatus(
		printSettings,
		status,
		myAccountID,
		rebloggedBy,
		favouritedBy,
	); err != nil {
		return fmt.Errorf("error printing the status: %w", err)
	}

	return nil
}

func statusDelete(
	client *rpc.Client,
	printSettings printer.Settings,
	cacheRoot string,
	flags []string,
) error {
	var (
		statusID string
		saveText bool
	)

	// Parse the remaining flags.
	if err := cli.ParseStatusDeleteFlags(
		&statusID,
		&saveText,
		flags,
	); err != nil {
		return err
	}

	if statusID == "" {
		return missingIDError{
			target: cli.TargetStatus,
			action: cli.ActionDelete,
		}
	}

	var status model.Status
	if err := client.Call(
		"GTSClient.GetStatus",
		statusID,
		&status,
	); err != nil {
		return fmt.Errorf("unable to retrieve the status: %w", err)
	}

	var myAccountID string
	if err := client.Call(
		"GTSClient.GetMyAccountID",
		gtsclient.NoRPCArgs{},
		&myAccountID,
	); err != nil {
		return fmt.Errorf("unable to get your account ID: %w", err)
	}

	if status.Account.ID != myAccountID {
		return forbiddenActionOnStatusError{action: cli.ActionDelete, includeNotMentioned: false}
	}

	var text string
	if err := client.Call(
		"GTSClient.DeleteStatus",
		statusID,
		&text,
	); err != nil {
		return fmt.Errorf("error deleting the status: %w", err)
	}

	printer.PrintSuccess(printSettings, "The status was successfully deleted.")

	if saveText {
		var instance string
		if err := client.Call("GTSClient.GetInstanceURL", gtsclient.NoRPCArgs{}, &instance); err != nil {
			return fmt.Errorf("unable to get the instance URL: %w", err)
		}

		cacheDir, err := utilities.CalculateStatusesCacheDir(cacheRoot, instance)
		if err != nil {
			return fmt.Errorf("unable to get the cache directory for the status: %w", err)
		}

		if err := utilities.EnsureDirectory(cacheDir); err != nil {
			return fmt.Errorf("unable to ensure the existence of the directory %q: %w", cacheDir, err)
		}

		path := filepath.Join(cacheDir, fmt.Sprintf("deleted-status-%s.txt", statusID))

		if err := utilities.SaveTextToFile(path, text); err != nil {
			return fmt.Errorf("unable to save the text to %q: %w", path, err)
		}

		printer.PrintSuccess(printSettings, "The text was successfully saved to '"+path+"'.")
	}

	return nil
}

func statusFind(
	client *rpc.Client,
	printSettings printer.Settings,
	flags []string,
) error {
	var (
		query       string
		limit       int
		accountName string
		resolve     bool
	)

	// Parse the remaining flags
	if err := cli.ParseStatusFindFlags(
		&query,
		&limit,
		&accountName,
		&resolve,
		flags,
	); err != nil {
		return err
	}

	if query == "" {
		return missingSearchQueryError{}
	}

	var results model.StatusList

	accountID := ""
	if accountName != "" {
		if err := client.Call(
			"GTSClient.GetAccountID",
			accountName,
			&accountID,
		); err != nil {
			return fmt.Errorf("error retrieving the account ID: %w", err)
		}
	}

	if err := client.Call(
		"GTSClient.SearchStatuses",
		gtsclient.SearchStatusesArgs{
			Limit:     limit,
			Query:     query,
			AccountID: accountID,
			Resolve:   resolve,
		},
		&results,
	); err != nil {
		return fmt.Errorf("error searching for statuses: %w", err)
	}

	var myAccountID string
	if err := client.Call("GTSClient.GetMyAccountID", gtsclient.NoRPCArgs{}, &myAccountID); err != nil {
		return fmt.Errorf("error retrieving your account ID: %w", err)
	}

	if err := printer.PrintStatusList(printSettings, results, myAccountID); err != nil {
		return fmt.Errorf("error printing the search result: %w", err)
	}

	return nil
}
