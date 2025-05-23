package printer

import "codeflow.dananglin.me.uk/apollo/enbas/internal/info"

func PrintUsageApp(
	settings Settings,
	targets map[string]string,
	flags map[string]string,
) error {
	data := struct {
		BinaryVersion string
		Name          string
		Targets       map[string]string
		Flags         map[string]string
	}{
		BinaryVersion: info.BinaryVersion,
		Name:          info.ApplicationName,
		Targets:       targets,
		Flags:         flags,
	}

	return renderTemplateToStdout(
		settings,
		"helpApp",
		"",
		data,
	)
}

func PrintUsageTarget(
	settings Settings,
	target string,
	description string,
	actions map[string]string,
	topLevelFlags map[string]string,
) error {
	data := struct {
		AppName       string
		Target        string
		Description   string
		Actions       map[string]string
		TopLevelFlags map[string]string
	}{
		AppName:       info.ApplicationName,
		Target:        target,
		Description:   description,
		Actions:       actions,
		TopLevelFlags: topLevelFlags,
	}

	return renderTemplateToStdout(
		settings,
		"helpTarget",
		"",
		data,
	)
}

func PrintUsageAction(
	settings Settings,
	action string,
	description string,
	availableTargets []string,
	topLevelFlags map[string]string,
) error {
	data := struct {
		AppName           string
		Action            string
		ActionDescription string
		AvailableTargets  []string
		TopLevelFlags     map[string]string
	}{
		AppName:           info.ApplicationName,
		Action:            action,
		ActionDescription: description,
		AvailableTargets:  availableTargets,
		TopLevelFlags:     topLevelFlags,
	}

	return renderTemplateToStdout(
		settings,
		"helpAction",
		"",
		data,
	)
}

func PrintUsageTargetAction(
	settings Settings,
	target string,
	action string,
	description string,
	flags map[string]string,
	topLevelFlags map[string]string,
) error {
	data := struct {
		AppName       string
		Target        string
		Action        string
		Description   string
		Flags         map[string]string
		TopLevelFlags map[string]string
	}{
		AppName:       info.ApplicationName,
		Target:        target,
		Action:        action,
		Description:   description,
		Flags:         flags,
		TopLevelFlags: topLevelFlags,
	}

	return renderTemplateToStdout(
		settings,
		"helpTargetAction",
		"",
		data,
	)
}
