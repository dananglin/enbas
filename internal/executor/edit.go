package executor

import (
	"fmt"
	"net/rpc"

	"codeflow.dananglin.me.uk/apollo/enbas/internal/gtsclient"
	"codeflow.dananglin.me.uk/apollo/enbas/internal/model"
	"codeflow.dananglin.me.uk/apollo/enbas/internal/server"
	"codeflow.dananglin.me.uk/apollo/enbas/internal/utilities"
)

func (e *EditExecutor) Execute() error {
	if e.resourceType == "" {
		return FlagNotSetError{flagText: flagType}
	}

	funcMap := map[string]func(*rpc.Client) error{
		resourceList:            e.editList,
		resourceMediaAttachment: e.editMediaAttachment,
	}

	doFunc, ok := funcMap[e.resourceType]
	if !ok {
		return UnsupportedTypeError{resourceType: e.resourceType}
	}

	client, err := server.Connect(e.config.Server, e.configDir)
	if err != nil {
		return fmt.Errorf("error creating the client for the daemon process: %w", err)
	}
	defer client.Close()

	return doFunc(client)
}

func (e *EditExecutor) editList(client *rpc.Client) error {
	if e.listID == "" {
		return MissingIDError{
			resource: resourceList,
			action:   "edit",
		}
	}

	var listToUpdate model.List
	if err := client.Call("GTSClient.GetList", e.listID, &listToUpdate); err != nil {
		return fmt.Errorf("unable to get the list: %w", err)
	}

	if e.listTitle != "" {
		listToUpdate.Title = e.listTitle
	}

	if e.listRepliesPolicy != "" {
		parsedListRepliesPolicy, err := model.ParseListRepliesPolicy(e.listRepliesPolicy)
		if err != nil {
			return err
		}

		listToUpdate.RepliesPolicy = parsedListRepliesPolicy
	}

	var updatedList model.List
	if err := client.Call("GTSClient.UpdateList", listToUpdate, &updatedList); err != nil {
		return fmt.Errorf("error updating the list: %w", err)
	}

	acctMap, err := getAccountsFromList(client, updatedList.ID)
	if err != nil {
		return err
	}

	if len(acctMap) > 0 {
		updatedList.Accounts = acctMap
	}

	e.printer.PrintSuccess("Successfully edited the list.")
	e.printer.PrintList(updatedList)

	return nil
}

func (e *EditExecutor) editMediaAttachment(client *rpc.Client) error {
	if !e.attachmentIDs.ExpectedLength(1) {
		return UnexpectedNumValuesError{
			name:     "media attachment IDs",
			expected: 1,
			actual:   len(e.attachmentIDs),
		}
	}

	var attachment model.Attachment
	if err := client.Call("GTSClient.GetMediaAttachment", e.attachmentIDs[0], &attachment); err != nil {
		return fmt.Errorf("unable to get the media attachment: %w", err)
	}

	description := attachment.Description
	if !e.mediaDescriptions.Empty() {
		if !e.mediaDescriptions.ExpectedLength(1) {
			return UnexpectedNumValuesError{
				name:     "media description",
				expected: 1,
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
		if !e.mediaFocusValues.ExpectedLength(1) {
			return UnexpectedNumValuesError{
				name:     "media focus values",
				expected: 1,
				actual:   len(e.mediaFocusValues),
			}
		}

		focus = e.mediaFocusValues[0]
	}

	var updatedAttachment model.Attachment
	if err := client.Call(
		"GTSClient.UpdateMediaAttachment",
		gtsclient.UpdateMediaAttachmentArgs{
			MediaAttachmentID: attachment.ID,
			Description: description,
			Focus: focus,
		},
		&updatedAttachment,
	); err != nil {
		return fmt.Errorf("error updating the media attachment: %w", err)
	}

	e.printer.PrintSuccess("Successfully edited the media attachment.")

	return nil
}
