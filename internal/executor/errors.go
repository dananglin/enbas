// SPDX-FileCopyrightText: 2024 Dan Anglin <d.n.i.anglin@gmail.com>
//
// SPDX-License-Identifier: GPL-3.0-or-later

package executor

type FlagNotSetError struct {
	flagText string
}

func (e FlagNotSetError) Error() string {
	return "please use the required --" + e.flagText + " flag"
}

type UnsupportedTypeError struct {
	resourceType string
}

func (e UnsupportedTypeError) Error() string {
	return "'" + e.resourceType + "' is not supported for this operation"
}

type NoAccountSpecifiedError struct{}

func (e NoAccountSpecifiedError) Error() string {
	return "no account specified in this request"
}

type UnsupportedAddOperationError struct {
	ResourceType      string
	AddToResourceType string
}

func (e UnsupportedAddOperationError) Error() string {
	return "adding '" +
		e.ResourceType +
		"' to '" +
		e.AddToResourceType +
		"' is not supported"
}

type UnsupportedRemoveOperationError struct {
	ResourceType           string
	RemoveFromResourceType string
}

func (e UnsupportedRemoveOperationError) Error() string {
	return "removing '" +
		e.ResourceType +
		"' from '" +
		e.RemoveFromResourceType +
		"' is not supported"
}

type EmptyContentError struct {
	ResourceType string
	Hint         string
}

func (e EmptyContentError) Error() string {
	message := "the content of this " + e.ResourceType + " should not be empty"

	if e.Hint != "" {
		message += ", " + e.Hint
	}

	return message
}

type UnknownCommandError struct {
	Command string
}

func (e UnknownCommandError) Error() string {
	return "unknown command '" + e.Command + "'"
}

type PollClosedError struct{}

func (e PollClosedError) Error() string {
	return "this poll is closed"
}

type MultipleChoiceError struct{}

func (e MultipleChoiceError) Error() string {
	return "this poll does not allow multiple choices"
}

type NoPollOptionError struct{}

func (e NoPollOptionError) Error() string {
	return "no options were provided for this poll, please use the --" +
		flagPollOption +
		" flag to add options to the poll"
}

type NotFollowingError struct {
	Account string
}

func (e NotFollowingError) Error() string {
	return "you are not following " + e.Account
}
