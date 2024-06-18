// SPDX-FileCopyrightText: 2024 Dan Anglin <d.n.i.anglin@gmail.com>
//
// SPDX-License-Identifier: GPL-3.0-or-later

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
	flagAttachmentID              = "attachment-id"
	flagBrowser                   = "browser"
	flagChoose                    = "choose"
	flagContentType               = "content-type"
	flagContent                   = "content"
	flagEnableFederation          = "enable-federation"
	flagEnableLikes               = "enable-likes"
	flagEnableReplies             = "enable-replies"
	flagEnableReposts             = "enable-reposts"
	flagFrom                      = "from"
	flagFromFile                  = "from-file"
	flagFull                      = "full"
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
	flagPollAllowsMultipleChoices = "poll-allows-multiple-choices"
	flagPollExpiresIn             = "poll-expires-in"
	flagPollHidesVoteCounts       = "poll-hides-vote-counts"
	flagPollID                    = "poll-id"
	flagPollOption                = "poll-option"
	flagSensitive                 = "sensitive"
	flagSkipRelationship          = "skip-relationship"
	flagShowPreferences           = "show-preferences"
	flagShowReposts               = "show-reposts"
	flagSpoilerText               = "spoiler-text"
	flagStatusID                  = "status-id"
	flagTag                       = "tag"
	flagTimelineCategory          = "timeline-category"
	flagTo                        = "to"
	flagType                      = "type"
	flagVisibility                = "visibility"
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
