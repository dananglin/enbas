package command

import (
	"regexp"
	"slices"
)

func (c Command) Validate() error {
	if err := c.ensureAction(); err != nil {
		return err
	}

	if err := c.ensureFocusedTarget(); err != nil {
		return err
	}

	if err := c.ensureRelatedTarget(); err != nil {
		return err
	}

	if err := c.ensureNoHelpFlag(); err != nil {
		return err
	}

	return nil
}

func (c Command) ensureAction() error {
	if c.Action == "" {
		return NewNoActionError()
	}

	return nil
}

func (c Command) ensureFocusedTarget() error {
	if c.FocusedTarget == "" {
		return NewNoFocusedTargetError(c.Action)
	}

	return nil
}

func (c Command) ensureRelatedTarget() error {
	if c.Preposition == "" {
		return nil
	}

	if c.RelatedTarget == "" {
		return NewNoRelatedTargetError(c.Action, c.FocusedTarget, c.Preposition)
	}

	return nil
}

// ensureNoHelpFlag returns true if any of the help flags
// (-h, --h, -help, --help) is detected.
func (c Command) ensureNoHelpFlag() error {
	pattern := regexp.MustCompile(`^--?h(?:elp)?$`)

	allFlags := slices.Concat(c.FocusedTargetFlags, c.RelatedTargetFlags)

	for _, text := range slices.All(allFlags) {
		if pattern.Match([]byte(text)) {
			return NewHelpFlagDetectedError(c.Action, c.FocusedTarget)
		}
	}

	return nil
}
