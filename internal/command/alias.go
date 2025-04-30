package command

import (
	"regexp"
	"strings"

	"codeflow.dananglin.me.uk/apollo/enbas/internal/cli"
)

func ExtractArgsFromAlias(
	command []string,
	userDefinedAliases map[string]string,
) ([]string, error) {
	if len(command) == 0 {
		return []string{}, NewAliasNoArgsError()
	}

	alias := command[0]

	// Check to see if the potential alias is actually built-in action word.
	// If it is, it is not an alias so return the original command.
	_, isAction := cli.ActionDesc(alias)
	if isAction {
		return command, nil
	}

	// Check to see if the potential alias matches one of the built-in alias.
	// If it is, build the final command using the result.
	parsedAlias, ok := cli.BuiltInAlias(alias)
	if ok {
		return buildFinalCommand(command, parsedAlias), nil
	}

	// At this point the potential alias may be defined by the user.
	// First, ensure that the alias is valid against our rules.
	if !validAlias(alias) {
		return []string{}, NewInvalidAliasError(alias)
	}

	// Now check to see if the potential alias is defined by the user
	outputStr, ok := userDefinedAliases[alias]
	if !ok {
		// The potential alias is unknown to the application and it also is
		// not an action word. Therefore command is invalid but this is
		// handled later in the process with the appropriate error handling.
		// So for now we return the original command back to the caller.
		return command, nil
	}

	// Return the results of the user defined alias.
	return buildFinalCommand(
		command,
		strings.Split(outputStr, " "),
	), nil
}

// ValidAlias checks to see if the given alias is valid.
// A valid alias is a least 3 characters long, can only contain lower
// cased alpha-numeric characters and hyphens, and must start and end
// with an alpha-numeric character.
func ValidAlias(alias string) bool {
	return validAlias(alias)
}

// See ValidAlias
func validAlias(alias string) bool {
	return regexp.MustCompile(`^[0-9a-z][0-9a-z-]+[0-9a-z]$`).MatchString(alias)
}

// buildFinalCommand by replacing the alias with it's corresponding extracted set of
// arguments in the original set of arguments.
func buildFinalCommand(originalArgs, extractedArgs []string) []string {
	return append(extractedArgs, originalArgs[1:]...)
}
