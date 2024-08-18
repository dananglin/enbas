package executor

import (
	"fmt"

	"codeflow.dananglin.me.uk/apollo/enbas/internal/client"
	"codeflow.dananglin.me.uk/apollo/enbas/internal/model"
	"codeflow.dananglin.me.uk/apollo/enbas/internal/utilities"
)

func (e *EditExecutor) Execute() error {
	if e.resourceType == "" {
		return FlagNotSetError{flagText: flagType}
	}

	funcMap := map[string]func(*client.Client) error{
		resourceList:            e.editList,
		resourceMediaAttachment: e.editMediaAttachment,
	}

	doFunc, ok := funcMap[e.resourceType]
	if !ok {
		return UnsupportedTypeError{resourceType: e.resourceType}
	}

	gtsClient, err := client.NewClientFromFile(e.config.CredentialsFile)
	if err != nil {
		return fmt.Errorf("unable to create the GoToSocial client: %w", err)
	}

	return doFunc(gtsClient)
}

func (e *EditExecutor) editList(gtsClient *client.Client) error {
	if e.listID == "" {
		return MissingIDError{
			resource: resourceList,
			action: "edit",
		}
	}

	list, err := gtsClient.GetList(e.listID)
	if err != nil {
		return fmt.Errorf("unable to get the list: %w", err)
	}

	if e.listTitle != "" {
		list.Title = e.listTitle
	}

	if e.listRepliesPolicy != "" {
		parsedListRepliesPolicy, err := model.ParseListRepliesPolicy(e.listRepliesPolicy)
		if err != nil {
			return err
		}

		list.RepliesPolicy = parsedListRepliesPolicy
	}

	updatedList, err := gtsClient.UpdateList(list)
	if err != nil {
		return fmt.Errorf("unable to update the list: %w", err)
	}

	e.printer.PrintSuccess("Successfully edited the list.")
	e.printer.PrintList(updatedList)

	return nil
}

func (e *EditExecutor) editMediaAttachment(gtsClient *client.Client) error {
	expectedNumValues := 1

	if !e.attachmentIDs.ExpectedLength(expectedNumValues) {
		return UnexpectedNumValuesError{
			name:     "media attachment IDs",
			expected: expectedNumValues,
			actual:   len(e.attachmentIDs),
		}
	}

	attachment, err := gtsClient.GetMediaAttachment(e.attachmentIDs[0])
	if err != nil {
		return fmt.Errorf("unable to get the media attachment: %w", err)
	}

	description := attachment.Description
	if !e.mediaDescriptions.Empty() {
		if !e.mediaDescriptions.ExpectedLength(expectedNumValues) {
			return UnexpectedNumValuesError{
				name:     "media description",
				expected: expectedNumValues,
				actual:   len(e.mediaDescriptions),
			}
		}

		var err error

		description, err = utilities.ReadContents(e.mediaDescriptions[0])
		if err != nil {
			return fmt.Errorf(
				"unable to read the contents from %s: %w",
				e.mediaDescriptions[0],
				err,
			)
		}
	}

	focus := fmt.Sprintf("%f,%f", attachment.Meta.Focus.X, attachment.Meta.Focus.Y)
	if !e.mediaFocusValues.Empty() {
		if !e.mediaFocusValues.ExpectedLength(expectedNumValues) {
			return UnexpectedNumValuesError{
				name:     "media focus values",
				expected: expectedNumValues,
				actual:   len(e.mediaFocusValues),
			}
		}

		focus = e.mediaFocusValues[0]
	}

	if _, err = gtsClient.UpdateMediaAttachment(attachment.ID, description, focus); err != nil {
		return fmt.Errorf("unable to update the media attachment: %w", err)
	}

	e.printer.PrintSuccess("Successfully edited the media attachment.")

	return nil
}
