/*
   This file is generated by the enbas-codegen
   DO NOT EDIT.
*/

package executor

import (
	"flag"

	"codeflow.dananglin.me.uk/apollo/enbas/internal/config"
	internalFlag "codeflow.dananglin.me.uk/apollo/enbas/internal/flag"
	"codeflow.dananglin.me.uk/apollo/enbas/internal/printer"
	"codeflow.dananglin.me.uk/apollo/enbas/internal/usage"
)

type Executor interface {
	Name() string
	Parse(args []string) error
	Execute() error
}

// AcceptExecutor is the executor for the accept command.
type AcceptExecutor struct {
	*flag.FlagSet
	printer      *printer.Printer
	config       *config.Config
	accountName  internalFlag.StringSliceValue
	resourceType string
}

func NewAcceptExecutor(
	printer *printer.Printer,
	config *config.Config,
) *AcceptExecutor {
	exe := AcceptExecutor{
		FlagSet:     flag.NewFlagSet("accept", flag.ExitOnError),
		printer:     printer,
		config:      config,
		accountName: internalFlag.NewStringSliceValue(),
	}

	exe.Usage = usage.ExecutorUsageFunc("accept", "Accepts a request (e.g. a follow request)", exe.FlagSet)

	exe.Var(&exe.accountName, "account-name", "The name of the account")
	exe.StringVar(&exe.resourceType, "type", "", "The type of resource you want to action on (e.g. account, status)")

	return &exe
}

// AddExecutor is the executor for the add command.
type AddExecutor struct {
	*flag.FlagSet
	printer        *printer.Printer
	config         *config.Config
	accountNames   internalFlag.StringSliceValue
	content        string
	listID         string
	pollID         string
	statusID       string
	toResourceType string
	resourceType   string
	votes          internalFlag.IntSliceValue
}

func NewAddExecutor(
	printer *printer.Printer,
	config *config.Config,
) *AddExecutor {
	exe := AddExecutor{
		FlagSet:      flag.NewFlagSet("add", flag.ExitOnError),
		printer:      printer,
		config:       config,
		accountNames: internalFlag.NewStringSliceValue(),
		votes:        internalFlag.NewIntSliceValue(),
	}

	exe.Usage = usage.ExecutorUsageFunc("add", "Adds a resource to another resource", exe.FlagSet)

	exe.Var(&exe.accountNames, "account-name", "The name of the account")
	exe.StringVar(&exe.content, "content", "", "The content of the created resource")
	exe.StringVar(&exe.listID, "list-id", "", "The ID of the list in question")
	exe.StringVar(&exe.pollID, "poll-id", "", "The ID of the poll")
	exe.StringVar(&exe.statusID, "status-id", "", "The ID of the status")
	exe.StringVar(&exe.toResourceType, "to", "", "TBC")
	exe.StringVar(&exe.resourceType, "type", "", "The type of resource you want to action on (e.g. account, status)")
	exe.Var(&exe.votes, "vote", "Add a vote to an option in a poll")

	return &exe
}

// BlockExecutor is the executor for the block command.
type BlockExecutor struct {
	*flag.FlagSet
	printer      *printer.Printer
	config       *config.Config
	accountName  internalFlag.StringSliceValue
	resourceType string
}

func NewBlockExecutor(
	printer *printer.Printer,
	config *config.Config,
) *BlockExecutor {
	exe := BlockExecutor{
		FlagSet:     flag.NewFlagSet("block", flag.ExitOnError),
		printer:     printer,
		config:      config,
		accountName: internalFlag.NewStringSliceValue(),
	}

	exe.Usage = usage.ExecutorUsageFunc("block", "Blocks a resource (e.g. an account)", exe.FlagSet)

	exe.Var(&exe.accountName, "account-name", "The name of the account")
	exe.StringVar(&exe.resourceType, "type", "", "The type of resource you want to action on (e.g. account, status)")

	return &exe
}

// CreateExecutor is the executor for the create command.
type CreateExecutor struct {
	*flag.FlagSet
	printer                   *printer.Printer
	config                    *config.Config
	addPoll                   bool
	content                   string
	contentType               string
	federated                 bool
	likeable                  bool
	replyable                 bool
	boostable                 bool
	fromFile                  string
	inReplyTo                 string
	language                  string
	listRepliesPolicy         string
	listTitle                 string
	pollAllowsMultipleChoices bool
	pollExpiresIn             internalFlag.TimeDurationValue
	pollHidesVoteCounts       bool
	pollOptions               internalFlag.StringSliceValue
	sensitive                 internalFlag.BoolPtrValue
	spoilerText               string
	resourceType              string
	visibility                string
}

func NewCreateExecutor(
	printer *printer.Printer,
	config *config.Config,
) *CreateExecutor {
	exe := CreateExecutor{
		FlagSet:       flag.NewFlagSet("create", flag.ExitOnError),
		printer:       printer,
		config:        config,
		pollExpiresIn: internalFlag.NewTimeDurationValue(),
		pollOptions:   internalFlag.NewStringSliceValue(),
		sensitive:     internalFlag.NewBoolPtrValue(),
	}

	exe.Usage = usage.ExecutorUsageFunc("create", "Creates a specific resource", exe.FlagSet)

	exe.BoolVar(&exe.addPoll, "add-poll", false, "Set to true to add a poll when creating a status")
	exe.StringVar(&exe.content, "content", "", "The content of the created resource")
	exe.StringVar(&exe.contentType, "content-type", "plain", "The type that the contents should be parsed from (valid values are plain and markdown)")
	exe.BoolVar(&exe.federated, "enable-federation", true, "Set to true to federate the status beyond the local timelines")
	exe.BoolVar(&exe.likeable, "enable-likes", true, "Set to true to allow the status to be liked (favourited)")
	exe.BoolVar(&exe.replyable, "enable-replies", true, "Set to true to allow viewers to reply to the status")
	exe.BoolVar(&exe.boostable, "enable-reposts", true, "Set to true to allow the status to be reposted (boosted) by others")
	exe.StringVar(&exe.fromFile, "from-file", "", "The file path where to read the contents from")
	exe.StringVar(&exe.inReplyTo, "in-reply-to", "", "The ID of the status that you want to reply to")
	exe.StringVar(&exe.language, "language", "", "The ISO 639 language code for this status")
	exe.StringVar(&exe.listRepliesPolicy, "list-replies-policy", "list", "The replies policy of the list")
	exe.StringVar(&exe.listTitle, "list-title", "", "The title of the list")
	exe.BoolVar(&exe.pollAllowsMultipleChoices, "poll-allows-multiple-choices", false, "Set to true to allow viewers to make multiple choices in the poll")
	exe.Var(&exe.pollExpiresIn, "poll-expires-in", "The duration in which the poll is open for")
	exe.BoolVar(&exe.pollHidesVoteCounts, "poll-hides-vote-counts", false, "Set to true to hide the vote count until the poll is closed")
	exe.Var(&exe.pollOptions, "poll-option", "A poll option. Use this multiple times to set multiple options")
	exe.Var(&exe.sensitive, "sensitive", "Set to true if the status should be marked as sensitive")
	exe.StringVar(&exe.spoilerText, "spoiler-text", "", "The text to display as the status' warning or subject")
	exe.StringVar(&exe.resourceType, "type", "", "The type of resource you want to action on (e.g. account, status)")
	exe.StringVar(&exe.visibility, "visibility", "", "The visibility of the posted status")

	return &exe
}

// DeleteExecutor is the executor for the delete command.
type DeleteExecutor struct {
	*flag.FlagSet
	printer      *printer.Printer
	config       *config.Config
	listID       string
	resourceType string
}

func NewDeleteExecutor(
	printer *printer.Printer,
	config *config.Config,
) *DeleteExecutor {
	exe := DeleteExecutor{
		FlagSet: flag.NewFlagSet("delete", flag.ExitOnError),
		printer: printer,
		config:  config,
	}

	exe.Usage = usage.ExecutorUsageFunc("delete", "Deletes a specific resource", exe.FlagSet)

	exe.StringVar(&exe.listID, "list-id", "", "The ID of the list in question")
	exe.StringVar(&exe.resourceType, "type", "", "The type of resource you want to action on (e.g. account, status)")

	return &exe
}

// EditExecutor is the executor for the edit command.
type EditExecutor struct {
	*flag.FlagSet
	printer           *printer.Printer
	config            *config.Config
	listID            string
	listTitle         string
	listRepliesPolicy string
	resourceType      string
}

func NewEditExecutor(
	printer *printer.Printer,
	config *config.Config,
) *EditExecutor {
	exe := EditExecutor{
		FlagSet: flag.NewFlagSet("edit", flag.ExitOnError),
		printer: printer,
		config:  config,
	}

	exe.Usage = usage.ExecutorUsageFunc("edit", "Edit a specific resource", exe.FlagSet)

	exe.StringVar(&exe.listID, "list-id", "", "The ID of the list in question")
	exe.StringVar(&exe.listTitle, "list-title", "", "The title of the list")
	exe.StringVar(&exe.listRepliesPolicy, "list-replies-policy", "", "The replies policy of the list")
	exe.StringVar(&exe.resourceType, "type", "", "The type of resource you want to action on (e.g. account, status)")

	return &exe
}

// FollowExecutor is the executor for the follow command.
type FollowExecutor struct {
	*flag.FlagSet
	printer      *printer.Printer
	config       *config.Config
	accountName  internalFlag.StringSliceValue
	notify       bool
	showReposts  bool
	resourceType string
}

func NewFollowExecutor(
	printer *printer.Printer,
	config *config.Config,
) *FollowExecutor {
	exe := FollowExecutor{
		FlagSet:     flag.NewFlagSet("follow", flag.ExitOnError),
		printer:     printer,
		config:      config,
		accountName: internalFlag.NewStringSliceValue(),
	}

	exe.Usage = usage.ExecutorUsageFunc("follow", "Follow a resource (e.g. an account)", exe.FlagSet)

	exe.Var(&exe.accountName, "account-name", "The name of the account")
	exe.BoolVar(&exe.notify, "notify", false, "Get notifications from statuses from the account you want to follow")
	exe.BoolVar(&exe.showReposts, "show-reposts", true, "Show reposts from the account you want to follow")
	exe.StringVar(&exe.resourceType, "type", "", "The type of resource you want to action on (e.g. account, status)")

	return &exe
}

// InitExecutor is the executor for the init command.
type InitExecutor struct {
	*flag.FlagSet
	printer   *printer.Printer
	configDir string
}

func NewInitExecutor(
	printer *printer.Printer,
	configDir string,
) *InitExecutor {
	exe := InitExecutor{
		FlagSet:   flag.NewFlagSet("init", flag.ExitOnError),
		printer:   printer,
		configDir: configDir,
	}

	exe.Usage = usage.ExecutorUsageFunc("init", "Creates a new configuration file in the specified configuration directory", exe.FlagSet)

	return &exe
}

// LoginExecutor is the executor for the login command.
type LoginExecutor struct {
	*flag.FlagSet
	printer  *printer.Printer
	config   *config.Config
	instance string
}

func NewLoginExecutor(
	printer *printer.Printer,
	config *config.Config,
) *LoginExecutor {
	exe := LoginExecutor{
		FlagSet: flag.NewFlagSet("login", flag.ExitOnError),
		printer: printer,
		config:  config,
	}

	exe.Usage = usage.ExecutorUsageFunc("login", "Logs into an account on GoToSocial", exe.FlagSet)

	exe.StringVar(&exe.instance, "instance", "", "The instance that you want to log into")

	return &exe
}

// MuteExecutor is the executor for the mute command.
type MuteExecutor struct {
	*flag.FlagSet
	printer           *printer.Printer
	config            *config.Config
	accountName       internalFlag.StringSliceValue
	muteDuration      internalFlag.TimeDurationValue
	muteNotifications bool
	resourceType      string
}

func NewMuteExecutor(
	printer *printer.Printer,
	config *config.Config,
) *MuteExecutor {
	exe := MuteExecutor{
		FlagSet:      flag.NewFlagSet("mute", flag.ExitOnError),
		printer:      printer,
		config:       config,
		accountName:  internalFlag.NewStringSliceValue(),
		muteDuration: internalFlag.NewTimeDurationValue(),
	}

	exe.Usage = usage.ExecutorUsageFunc("mute", "Mutes a specific resource (e.g. an account)", exe.FlagSet)

	exe.Var(&exe.accountName, "account-name", "The name of the account")
	exe.Var(&exe.muteDuration, "mute-duration", "Specify how long the mute should last for. To mute indefinitely, set this to 0s")
	exe.BoolVar(&exe.muteNotifications, "mute-notifications", false, "Set to true to mute notifications as well as posts")
	exe.StringVar(&exe.resourceType, "type", "", "The type of resource you want to action on (e.g. account, status)")

	return &exe
}

// RejectExecutor is the executor for the reject command.
type RejectExecutor struct {
	*flag.FlagSet
	printer      *printer.Printer
	config       *config.Config
	accountName  internalFlag.StringSliceValue
	resourceType string
}

func NewRejectExecutor(
	printer *printer.Printer,
	config *config.Config,
) *RejectExecutor {
	exe := RejectExecutor{
		FlagSet:     flag.NewFlagSet("reject", flag.ExitOnError),
		printer:     printer,
		config:      config,
		accountName: internalFlag.NewStringSliceValue(),
	}

	exe.Usage = usage.ExecutorUsageFunc("reject", "Rejects a request (e.g. a follow request)", exe.FlagSet)

	exe.Var(&exe.accountName, "account-name", "The name of the account")
	exe.StringVar(&exe.resourceType, "type", "", "The type of resource you want to action on (e.g. account, status)")

	return &exe
}

// RemoveExecutor is the executor for the remove command.
type RemoveExecutor struct {
	*flag.FlagSet
	printer          *printer.Printer
	config           *config.Config
	accountNames     internalFlag.StringSliceValue
	fromResourceType string
	listID           string
	statusID         string
	resourceType     string
}

func NewRemoveExecutor(
	printer *printer.Printer,
	config *config.Config,
) *RemoveExecutor {
	exe := RemoveExecutor{
		FlagSet:      flag.NewFlagSet("remove", flag.ExitOnError),
		printer:      printer,
		config:       config,
		accountNames: internalFlag.NewStringSliceValue(),
	}

	exe.Usage = usage.ExecutorUsageFunc("remove", "Removes a resource from another resource", exe.FlagSet)

	exe.Var(&exe.accountNames, "account-name", "The name of the account")
	exe.StringVar(&exe.fromResourceType, "from", "", "Specify the resource type to action the target resource from")
	exe.StringVar(&exe.listID, "list-id", "", "The ID of the list in question")
	exe.StringVar(&exe.statusID, "status-id", "", "The ID of the status")
	exe.StringVar(&exe.resourceType, "type", "", "The type of resource you want to action on (e.g. account, status)")

	return &exe
}

// ShowExecutor is the executor for the show command.
type ShowExecutor struct {
	*flag.FlagSet
	printer                 *printer.Printer
	config                  *config.Config
	accountName             internalFlag.StringSliceValue
	getAllImages            bool
	getAllVideos            bool
	attachmentIDs           internalFlag.StringSliceValue
	showInBrowser           bool
	excludeBoosts           bool
	excludeReplies          bool
	fromResourceType        string
	limit                   int
	listID                  string
	myAccount               bool
	onlyMedia               bool
	onlyPinned              bool
	onlyPublic              bool
	pollID                  string
	showUserPreferences     bool
	showStatuses            bool
	skipAccountRelationship bool
	statusID                string
	timelineCategory        string
	tag                     string
	resourceType            string
}

func NewShowExecutor(
	printer *printer.Printer,
	config *config.Config,
) *ShowExecutor {
	exe := ShowExecutor{
		FlagSet:       flag.NewFlagSet("show", flag.ExitOnError),
		printer:       printer,
		config:        config,
		accountName:   internalFlag.NewStringSliceValue(),
		attachmentIDs: internalFlag.NewStringSliceValue(),
	}

	exe.Usage = usage.ExecutorUsageFunc("show", "Shows details about a specified resource", exe.FlagSet)

	exe.Var(&exe.accountName, "account-name", "The name of the account")
	exe.BoolVar(&exe.getAllImages, "all-images", false, "Set to true to show all images from a status")
	exe.BoolVar(&exe.getAllVideos, "all-videos", false, "Set to true to show all videos from a status")
	exe.Var(&exe.attachmentIDs, "attachment-id", "The ID of the media attachment")
	exe.BoolVar(&exe.showInBrowser, "browser", false, "Set to true to view in the your favourite browser")
	exe.BoolVar(&exe.excludeBoosts, "exclude-boosts", false, "Set to true to exclude statuses that are boosts of another status")
	exe.BoolVar(&exe.excludeReplies, "exclude-replies", false, "Set to true to exclude statuses that are a reply to another status")
	exe.StringVar(&exe.fromResourceType, "from", "", "Specify the resource type to action the target resource from")
	exe.IntVar(&exe.limit, "limit", 20, "Specify the limit of items to display")
	exe.StringVar(&exe.listID, "list-id", "", "The ID of the list in question")
	exe.BoolVar(&exe.myAccount, "my-account", false, "Set to true to specify your account")
	exe.BoolVar(&exe.onlyMedia, "only-media", false, "Set to true to show only the statuses with media attachments")
	exe.BoolVar(&exe.onlyPinned, "only-pinned", false, "Set to true to show only the account's pinned statuses")
	exe.BoolVar(&exe.onlyPublic, "only-public", false, "Set to true to show only the account's public posts")
	exe.StringVar(&exe.pollID, "poll-id", "", "The ID of the poll")
	exe.BoolVar(&exe.showUserPreferences, "show-preferences", false, "Set to true to view your posting preferences when viewing your account information")
	exe.BoolVar(&exe.showStatuses, "show-statuses", false, "Set to true to view the statuses created from the account you are viewing")
	exe.BoolVar(&exe.skipAccountRelationship, "skip-relationship", false, "Set to true to skip showing your relationship to the account that you are viewing")
	exe.StringVar(&exe.statusID, "status-id", "", "The ID of the status")
	exe.StringVar(&exe.timelineCategory, "timeline-category", "home", "The timeline category")
	exe.StringVar(&exe.tag, "tag", "", "The name of the tag")
	exe.StringVar(&exe.resourceType, "type", "", "The type of resource you want to action on (e.g. account, status)")

	return &exe
}

// SwitchExecutor is the executor for the switch command.
type SwitchExecutor struct {
	*flag.FlagSet
	printer     *printer.Printer
	config      *config.Config
	accountName internalFlag.StringSliceValue
	to          string
}

func NewSwitchExecutor(
	printer *printer.Printer,
	config *config.Config,
) *SwitchExecutor {
	exe := SwitchExecutor{
		FlagSet:     flag.NewFlagSet("switch", flag.ExitOnError),
		printer:     printer,
		config:      config,
		accountName: internalFlag.NewStringSliceValue(),
	}

	exe.Usage = usage.ExecutorUsageFunc("switch", "Performs a switch operation (e.g. switching between logged in accounts)", exe.FlagSet)

	exe.Var(&exe.accountName, "account-name", "The name of the account")
	exe.StringVar(&exe.to, "to", "", "TBC")

	return &exe
}

// UnblockExecutor is the executor for the unblock command.
type UnblockExecutor struct {
	*flag.FlagSet
	printer      *printer.Printer
	config       *config.Config
	accountName  internalFlag.StringSliceValue
	resourceType string
}

func NewUnblockExecutor(
	printer *printer.Printer,
	config *config.Config,
) *UnblockExecutor {
	exe := UnblockExecutor{
		FlagSet:     flag.NewFlagSet("unblock", flag.ExitOnError),
		printer:     printer,
		config:      config,
		accountName: internalFlag.NewStringSliceValue(),
	}

	exe.Usage = usage.ExecutorUsageFunc("unblock", "Unblocks a resource (e.g. an account)", exe.FlagSet)

	exe.Var(&exe.accountName, "account-name", "The name of the account")
	exe.StringVar(&exe.resourceType, "type", "", "The type of resource you want to action on (e.g. account, status)")

	return &exe
}

// UnfollowExecutor is the executor for the unfollow command.
type UnfollowExecutor struct {
	*flag.FlagSet
	printer      *printer.Printer
	config       *config.Config
	accountName  internalFlag.StringSliceValue
	resourceType string
}

func NewUnfollowExecutor(
	printer *printer.Printer,
	config *config.Config,
) *UnfollowExecutor {
	exe := UnfollowExecutor{
		FlagSet:     flag.NewFlagSet("unfollow", flag.ExitOnError),
		printer:     printer,
		config:      config,
		accountName: internalFlag.NewStringSliceValue(),
	}

	exe.Usage = usage.ExecutorUsageFunc("unfollow", "Unfollows a resource (e.g. an account)", exe.FlagSet)

	exe.Var(&exe.accountName, "account-name", "The name of the account")
	exe.StringVar(&exe.resourceType, "type", "", "The type of resource you want to action on (e.g. account, status)")

	return &exe
}

// UnmuteExecutor is the executor for the unmute command.
type UnmuteExecutor struct {
	*flag.FlagSet
	printer      *printer.Printer
	config       *config.Config
	accountName  internalFlag.StringSliceValue
	resourceType string
}

func NewUnmuteExecutor(
	printer *printer.Printer,
	config *config.Config,
) *UnmuteExecutor {
	exe := UnmuteExecutor{
		FlagSet:     flag.NewFlagSet("unmute", flag.ExitOnError),
		printer:     printer,
		config:      config,
		accountName: internalFlag.NewStringSliceValue(),
	}

	exe.Usage = usage.ExecutorUsageFunc("unmute", "Umutes a specific resource (e.g. an account)", exe.FlagSet)

	exe.Var(&exe.accountName, "account-name", "The name of the account")
	exe.StringVar(&exe.resourceType, "type", "", "The type of resource you want to action on (e.g. account, status)")

	return &exe
}

// VersionExecutor is the executor for the version command.
type VersionExecutor struct {
	*flag.FlagSet
	printer *printer.Printer
	full    bool
}

func NewVersionExecutor(
	printer *printer.Printer,
) *VersionExecutor {
	exe := VersionExecutor{
		FlagSet: flag.NewFlagSet("version", flag.ExitOnError),
		printer: printer,
	}

	exe.Usage = usage.ExecutorUsageFunc("version", "Prints the application's version and build information", exe.FlagSet)

	exe.BoolVar(&exe.full, "full", false, "Set to true to print the build information in full")

	return &exe
}

// WhoamiExecutor is the executor for the whoami command.
type WhoamiExecutor struct {
	*flag.FlagSet
	printer *printer.Printer
	config  *config.Config
}

func NewWhoamiExecutor(
	printer *printer.Printer,
	config *config.Config,
) *WhoamiExecutor {
	exe := WhoamiExecutor{
		FlagSet: flag.NewFlagSet("whoami", flag.ExitOnError),
		printer: printer,
		config:  config,
	}

	exe.Usage = usage.ExecutorUsageFunc("whoami", "Prints the account that you are currently logged into", exe.FlagSet)

	return &exe
}