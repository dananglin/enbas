package executor

import (
	"flag"
	"fmt"

	"codeflow.dananglin.me.uk/apollo/enbas/internal/client"
	"codeflow.dananglin.me.uk/apollo/enbas/internal/model"
)

type CreateExecutor struct {
	*flag.FlagSet

	topLevelFlags     TopLevelFlags
	boostable         bool
	federated         bool
	likeable          bool
	replyable         bool
	sensitive         bool
	content           string
	contentType       string
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

	createExe.BoolVar(&createExe.boostable, flagEnableReposts, true, "specify if the status can be reposted/boosted by others")
	createExe.BoolVar(&createExe.federated, flagEnableFederation, true, "specify if the status can be federated beyond the local timelines")
	createExe.BoolVar(&createExe.likeable, flagEnableLikes, true, "specify if the status can be liked/favourited")
	createExe.BoolVar(&createExe.replyable, flagEnableReplies, true, "specify if the status can be replied to")
	createExe.BoolVar(&createExe.sensitive, flagSensitive, false, "specify if the status should be marked as sensitive")
	createExe.StringVar(&createExe.content, flagContent, "", "the content of the status to create")
	createExe.StringVar(&createExe.contentType, flagContentType, "plain", "the type that the contents should be parsed from (valid values are plain and markdown)")
	createExe.StringVar(&createExe.language, flagLanguage, "", "the ISO 639 language code for this status")
	createExe.StringVar(&createExe.spoilerText, flagSpoilerText, "", "the text to display as the status' warning or subject")
	createExe.StringVar(&createExe.visibility, flagVisibility, "", "the visibility of the posted status")
	createExe.StringVar(&createExe.resourceType, flagType, "", "specify the type of resource to create")
	createExe.StringVar(&createExe.listTitle, flagListTitle, "", "specify the title of the list")
	createExe.StringVar(&createExe.listRepliesPolicy, flagListRepliesPolicy, "list", "specify the policy of the replies for this list (valid values are followed, list and none)")

	createExe.Usage = commandUsageFunc(name, summary, createExe.FlagSet)

	return &createExe
}

func (c *CreateExecutor) Execute() error {
	if c.resourceType == "" {
		return FlagNotSetError{flagText: flagType}
	}

	gtsClient, err := client.NewClientFromConfig(c.topLevelFlags.ConfigDir)
	if err != nil {
		return fmt.Errorf("unable to create the GoToSocial client; %w", err)
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

	repliesPolicy, err := model.ParseListRepliesPolicy(c.listRepliesPolicy)
	if err != nil {
		return fmt.Errorf("unable to parse the list replies policy; %w", err)
	}

	list, err := gtsClient.CreateList(c.listTitle, repliesPolicy)
	if err != nil {
		return fmt.Errorf("unable to create the list; %w", err)
	}

	fmt.Println("Successfully created the following list:")
	fmt.Printf("\n%s\n", list)

	return nil
}

func (c *CreateExecutor) createStatus(gtsClient *client.Client) error {
	if c.content == "" {
		return FlagNotSetError{flagText: flagContent}
	}

	var (
		language   string
		visibility string
	)

	preferences, err := gtsClient.GetUserPreferences()
	if err != nil {
		fmt.Println("WARNING: Unable to get your posting preferences; %w", err)
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

	parsedVisibility := model.ParseStatusVisibility(visibility)
	if parsedVisibility == model.StatusVisibilityUnknown {
		return InvalidStatusVisibilityError{Visibility: visibility}
	}

	parsedContentType := model.ParseStatusContentType(c.contentType)
	if parsedContentType == model.StatusContentTypeUnknown {
		return InvalidStatusContentTypeError{ContentType: c.contentType}
	}

	form := client.CreateStatusForm{
		Content:     c.content,
		ContentType: parsedContentType,
		Language:    language,
		SpoilerText: c.spoilerText,
		Boostable:   c.boostable,
		Federated:   c.federated,
		Likeable:    c.likeable,
		Replyable:   c.replyable,
		Sensitive:   c.sensitive,
		Visibility:  parsedVisibility,
	}

	status, err := gtsClient.CreateStatus(form)
	if err != nil {
		return fmt.Errorf("unable to create the status; %w", err)
	}

	fmt.Println("Successfully created the following status:")
	fmt.Println(status)

	return nil
}