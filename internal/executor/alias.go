package executor

import (
	"fmt"

	"codeflow.dananglin.me.uk/apollo/enbas/internal/cli"
	"codeflow.dananglin.me.uk/apollo/enbas/internal/command"
	"codeflow.dananglin.me.uk/apollo/enbas/internal/config"
	"codeflow.dananglin.me.uk/apollo/enbas/internal/printer"
)

// aliasFunc is the function for the 'alias' target for interacting
// with a single alias.
func aliasFunc(
	cfg config.Config,
	printSettings printer.Settings,
	cmd command.Command,
) error {
	if cfg.IsZero() {
		return zeroConfigurationError{path: cfg.Path}
	}

	switch cmd.Action {
	case cli.ActionCreate:
		return aliasCreate(
			cfg.Path,
			printSettings,
			cmd.FocusedTargetFlags,
		)
	case cli.ActionDelete:
		return aliasDelete(
			cfg.Path,
			printSettings,
			cmd.FocusedTargetFlags,
		)
	case cli.ActionEdit:
		return aliasEdit(
			cfg.Path,
			printSettings,
			cmd.FocusedTargetFlags,
		)
	case cli.ActionRename:
		return aliasRename(
			cfg.Path,
			printSettings,
			cmd.FocusedTargetFlags,
		)
	default:
		return unsupportedActionError{action: cmd.Action, target: cli.TargetAlias}
	}
}

func aliasCreate(
	configFilepath string,
	printSettings printer.Settings,
	flags []string,
) error {
	var (
		name      string
		arguments string
	)

	// Parse the remaining flags.
	if err := cli.ParseAliasCreateFlags(
		&name,
		&arguments,
		flags,
	); err != nil {
		return err //nolint:wrapcheck
	}

	if name == "" {
		return aliasNameUnsetError{}
	}

	if !command.ValidAlias(name) {
		return fmt.Errorf(
			"alias validation error: %w",
			command.NewInvalidAliasError(name),
		)
	}

	if arguments == "" {
		return aliasArgsUnsetError{}
	}

	if _, exists := cli.ActionDesc(name); exists {
		return aliasActionKeywordError{alias: name}
	}

	if _, exists := cli.BuiltInAlias(name); exists {
		return aliasBuiltinAliasError{alias: name}
	}

	if err := config.CreateAlias(
		configFilepath,
		name,
		arguments,
	); err != nil {
		return fmt.Errorf("error creating the alias: %w", err)
	}

	printer.PrintSuccess(printSettings, "Successfully created the alias.")

	return nil
}

func aliasEdit(
	configFilepath string,
	printSettings printer.Settings,
	flags []string,
) error {
	var (
		name      string
		arguments string
	)

	// Parse the remaining flags.
	if err := cli.ParseAliasEditFlags(
		&name,
		&arguments,
		flags,
	); err != nil {
		return err //nolint:wrapcheck
	}

	if name == "" {
		return aliasNameUnsetError{}
	}

	if arguments == "" {
		return aliasArgsUnsetError{}
	}

	if !command.ValidAlias(name) {
		return fmt.Errorf(
			"alias validation error: %w",
			command.NewInvalidAliasError(name),
		)
	}

	if _, exists := cli.ActionDesc(name); exists {
		return aliasActionKeywordError{alias: name}
	}

	if _, exists := cli.BuiltInAlias(name); exists {
		return aliasBuiltinAliasError{alias: name}
	}

	if err := config.EditAlias(
		configFilepath,
		name,
		arguments,
	); err != nil {
		return fmt.Errorf("error editing the alias: %w", err)
	}

	printer.PrintSuccess(printSettings, "Successfully edited the alias.")

	return nil
}

func aliasDelete(
	configFilepath string,
	printSettings printer.Settings,
	flags []string,
) error {
	var name string

	// Parse the remaining flags.
	if err := cli.ParseAliasDeleteFlags(
		&name,
		flags,
	); err != nil {
		return err //nolint:wrapcheck
	}

	if name == "" {
		return aliasNameUnsetError{}
	}

	if err := config.DeleteAlias(
		configFilepath,
		name,
	); err != nil {
		return fmt.Errorf("error deleting the alias: %w", err)
	}

	printer.PrintSuccess(printSettings, "Successfully deleted the alias.")

	return nil
}

func aliasRename(
	configFilepath string,
	printSettings printer.Settings,
	flags []string,
) error {
	var (
		oldName string
		newName string
	)

	// Parse the remaining flags.
	if err := cli.ParseAliasRenameFlags(
		&oldName,
		&newName,
		flags,
	); err != nil {
		return err //nolint:wrapcheck
	}

	if oldName == "" {
		return aliasOldNameUnsetError{}
	}

	if newName == "" {
		return aliasNewNameUnsetError{}
	}

	if newName == oldName {
		return aliasNewNameUnsetError{}
	}

	if !command.ValidAlias(newName) {
		return fmt.Errorf(
			"alias validation error: %w",
			command.NewInvalidAliasError(newName),
		)
	}

	if _, exists := cli.ActionDesc(newName); exists {
		return aliasActionKeywordError{alias: newName}
	}

	if _, exists := cli.BuiltInAlias(newName); exists {
		return aliasBuiltinAliasError{alias: newName}
	}

	if err := config.RenameAlias(
		configFilepath,
		oldName,
		newName,
	); err != nil {
		return fmt.Errorf("error renaming the alias: %w", err)
	}

	printer.PrintSuccess(printSettings, "Successfully renamed the alias.")

	return nil
}
