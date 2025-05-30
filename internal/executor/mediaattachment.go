package executor

import (
	"fmt"
	"net/rpc"

	"codeflow.dananglin.me.uk/apollo/enbas/internal/cli"
	"codeflow.dananglin.me.uk/apollo/enbas/internal/command"
	"codeflow.dananglin.me.uk/apollo/enbas/internal/config"
	"codeflow.dananglin.me.uk/apollo/enbas/internal/gtsclient"
	"codeflow.dananglin.me.uk/apollo/enbas/internal/model"
	"codeflow.dananglin.me.uk/apollo/enbas/internal/printer"
	"codeflow.dananglin.me.uk/apollo/enbas/internal/server"
	"codeflow.dananglin.me.uk/apollo/enbas/internal/utilities"
)

// mediaAttachmentFunc is the function for the 'media-attachment' target for
// interacting with a single media attachment.
func mediaAttachmentFunc(
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
		return mediaAttachmentCreate(session.Client(), printSettings, cmd.FocusedTargetFlags)
	case cli.ActionEdit:
		return mediaAttachmentEdit(session.Client(), printSettings, cmd.FocusedTargetFlags)
	case cli.ActionShow:
		return mediaAttachmentShow(session.Client(), printSettings, cmd.FocusedTargetFlags)
	default:
		return unsupportedActionError{action: cmd.Action, target: cli.TargetMediaAttachment}
	}
}

func mediaAttachmentCreate(
	client *rpc.Client,
	printSettings printer.Settings,
	flags []string,
) error {
	var (
		description string
		file        string
		focus       string
		err         error
	)

	// Parse the remaining flags.
	if err := cli.ParseMediaAttachmentCreateFlags(
		&description,
		&file,
		&focus,
		flags,
	); err != nil {
		return err
	}

	if file == "" {
		return missingMediaFileError{}
	}

	parsedDescription := ""

	if description != "" {
		parsedDescription, err = utilities.ReadContents(description)
		if err != nil {
			return fmt.Errorf(
				"error reading the contents from %s: %w",
				description,
				err,
			)
		}
	}

	var attachment model.MediaAttachment
	if err := client.Call(
		"GTSClient.CreateMediaAttachment",
		gtsclient.CreateMediaAttachmentArgs{
			Path:        file,
			Description: parsedDescription,
			Focus:       focus,
		},
		&attachment,
	); err != nil {
		return fmt.Errorf("error creating the media attachment: %w", err)
	}

	printer.PrintSuccess(printSettings, "Successfully created the media attachment with ID: "+attachment.ID)

	return nil
}

func mediaAttachmentEdit(
	client *rpc.Client,
	printSettings printer.Settings,
	flags []string,
) error {
	var (
		attachmentID string
		description  string
		focus        string
		err          error
	)

	// Parse the remaining flags.
	if err := cli.ParseMediaAttachmentEditFlags(
		&attachmentID,
		&description,
		&focus,
		flags,
	); err != nil {
		return err
	}

	if attachmentID == "" {
		return missingIDError{
			target: cli.TargetMediaAttachment,
			action: cli.ActionEdit,
		}
	}

	var attachment model.MediaAttachment
	if err = client.Call(
		"GTSClient.GetMediaAttachment",
		attachmentID,
		&attachment,
	); err != nil {
		return fmt.Errorf("error retrieving the media attachment: %w", err)
	}

	updatedDescription := attachment.Description
	if description != "" {
		updatedDescription, err = utilities.ReadContents(description)
		if err != nil {
			return fmt.Errorf(
				"unable to read the contents from %s: %w",
				description,
				err,
			)
		}
	}

	updatedFocus := fmt.Sprintf("%f,%f", attachment.Meta.Focus.X, attachment.Meta.Focus.Y)
	if focus != "" {
		updatedFocus = focus
	}

	var updatedAttachment model.MediaAttachment
	if err = client.Call(
		"GTSClient.UpdateMediaAttachment",
		gtsclient.UpdateMediaAttachmentArgs{
			MediaAttachmentID: attachment.ID,
			Description:       updatedDescription,
			Focus:             updatedFocus,
		},
		&updatedAttachment,
	); err != nil {
		return fmt.Errorf("error updating the media attachment: %w", err)
	}

	printer.PrintSuccess(printSettings, "Successfully edited the media attachment.")

	return nil
}

func mediaAttachmentShow(
	client *rpc.Client,
	printSettings printer.Settings,
	flags []string,
) error {
	var attachmentID string

	// Parse the remaining flags.
	if err := cli.ParseMediaAttachmentShowFlags(
		&attachmentID,
		flags,
	); err != nil {
		return err
	}

	if attachmentID == "" {
		return missingIDError{
			target: cli.TargetMediaAttachment,
			action: cli.ActionShow,
		}
	}

	var attachment model.MediaAttachment
	if err := client.Call(
		"GTSClient.GetMediaAttachment",
		attachmentID,
		&attachment,
	); err != nil {
		return fmt.Errorf("unable to retrieve the media attachment: %w", err)
	}

	if err := printer.PrintMediaAttachment(printSettings, attachment); err != nil {
		return fmt.Errorf("error printing the media attachment details: %w", err)
	}

	return nil
}
