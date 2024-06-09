// SPDX-FileCopyrightText: 2024 Dan Anglin <d.n.i.anglin@gmail.com>
//
// SPDX-License-Identifier: GPL-3.0-or-later

package executor

import (
	"flag"
	"fmt"
	"strconv"

	"codeflow.dananglin.me.uk/apollo/enbas/internal/client"
	"codeflow.dananglin.me.uk/apollo/enbas/internal/model"
	"codeflow.dananglin.me.uk/apollo/enbas/internal/utilities"
)

type CreateExecutor struct {
	*flag.FlagSet

	topLevelFlags     TopLevelFlags
	boostable         bool
	federated         bool
	likeable          bool
	replyable         bool
	sensitive         *bool
	content           string
	contentType       string
	fromFile          string
	language          string
	spoilerText       string
	resourceType      string
	listTitle         string
	listRepliesPolicy string
	visibility        string
}

func NewCreateExecutor(tlf TopLevelFlags, name, summary string) *CreateExecutor {
	createExe := CreateExecutor{
		FlagSet: flag.NewFlagSet(name, flag.ExitOnError),

		topLevelFlags: tlf,
	}

	createExe.BoolVar(&createExe.boostable, flagEnableReposts, true, "Specify if the status can be reposted/boosted by others")
	createExe.BoolVar(&createExe.federated, flagEnableFederation, true, "Specify if the status can be federated beyond the local timelines")
	createExe.BoolVar(&createExe.likeable, flagEnableLikes, true, "Specify if the status can be liked/favourited")
	createExe.BoolVar(&createExe.replyable, flagEnableReplies, true, "Specify if the status can be replied to")
	createExe.StringVar(&createExe.content, flagContent, "", "The content of the status to create")
	createExe.StringVar(&createExe.contentType, flagContentType, "plain", "The type that the contents should be parsed from (valid values are plain and markdown)")
	createExe.StringVar(&createExe.fromFile, flagFromFile, "", "The file path where to read the contents from")
	createExe.StringVar(&createExe.language, flagLanguage, "", "The ISO 639 language code for this status")
	createExe.StringVar(&createExe.spoilerText, flagSpoilerText, "", "The text to display as the status' warning or subject")
	createExe.StringVar(&createExe.visibility, flagVisibility, "", "The visibility of the posted status")
	createExe.StringVar(&createExe.resourceType, flagType, "", "Specify the type of resource to create")
	createExe.StringVar(&createExe.listTitle, flagListTitle, "", "Specify the title of the list")
	createExe.StringVar(&createExe.listRepliesPolicy, flagListRepliesPolicy, "list", "Specify the policy of the replies for this list (valid values are followed, list and none)")

	createExe.BoolFunc(flagSensitive, "Specify if the status should be marked as sensitive", func(value string) error {
		boolVal, err := strconv.ParseBool(value)
		if err != nil {
			return fmt.Errorf("unable to parse %q as a boolean value: %w", value, err)
		}

		createExe.sensitive = new(bool)
		*createExe.sensitive = boolVal

		return nil
	})

	createExe.Usage = commandUsageFunc(name, summary, createExe.FlagSet)

	return &createExe
}

func (c *CreateExecutor) Execute() error {
	if c.resourceType == "" {
		return FlagNotSetError{flagText: flagType}
	}

	gtsClient, err := client.NewClientFromConfig(c.topLevelFlags.ConfigDir)
	if err != nil {
		return fmt.Errorf("unable to create the GoToSocial client: %w", err)
	}

	funcMap := map[string]func(*client.Client) error{
		resourceList:   c.createList,
		resourceStatus: c.createStatus,
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

	fmt.Println("Successfully created the following list:")
	utilities.Display(list, *c.topLevelFlags.NoColor, c.topLevelFlags.Pager)

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

	switch {
	case c.content != "":
		content = c.content
	case c.fromFile != "":
		content, err = utilities.ReadFile(c.fromFile)
		if err != nil {
			return fmt.Errorf("unable to get the status contents from %q: %w", c.fromFile, err)
		}
	default:
		return EmptyContentError{
			ResourceType: resourceStatus,
			Hint:         "please use --" + flagContent + " or --" + flagFromFile,
		}
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

	if c.sensitive != nil {
		sensitive = *c.sensitive
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
		Content:     content,
		ContentType: parsedContentType,
		Language:    language,
		SpoilerText: c.spoilerText,
		Boostable:   c.boostable,
		Federated:   c.federated,
		Likeable:    c.likeable,
		Replyable:   c.replyable,
		Sensitive:   sensitive,
		Visibility:  parsedVisibility,
	}

	status, err := gtsClient.CreateStatus(form)
	if err != nil {
		return fmt.Errorf("unable to create the status: %w", err)
	}

	fmt.Println("Successfully created the following status:")
	utilities.Display(status, *c.topLevelFlags.NoColor, c.topLevelFlags.Pager)

	return nil
}
