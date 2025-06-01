package executor

import "fmt"

type unsupportedActionError struct {
	action string
	target string
}

func (e unsupportedActionError) Error() string {
	return "unsupported action (" +
		e.action +
		") on target (" +
		e.target +
		")"
}

type unrecognisedTargetError struct {
	target string
}

func (e unrecognisedTargetError) Error() string {
	return "unrecognised target: " + e.target
}

type missingIDError struct {
	target string
	action string
}

func (e missingIDError) Error() string {
	return "please provide the ID of the " + e.target + " you want to " + e.action
}

type missingValueError struct {
	valueType string
	target    string
	action    string
}

func (e missingValueError) Error() string {
	return "please specify the " + e.valueType +
		" of the " + e.target +
		" you want to " + e.action
}

type unsupportedTargetToTargetError struct {
	action        string
	focusedTarget string
	preposition   string
	relatedTarget string
}

func (e unsupportedTargetToTargetError) Error() string {
	return "'" +
		e.action + " " +
		e.focusedTarget + " " +
		e.preposition + " " +
		e.relatedTarget +
		"' is not a supported operation"
}

type forbiddenActionOnStatusError struct {
	action              string
	includeNotMentioned bool
}

func (e forbiddenActionOnStatusError) Error() string {
	msg := "unable to " + e.action + " the status because you are not the owner"

	if e.includeNotMentioned {
		msg += " and you are not mentioned in it"
	}

	return msg
}

type missingAccountInCredentialsError struct{}

func (e missingAccountInCredentialsError) Error() string {
	return "this account is not present in the credentials file"
}

type zeroValuesError struct {
	valueType string
	action    string
}

func (e zeroValuesError) Error() string {
	msg := "please specify one or more " + e.valueType + "(s)"

	if e.action != "" {
		msg += " to " + e.action
	}

	return msg
}

type notFollowingError struct {
	account string
}

func (e notFollowingError) Error() string {
	return "you are not following " + e.account
}

type loginNoInstanceError struct{}

func (e loginNoInstanceError) Error() string {
	return "please specify the instance that you want to log into"
}

type pollMissingError struct{}

func (e pollMissingError) Error() string {
	return "this status does not have a poll"
}

type pollClosedError struct{}

func (e pollClosedError) Error() string {
	return "this poll is closed"
}

type pollNoMultipleChoiceError struct{}

func (e pollNoMultipleChoiceError) Error() string {
	return "this poll does not allow multiple choices"
}

type voteInOwnPollError struct{}

func (e voteInOwnPollError) Error() string {
	return "you cannot vote in your own poll"
}

type noPollOptionsError struct{}

func (e noPollOptionsError) Error() string {
	return "no options were provided for this poll"
}

type missingSearchQueryError struct{}

func (e missingSearchQueryError) Error() string {
	return "please enter a search query"
}

type missingMediaFileError struct{}

func (e missingMediaFileError) Error() string {
	return "please provide the path to the media file"
}

type noContentOrMediaError struct{}

func (e noContentOrMediaError) Error() string {
	return "please add content or attach at least one media to the status that you want to create"
}

type statusHasPollAndMediaError struct{}

func (e statusHasPollAndMediaError) Error() string {
	return "you cannot create a status with both a poll and media attachments"
}

type mismatchedMediaFlagsError struct {
	kind string
	want int
	got  int
}

func (e mismatchedMediaFlagsError) Error() string {
	return fmt.Sprintf(
		"the number of %s provided does not match the number of media files provided: want %d, got %d",
		e.kind,
		e.want,
		e.got,
	)
}

type zeroConfigurationError struct {
	path string
}

func (e zeroConfigurationError) Error() string {
	return "configuration not set: please ensure that the configuration file is present at '" + e.path + "'"
}

type aliasActionKeywordError struct {
	alias string
}

func (e aliasActionKeywordError) Error() string {
	return "'" + e.alias + "' is a built-in action keyword"
}

type aliasBuiltinAliasError struct {
	alias string
}

func (e aliasBuiltinAliasError) Error() string {
	return "'" + e.alias + "' is a built-in alias"
}

type aliasNewNameUnsetError struct{}

func (e aliasNewNameUnsetError) Error() string {
	return "please specify the alias' new name"
}

type invalidTimelineCategoryError struct {
	category string
}

func (e invalidTimelineCategoryError) Error() string {
	return "'" + e.category + "' is not a valid timeline category"
}

type usageNoOpFoundError struct {
	target string
}

func (e usageNoOpFoundError) Error() string {
	return "no operations found for '" + e.target + "'"
}

type usageNoOpForTargetError struct {
	operation string
	target    string
}

func (e usageNoOpForTargetError) Error() string {
	return "unable to find the operation '" +
		e.operation +
		"' for the target '" +
		e.target +
		"'"
}
