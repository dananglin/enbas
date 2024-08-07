package executor

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

const (
	flagAddPoll                   = "add-poll"
	flagAccountName               = "account-name"
	flagAllImages                 = "all-images"
	flagAllVideos                 = "all-videos"
	flagAttachmentID              = "attachment-id"
	flagBrowser                   = "browser"
	flagContentType               = "content-type"
	flagContent                   = "content"
	flagEnableFederation          = "enable-federation"
	flagEnableLikes               = "enable-likes"
	flagEnableReplies             = "enable-replies"
	flagEnableReposts             = "enable-reposts"
	flagExcludeBoosts             = "exclude-boosts"
	flagExcludeReplies            = "exclude-replies"
	flagFrom                      = "from"
	flagFromFile                  = "from-file"
	flagFull                      = "full"
	flagInReplyTo                 = "in-reply-to"
	flagInstance                  = "instance"
	flagLanguage                  = "language"
	flagLimit                     = "limit"
	flagListID                    = "list-id"
	flagListTitle                 = "list-title"
	flagListRepliesPolicy         = "list-replies-policy"
	flagMyAccount                 = "my-account"
	flagMuteDuration              = "mute-duration"
	flagMuteNotifications         = "mute-notifications"
	flagNotify                    = "notify"
	flagOnlyMedia                 = "only-media"
	flagOnlyPinned                = "only-pinned"
	flagOnlyPublic                = "only-public"
	flagPollAllowsMultipleChoices = "poll-allows-multiple-choices"
	flagPollExpiresIn             = "poll-expires-in"
	flagPollHidesVoteCounts       = "poll-hides-vote-counts"
	flagPollID                    = "poll-id"
	flagPollOption                = "poll-option"
	flagSensitive                 = "sensitive"
	flagSkipRelationship          = "skip-relationship"
	flagShowPreferences           = "show-preferences"
	flagShowReposts               = "show-reposts"
	flagShowStatuses              = "show-statuses"
	flagSpoilerText               = "spoiler-text"
	flagStatusID                  = "status-id"
	flagTag                       = "tag"
	flagTimelineCategory          = "timeline-category"
	flagTo                        = "to"
	flagType                      = "type"
	flagVisibility                = "visibility"
	flagVote                      = "vote"
)

type MultiStringFlagValue []string

func (v *MultiStringFlagValue) String() string {
	return strings.Join(*v, ", ")
}

func (v *MultiStringFlagValue) Set(value string) error {
	if len(value) > 0 {
		*v = append(*v, value)
	}

	return nil
}

type MultiIntFlagValue []int

func (v *MultiIntFlagValue) String() string {
	value := "Choices: "

	for ind, vote := range *v {
		if ind == len(*v)-1 {
			value += strconv.Itoa(vote)
		} else {
			value += strconv.Itoa(vote) + ", "
		}
	}

	return value
}

func (v *MultiIntFlagValue) Set(text string) error {
	value, err := strconv.Atoi(text)
	if err != nil {
		return fmt.Errorf("unable to parse the value to an integer: %w", err)
	}

	*v = append(*v, value)

	return nil
}

type TimeDurationFlagValue struct {
	Duration time.Duration
}

func (v TimeDurationFlagValue) String() string {
	return ""
}

func (v *TimeDurationFlagValue) Set(text string) error {
	duration, err := time.ParseDuration(text)
	if err != nil {
		return fmt.Errorf("unable to parse the value as time duration: %w", err)
	}

	v.Duration = duration

	return nil
}
