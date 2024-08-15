package executor

import (
	"errors"
	"fmt"

	"codeflow.dananglin.me.uk/apollo/enbas/internal/client"
	"codeflow.dananglin.me.uk/apollo/enbas/internal/model"
	"codeflow.dananglin.me.uk/apollo/enbas/internal/utilities"
)

func (c *CreateExecutor) Execute() error {
	if c.resourceType == "" {
		return FlagNotSetError{flagText: flagType}
	}

	gtsClient, err := client.NewClientFromFile(c.config.CredentialsFile)
	if err != nil {
		return fmt.Errorf("unable to create the GoToSocial client: %w", err)
	}

	funcMap := map[string]func(*client.Client) error{
		resourceList:            c.createList,
		resourceStatus:          c.createStatus,
		resourceMediaAttachment: c.createMediaAttachment,
	}

	doFunc, ok := funcMap[c.resourceType]
	if !ok {
		return UnsupportedTypeError{resourceType: c.resourceType}
	}

	return doFunc(gtsClient)
}

func (c *CreateExecutor) createList(gtsClient *client.Client) error {
	if c.listTitle == "" {
		return FlagNotSetError{flagText: flagListTitle}
	}

	parsedListRepliesPolicy, err := model.ParseListRepliesPolicy(c.listRepliesPolicy)
	if err != nil {
		return err
	}

	form := client.CreateListForm{
		Title:         c.listTitle,
		RepliesPolicy: parsedListRepliesPolicy,
	}

	list, err := gtsClient.CreateList(form)
	if err != nil {
		return fmt.Errorf("unable to create the list: %w", err)
	}

	c.printer.PrintSuccess("Successfully created the following list:")
	c.printer.PrintList(list)

	return nil
}

func (c *CreateExecutor) createStatus(gtsClient *client.Client) error {
	var (
		err        error
		content    string
		language   string
		visibility string
		sensitive  bool
	)

	attachmentIDs := []string(c.attachmentIDs)

	if !c.mediaFiles.Empty() {
		descriptionsExists := false
		focusValuesExists := false
		expectedLength := len(c.mediaFiles)
		mediaDescriptions := make([]string, expectedLength)

		if !c.mediaDescriptions.Empty() {
			descriptionsExists = true

			if !c.mediaDescriptions.ExpectedLength(expectedLength) {
				return errors.New("the number of media descriptions does not match the number of media files")
			}
		}

		if !c.mediaFocusValues.Empty() {
			focusValuesExists = true

			if !c.mediaFocusValues.ExpectedLength(expectedLength) {
				return errors.New("the number of media focus values does not match the number of media files")
			}
		}

		if descriptionsExists {
			for ind := 0; ind < expectedLength; ind++ {
				content, err := utilities.ReadContents(c.mediaDescriptions[ind])
				if err != nil {
					return fmt.Errorf("unable to read the contents from %s: %w", c.mediaDescriptions[ind], err)
				}

				mediaDescriptions[ind] = content
			}
		}

		for ind := 0; ind < expectedLength; ind++ {
			var (
				mediaFile   string
				description string
				focus       string
			)

			mediaFile = c.mediaFiles[ind]

			if descriptionsExists {
				description = mediaDescriptions[ind]
			}

			if focusValuesExists {
				focus = c.mediaFocusValues[ind]
			}

			attachment, err := gtsClient.CreateMediaAttachment(
				mediaFile,
				description,
				focus,
			)
			if err != nil {
				return fmt.Errorf("unable to create the media attachment for %s: %w", mediaFile, err)
			}

			attachmentIDs = append(attachmentIDs, attachment.ID)
		}
	}

	switch {
	case c.content != "":
		content = c.content
	case c.fromFile != "":
		content, err = utilities.ReadTextFile(c.fromFile)
		if err != nil {
			return fmt.Errorf("unable to get the status contents from %q: %w", c.fromFile, err)
		}
	default:
		if len(attachmentIDs) == 0 {
			// TODO: revisit this error type
			return EmptyContentError{
				ResourceType: resourceStatus,
				Hint:         "please use --" + flagContent + " or --" + flagFromFile,
			}
		}
	}

	numAttachmentIDs := len(attachmentIDs)

	if c.addPoll && numAttachmentIDs > 0 {
		return fmt.Errorf("attaching media to a poll is not allowed")
	}

	preferences, err := gtsClient.GetUserPreferences()
	if err != nil {
		fmt.Println("WARNING: Unable to get your posting preferences: %w", err)
	}

	if c.language != "" {
		language = c.language
	} else {
		language = preferences.PostingDefaultLanguage
	}

	if c.visibility != "" {
		visibility = c.visibility
	} else {
		visibility = preferences.PostingDefaultVisibility
	}

	if c.sensitive.Value != nil {
		sensitive = *c.sensitive.Value
	} else {
		sensitive = preferences.PostingDefaultSensitive
	}

	parsedVisibility, err := model.ParseStatusVisibility(visibility)
	if err != nil {
		return err
	}

	parsedContentType, err := model.ParseStatusContentType(c.contentType)
	if err != nil {
		return err
	}

	form := client.CreateStatusForm{
		Content:       content,
		ContentType:   parsedContentType,
		Language:      language,
		SpoilerText:   c.spoilerText,
		Boostable:     c.boostable,
		Federated:     c.federated,
		InReplyTo:     c.inReplyTo,
		Likeable:      c.likeable,
		Replyable:     c.replyable,
		Sensitive:     sensitive,
		Visibility:    parsedVisibility,
		Poll:          nil,
		AttachmentIDs: nil,
	}

	if numAttachmentIDs > 0 {
		form.AttachmentIDs = attachmentIDs
	}

	if c.addPoll {
		if len(c.pollOptions) == 0 {
			return NoPollOptionError{}
		}

		poll := client.CreateStatusPollForm{
			Options:    c.pollOptions,
			Multiple:   c.pollAllowsMultipleChoices,
			HideTotals: c.pollHidesVoteCounts,
			ExpiresIn:  int(c.pollExpiresIn.Duration.Seconds()),
		}

		form.Poll = &poll
	}

	status, err := gtsClient.CreateStatus(form)
	if err != nil {
		return fmt.Errorf("unable to create the status: %w", err)
	}

	c.printer.PrintSuccess("Successfully created the status with ID: " + status.ID)

	return nil
}

func (c *CreateExecutor) createMediaAttachment(gtsClient *client.Client) error {
	expectedNumValues := 1
	if !c.mediaFiles.ExpectedLength(expectedNumValues) {
		return fmt.Errorf(
			"received an unexpected number of media files: want %d",
			expectedNumValues,
		)
	}

	description := ""
	if !c.mediaDescriptions.Empty() {
		if !c.mediaDescriptions.ExpectedLength(expectedNumValues) {
			return fmt.Errorf(
				"received an unexpected number of media descriptions: want %d",
				expectedNumValues,
			)
		}

		var err error
		description, err = utilities.ReadContents(c.mediaDescriptions[0])
		if err != nil {
			return fmt.Errorf(
				"unable to read the contents from %s: %w",
				c.mediaDescriptions[0],
			)
		}
	}

	focus := ""
	if !c.mediaFocusValues.Empty() {
		if !c.mediaFocusValues.ExpectedLength(expectedNumValues) {
			return fmt.Errorf(
				"received an unexpected number of media focus values: want %d",
				expectedNumValues,
			)
		}
		focus = c.mediaFocusValues[0]
	}

	attachment, err := gtsClient.CreateMediaAttachment(
		c.mediaFiles[0],
		description,
		focus,
	)
	if err != nil {
		return fmt.Errorf("unable to create the media attachment: %w", err)
	}

	c.printer.PrintSuccess("Successfully created the media attachment with ID: " + attachment.ID)

	return nil
}
