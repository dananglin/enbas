package executor

import "fmt"

type Error struct {
	message string
}

func (e Error) Error() string {
	return e.message
}

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
	resourceType      string
	addToResourceType string
}

func (e UnsupportedAddOperationError) Error() string {
	return "adding '" +
		e.resourceType +
		"' to '" +
		e.addToResourceType +
		"' is not supported"
}

type UnsupportedRemoveOperationError struct {
	resourceType           string
	removeFromResourceType string
}

func (e UnsupportedRemoveOperationError) Error() string {
	return "removing '" +
		e.resourceType +
		"' from '" +
		e.removeFromResourceType +
		"' is not supported"
}

type UnsupportedShowOperationError struct {
	resourceType         string
	showFromResourceType string
}

func (e UnsupportedShowOperationError) Error() string {
	return "showing '" +
		e.resourceType +
		"' from '" +
		e.showFromResourceType +
		"' is not supported"
}

type UnknownCommandError struct {
	command string
}

func (e UnknownCommandError) Error() string {
	return "unknown command '" + e.command + "'"
}

type NotFollowingError struct {
	account string
}

func (e NotFollowingError) Error() string {
	return "you are not following " + e.account
}

type MismatchedNumMediaValuesError struct {
	valueType     string
	numValues     int
	numMediaFiles int
}

func (e MismatchedNumMediaValuesError) Error() string {
	return fmt.Sprintf(
		"unexpected number of %s: received %d media files but got %d %s",
		e.valueType,
		e.numMediaFiles,
		e.numValues,
		e.valueType,
	)
}

type UnexpectedNumValuesError struct {
	name     string
	actual   int
	expected int
}

func (e UnexpectedNumValuesError) Error() string {
	return fmt.Sprintf(
		"received an unexpected number of %s: received %d, expected %d",
		e.name,
		e.actual,
		e.expected,
	)
}
