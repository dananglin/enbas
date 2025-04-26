package executor

import (
	"fmt"
	"maps"
	"strings"

	"codeflow.dananglin.me.uk/apollo/enbas/internal/cli"
	"codeflow.dananglin.me.uk/apollo/enbas/internal/command"
	"codeflow.dananglin.me.uk/apollo/enbas/internal/printer"
)

// helpFunc is the function for the help target for printing the
// help documentation to the screen for the user.
func helpFunc(
	opts topLevelOpts,
	cmd command.Command,
) error {
	printSettings := printer.NewSettings(opts.noColor, "", 0)

	switch cmd.Action {
	case cli.ActionShow:
		return helpShow(printSettings, cmd.FocusedTargetFlags)
	default:
		return unsupportedActionError{action: cmd.Action, target: cli.TargetHelp}
	}
}

func helpShow(
	printSettings printer.Settings,
	flags []string,
) error {
	var (
		action string
		target string
	)

	// Parse the remaining flags.
	if err := cli.ParseHelpShowFlags(
		&action,
		&target,
		flags,
	); err != nil {
		return err
	}

	switch {
	case (target == "") && (action == ""):
		return helpShowApp(printSettings)
	case (target != "") && (action == ""):
		return helpShowTarget(printSettings, target)
	case (target == "") && (action != ""):
		return helpShowAction(printSettings, action)
	default:
		return helpShowTargetAction(printSettings, target, action)
	}
}

func helpShowApp(printSettings printer.Settings) error {
	if err := printer.PrintHelpApp(
		printSettings,
		cli.TargetDescMap(),
		cli.TopLevelFlagsUsageMap(),
	); err != nil {
		return fmt.Errorf("error printing the help documentation: %w", err)
	}

	return nil
}

func helpShowTarget(
	printSettings printer.Settings,
	target string,
) error {
	desc, ok := cli.TargetDesc(target)
	if !ok {
		return unrecognisedTargetError{target: target}
	}

	if err := printer.PrintHelpTarget(
		printSettings,
		target,
		desc,
		renderTemplatesinMap(cli.TargetActions(target), target, ""),
		cli.TopLevelFlagsUsageMap(),
	); err != nil {
		return fmt.Errorf("error printing the help documentation: %w", err)
	}

	return nil
}

func helpShowAction(
	printSettings printer.Settings,
	action string,
) error {
	desc, ok := cli.ActionDesc(action)
	if !ok {
		return unrecognisedActionError{action: action}
	}

	availableTargets := cli.ActionTargets(action)

	if err := printer.PrintHelpAction(
		printSettings,
		action,
		renderTargetTemplate(desc, "target"),
		availableTargets,
		cli.TopLevelFlagsUsageMap(),
	); err != nil {
		return fmt.Errorf("error printing the help documentation: %w", err)
	}

	return nil
}

func helpShowTargetAction(
	printSettings printer.Settings,
	target string,
	action string,
) error {
	if _, ok := cli.TargetDesc(target); !ok {
		return unrecognisedTargetError{target: target}
	}

	desc, ok := cli.ActionDesc(action)
	if !ok {
		return unrecognisedActionError{action: action}
	}

	flags, ok := cli.TargetActionFlags(target, action)
	if !ok {
		return unsupportedActionError{
			action: action,
			target: target,
		}
	}

	if err := printer.PrintHelpTargetAction(
		printSettings,
		target,
		action,
		renderTargetTemplate(desc, target),
		renderTemplatesinMap(flags, target, action),
		cli.TopLevelFlagsUsageMap(),
	); err != nil {
		return fmt.Errorf("error printing the help documentation: %w", err)
	}

	return nil
}

func renderTemplatesinMap(
	input map[string]string,
	target string,
	action string,
) map[string]string {
	output := make(map[string]string)

	for key := range maps.All(input) {
		output[key] = renderTemplate(input[key], target, action)
	}

	return output
}

func renderTemplate(text, target, action string) string {
	return renderTargetTemplate(renderActionTemplate(text, action), target)
}

func renderTargetTemplate(text, target string) string {
	if target == "" {
		return text
	}

	return strings.ReplaceAll(text, "{target}", target)
}

func renderActionTemplate(text, action string) string {
	if action == "" {
		return text
	}

	return strings.ReplaceAll(text, "{action}", action)
}
