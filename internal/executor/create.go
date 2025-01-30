package executor

import (
	"fmt"
	"net/rpc"

	"codeflow.dananglin.me.uk/apollo/enbas/internal/gtsclient"
	"codeflow.dananglin.me.uk/apollo/enbas/internal/model"
	"codeflow.dananglin.me.uk/apollo/enbas/internal/server"
	"codeflow.dananglin.me.uk/apollo/enbas/internal/utilities"
)

func (c *CreateExecutor) Execute() error {
	if c.resourceType == "" {
		return FlagNotSetError{flagText: flagType}
	}

	funcMap := map[string]func(*rpc.Client) error{
		resourceList:            c.createList,
		resourceStatus:          c.createStatus,
		resourceMediaAttachment: c.createMediaAttachment,
	}

	doFunc, ok := funcMap[c.resourceType]
	if !ok {
		return UnsupportedTypeError{resourceType: c.resourceType}
	}

	client, err := server.Connect(c.config.Server, c.configDir)
	if err != nil {
		return fmt.Errorf("error creating the client for the daemon process: %w", err)
	}
	defer client.Close()

	return doFunc(client)
}

func (c *CreateExecutor) createList(client *rpc.Client) error {
	if c.listTitle == "" {
		return Error{"please provide the title of the list that you want to create"}
	}

	parsedListRepliesPolicy, err := model.ParseListRepliesPolicy(c.listRepliesPolicy)
	if err != nil {
		return err //nolint:wrapcheck
	}

	var list model.List
	if err := client.Call(
		"GTSClient.CreateList",
		gtsclient.CreateListArgs{
			Title:         c.listTitle,
			RepliesPolicy: parsedListRepliesPolicy,
			Exclusive:     c.listExclusive,
		},
		&list,
	); err != nil {
		return fmt.Errorf("unable to create the list: %w", err)
	}

	c.printer.PrintSuccess("Successfully created the following list:")
	c.printer.PrintList(list)

	return nil
}

func (c *CreateExecutor) createStatus(client *rpc.Client) error {
	var (
		err        error
		language   string
		visibility string
		sensitive  bool
	)

	attachmentIDs := []string(c.attachmentIDs)

	if !c.mediaFiles.Empty() {
		descriptionsExists := false
		focusValuesExists := false
		numMediaFiles := len(c.mediaFiles)
		mediaDescriptions := make([]string, numMediaFiles)

		if !c.mediaDescriptions.Empty() {
			descriptionsExists = true

			if !c.mediaDescriptions.ExpectedLength(numMediaFiles) {
				return MismatchedNumMediaValuesError{
					valueType:     "media descriptions",
					numValues:     len(c.mediaDescriptions),
					numMediaFiles: numMediaFiles,
				}
			}
		}

		if !c.mediaFocusValues.Empty() {
			focusValuesExists = true

			if !c.mediaFocusValues.ExpectedLength(numMediaFiles) {
				return MismatchedNumMediaValuesError{
					valueType:     "media focus values",
					numValues:     len(c.mediaFocusValues),
					numMediaFiles: numMediaFiles,
				}
			}
		}

		if descriptionsExists {
			for ind := range numMediaFiles {
				mediaDesc, err := utilities.ReadContents(c.mediaDescriptions[ind])
				if err != nil {
					return fmt.Errorf(
						"unable to read the contents from %s: %w",
						c.mediaDescriptions[ind],
						err,
					)
				}

				mediaDescriptions[ind] = mediaDesc
			}
		}

		for ind := range numMediaFiles {
			var (
				mediaFile   string
				description string
				focus       string
				attachment  model.MediaAttachment
			)

			mediaFile = c.mediaFiles[ind]

			if descriptionsExists {
				description = mediaDescriptions[ind]
			}

			if focusValuesExists {
				focus = c.mediaFocusValues[ind]
			}

			if err := client.Call(
				"GTSClient.CreateMediaAttachment",
				gtsclient.CreateMediaAttachmentArgs{
					Path:        mediaFile,
					Description: description,
					Focus:       focus,
				},
				&attachment,
			); err != nil {
				return fmt.Errorf("unable to create the media attachment for %s: %w", mediaFile, err)
			}

			attachmentIDs = append(attachmentIDs, attachment.ID)
		}
	}

	if c.content == "" && len(attachmentIDs) == 0 {
		return Error{"please add content to the status that you want to create"}
	}

	content, err := utilities.ReadContents(c.content)
	if err != nil {
		return fmt.Errorf("unable to read the contents from %s: %w", c.content, err)
	}

	numAttachmentIDs := len(attachmentIDs)

	if c.addPoll && numAttachmentIDs > 0 {
		return Error{"attaching media to a poll is not allowed"}
	}

	var preferences model.Preferences
	if err := client.Call("GTSClient.GetUserPreferences", gtsclient.NoRPCArgs{}, &preferences); err != nil {
		c.printer.PrintInfo("WARNING: Unable to get your posting preferences: " + err.Error() + ".\n")
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
		return err //nolint:wrapcheck
	}

	parsedContentType, err := model.ParseStatusContentType(c.contentType)
	if err != nil {
		return err //nolint:wrapcheck
	}

	form := gtsclient.CreateStatusForm{
		Content:       content,
		ContentType:   parsedContentType,
		Language:      language,
		SpoilerText:   c.summary,
		Boostable:     c.boostable,
		LocalOnly:     c.localOnly,
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
			return Error{"no options were provided for this poll"}
		}

		poll := gtsclient.CreateStatusPollForm{
			Options:    c.pollOptions,
			Multiple:   c.pollAllowsMultipleChoices,
			HideTotals: c.pollHidesVoteCounts,
			ExpiresIn:  int(c.pollExpiresIn.Duration.Seconds()),
		}

		form.Poll = &poll
	}

	var status model.Status
	if err := client.Call("GTSClient.CreateStatus", form, &status); err != nil {
		return fmt.Errorf("error creating the status: %w", err)
	}

	c.printer.PrintSuccess("Successfully created the status with ID: " + status.ID)

	return nil
}

func (c *CreateExecutor) createMediaAttachment(client *rpc.Client) error {
	expectedNumValues := 1

	if !c.mediaFiles.ExpectedLength(expectedNumValues) {
		return UnexpectedNumValuesError{
			name:     "media files",
			expected: expectedNumValues,
			actual:   len(c.mediaFiles),
		}
	}

	description := ""
	if !c.mediaDescriptions.Empty() {
		if !c.mediaDescriptions.ExpectedLength(expectedNumValues) {
			return UnexpectedNumValuesError{
				name:     "media descriptions",
				expected: expectedNumValues,
				actual:   len(c.mediaDescriptions),
			}
		}

		var err error

		description, err = utilities.ReadContents(c.mediaDescriptions[0])
		if err != nil {
			return fmt.Errorf(
				"unable to read the contents from %s: %w",
				c.mediaDescriptions[0],
				err,
			)
		}
	}

	focus := ""
	if !c.mediaFocusValues.Empty() {
		if !c.mediaFocusValues.ExpectedLength(expectedNumValues) {
			return UnexpectedNumValuesError{
				name:     "media focus values",
				expected: expectedNumValues,
				actual:   len(c.mediaFocusValues),
			}
		}

		focus = c.mediaFocusValues[0]
	}

	var attachment model.MediaAttachment
	if err := client.Call(
		"GTSClient.CreateMediaAttachment",
		gtsclient.CreateMediaAttachmentArgs{
			Path:        c.mediaFiles[0],
			Description: description,
			Focus:       focus,
		},
		&attachment,
	); err != nil {
		return fmt.Errorf("unable to create the media attachment: %w", err)
	}

	c.printer.PrintSuccess("Successfully created the media attachment with ID: " + attachment.ID)

	return nil
}
