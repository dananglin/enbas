package printer

import (
	"codeflow.dananglin.me.uk/apollo/enbas/internal/cli"
	"codeflow.dananglin.me.uk/apollo/enbas/internal/info"
)

func PrintUsageRoot(
	settings Settings,
	targets map[string]string,
	flags map[string]string,
	exampleTarget string,
) error {
	data := struct {
		BinaryVersion string
		Name          string
		Targets       map[string]string
		Flags         map[string]string
		ExampleTarget string
	}{
		BinaryVersion: info.BinaryVersion,
		Name:          info.ApplicationName,
		Targets:       targets,
		Flags:         flags,
		ExampleTarget: exampleTarget,
	}

	return renderTemplateToPager(
		settings,
		"usageRoot",
		"",
		data,
	)
}

func PrintUsageTarget(
	settings Settings,
	target string,
	description string,
	topLevelFlags map[string]string,
	operations map[string]cli.UsageOperation,
) error {
	data := struct {
		AppName       string
		Target        string
		Description   string
		TopLevelFlags map[string]string
		Operations    map[string]cli.UsageOperation
	}{
		AppName:       info.ApplicationName,
		Target:        target,
		Description:   description,
		TopLevelFlags: topLevelFlags,
		Operations:    operations,
	}

	return renderTemplateToPager(
		settings,
		"usageTarget",
		"",
		data,
	)
}

func PrintUsageOperation(
	settings Settings,
	topLevelFlags map[string]string,
	name string,
	operation cli.UsageOperation,
	flagDescriptions []string,
) error {
	if len(operation.Flags) != len(flagDescriptions) {
		return NumFlagDescriptionsError{
			NumFlags:            len(operation.Flags),
			NumFlagDescriptions: len(flagDescriptions),
		}
	}

	type flag struct {
		Name  string
		Usage string
	}

	flags := make([]flag, 0)

	for idx := range operation.Flags {
		flags = append(
			flags,
			flag{
				Name:  operation.Flags[idx],
				Usage: flagDescriptions[idx],
			},
		)
	}

	data := struct {
		AppName       string
		TopLevelFlags map[string]string
		Name          string
		Operation     cli.UsageOperation
		Flags         []flag
	}{
		AppName:       info.ApplicationName,
		TopLevelFlags: topLevelFlags,
		Name:          name,
		Operation:     operation,
		Flags:         flags,
	}

	return renderTemplateToPager(
		settings,
		"usageOperation",
		"",
		data,
	)
}
