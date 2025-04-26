package command

import (
	"slices"
	"strings"

	"codeflow.dananglin.me.uk/apollo/enbas/internal/cli"
)

func Parse(args []string) (Command, error) {
	if len(args) == 0 {
		return Command{}, NewNoActionError()
	}

	// No target detected after action specified
	if len(args) == 1 {
		return Command{}, NewNoFocusedTargetError(args[0])
	}

	if strings.HasPrefix(args[1], "-") {
		// The second argument is a flag.
		// Flags after an action word is not supported.
		action := args[0]
		flag := args[1]
		return Command{}, NewFlagAfterActionError(action, flag)
	}

	// Start constructing the command value
	cmd := Command{
		Action:             args[0],
		FocusedTarget:      args[1],
		FocusedTargetFlags: []string{},
		Preposition:        "",
		RelatedTarget:      "",
		RelatedTargetFlags: []string{},
	}

	remaining := args[2:]
	preposition := cli.TargetActionPreposition(cmd.FocusedTarget, cmd.Action)

	if len(remaining) == 0 && preposition == "" {
		return cmd, nil
	}

	// Here we are expecting a preposition word in the argument list
	// but there are no more remaining arguments therefore the preposition
	// keyword is missing and the command is incomplete.
	if len(remaining) == 0 && preposition != "" {
		return Command{}, NewPrepositionKeywordMissingError(preposition)
	}

	if preposition == "" {
		cmd.FocusedTargetFlags = remaining

		return cmd, nil
	}

	cmd.Preposition = preposition

	extracted, err := extractTargetFlags(remaining, preposition)
	if err != nil {
		return Command{}, err
	}

	cmd.FocusedTargetFlags = extracted.flags

	if len(extracted.remainingArgs) == 0 {
		return cmd, nil
	}

	// The next argument must not be a flag
	if strings.HasPrefix(args[0], "-") {
		return Command{}, NewFlagAfterPrepositionError(args[0], cmd.Preposition)
	}

	cmd.RelatedTarget, cmd.RelatedTargetFlags = extracted.remainingArgs[0], extracted.remainingArgs[1:]

	return cmd, nil
}

type extractTargetFlagsResult struct {
	flags         []string
	remainingArgs []string
}

func extractTargetFlags(args []string, preposition string) (extractTargetFlagsResult, error) {
	idx := slices.IndexFunc(args, func(s string) bool {
		return s == preposition
	})

	if idx < 0 {
		return extractTargetFlagsResult{}, NewPrepositionKeywordMissingError(preposition)
	}

	return extractTargetFlagsResult{
		flags:         args[:idx],
		remainingArgs: args[idx+1:],
	}, nil
}
