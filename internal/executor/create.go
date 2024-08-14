package executor

import (
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
		Content:     content,
		ContentType: parsedContentType,
		Language:    language,
		SpoilerText: c.spoilerText,
		Boostable:   c.boostable,
		Federated:   c.federated,
		InReplyTo:   c.inReplyTo,
		Likeable:    c.likeable,
		Replyable:   c.replyable,
		Sensitive:   sensitive,
		Visibility:  parsedVisibility,
		Poll:        nil,
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
