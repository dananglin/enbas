package command

type helpFlagDetectedError struct {
	action string
	target string
}

func (e helpFlagDetectedError) Error() string {
	msg := "help flag detected: please use 'enbas show help"

	if e.action != "" {
		msg += " --action " + e.action
	}

	if e.target != "" {
		msg += " --target " + e.target
	}

	msg += "' instead"

	return msg
}

func NewHelpFlagDetectedError(action, target string) error {
	return helpFlagDetectedError{action: action, target: target}
}

type noActionError struct{}

func (e noActionError) Error() string {
	return "please specify an action"
}

func NewNoActionError() error {
	return noActionError{}
}

type noFocusedTargetError struct {
	action string
}

func (e noFocusedTargetError) Error() string {
	return "please specify a target to " + e.action
}

func NewNoFocusedTargetError(action string) error {
	return noFocusedTargetError{action: action}
}

type noRelatedTargetError struct {
	action        string
	focusedTarget string
	preposition   string
}

func (e noRelatedTargetError) Error() string {
	return "please specify a target to " +
		e.action +
		" the " +
		e.focusedTarget + " " + e.preposition
}

func NewNoRelatedTargetError(action, focusedTarget, preposition string) error {
	return noRelatedTargetError{action: action, focusedTarget: focusedTarget, preposition: preposition}
}

type flagAfterActionError struct {
	action string
	flag   string
}

func (e flagAfterActionError) Error() string {
	return "the flag (" + e.flag + ") was specified after the action word (" + e.action + ")"
}

func NewFlagAfterActionError(action, flag string) error {
	return flagAfterActionError{action: action, flag: flag}
}

type flagAfterPrepositionError struct {
	preposition string
	flag        string
}

func (e flagAfterPrepositionError) Error() string {
	return "the flag (" + e.flag + ") was specified after the preposition word (" + e.preposition + ")"
}

func NewFlagAfterPrepositionError(preposition, flag string) error {
	return flagAfterPrepositionError{preposition: preposition, flag: flag}
}

type prepositionKeywordMissingError struct {
	preposition string
}

func (e prepositionKeywordMissingError) Error() string {
	return "the preposition keyword \"" + e.preposition + "\" was not found"
}

func NewPrepositionKeywordMissingError(preposition string) error {
	return prepositionKeywordMissingError{preposition: preposition}
}
