package printer

import (
	"embed"
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"text/template"

	"codeflow.dananglin.me.uk/apollo/enbas/internal/info"
	"codeflow.dananglin.me.uk/apollo/enbas/internal/model"
)

//go:embed templates/*
var templatesFS embed.FS

type Settings struct {
	noColor                bool
	lineWrapCharacterLimit int
	pager                  string
}

func NewSettings(
	noColor bool,
	pager string,
	lineWrapCharacterLimit int,
) Settings {
	if lineWrapCharacterLimit < minTerminalWidth {
		lineWrapCharacterLimit = minTerminalWidth
	}

	return Settings{
		noColor:                noColor,
		lineWrapCharacterLimit: lineWrapCharacterLimit,
		pager:                  pager,
	}
}

// PrintSuccess prints the successful message to standard output.
func PrintSuccess(settings Settings, text string) {
	const icon = "\u2714"

	success := boldgreen + icon + " " + reset
	if settings.noColor {
		success = icon + " "
	}

	printToStdout(success + " " + text + "\n")
}

// PrintFailure prints the failure message to standard error.
func PrintFailure(settings Settings, text string) {
	const icon = "\u2717"

	failure := boldred + icon + " " + reset
	if settings.noColor {
		failure = icon + " "
	}

	printToStderr(failure + " " + text + "\n")
}

// PrintInfo prints the message to standard output.
func PrintInfo(text string) {
	printToStdout(text)
}

// PrintVersion prints the binary build information.
func PrintVersion(settings Settings, showFullVersion bool) error {
	if !showFullVersion {
		printToStdout(info.ApplicationTitledName + " " + info.BinaryVersion + "\n")

		return nil
	}

	data := struct {
		Name          string
		BinaryVersion string
		GitCommit     string
		GoVersion     string
		BuildTime     string
	}{
		Name:          info.ApplicationTitledName,
		BinaryVersion: info.BinaryVersion,
		GitCommit:     info.GitCommit,
		GoVersion:     runtime.Version(),
		BuildTime:     info.BuildTime,
	}

	return renderTemplateToStdout(
		settings,
		"version",
		"",
		data,
	)
}

// PrintAccount prints the account details to the pager.
func PrintAccount(
	settings Settings,
	account model.Account,
	relationship model.AccountRelationship,
	preferences model.Preferences,
	statusList model.StatusList,
	myAccountID string,
) error {
	data := struct {
		Account      model.Account
		Relationship model.AccountRelationship
		Preferences  model.Preferences
		StatusList   model.StatusList
	}{
		Account:      account,
		Relationship: relationship,
		Preferences:  preferences,
		StatusList:   statusList,
	}

	return renderTemplateToPager(settings, "account", myAccountID, data)
}

// PrintAccountList prints the list of accounts to the pager.
func PrintAccountList(settings Settings, list model.AccountList) error {
	if list.BlockedAccounts {
		return renderTemplateToPager(settings, "blockedAccounts", "", list)
	}

	return renderTemplateToPager(settings, "accountList", "", list)
}

// PrintStatus prints the status to the pager.
func PrintStatus(
	settings Settings,
	status model.Status,
	myAccountID string,
	boostedBy model.AccountList,
	likedBy model.AccountList,
) error {
	data := struct {
		Status    model.Status
		BoostedBy model.AccountList
		LikedBy   model.AccountList
	}{
		Status:    status,
		BoostedBy: boostedBy,
		LikedBy:   likedBy,
	}

	return renderTemplateToPager(settings, "statusDoc", myAccountID, data)
}

// PrintStatusList prints a list of status cards to the pager.
func PrintStatusList(settings Settings, list model.StatusList, myAccountID string) error {
	return renderTemplateToPager(settings, "statusList", myAccountID, list)
}

// PrintInstance prints the instance information to the pager.
func PrintInstance(settings Settings, instance model.InstanceV2) error {
	return renderTemplateToPager(settings, "instance", "", instance)
}

// PrintMediaAttachment prints the details of the media attachment to the pager.
func PrintMediaAttachment(settings Settings, attachement model.MediaAttachment) error {
	return renderTemplateToPager(settings, "mediaAttachmentDoc", "", attachement)
}

// PrintTag prints the details of the tag to the pager.
func PrintTag(settings Settings, tag model.Tag) error {
	return renderTemplateToPager(settings, "tag", "", tag)
}

// PrintTagList prints the list of tags to the pager.
func PrintTagList(settings Settings, list model.TagList) error {
	return renderTemplateToPager(settings, "tagList", "", list)
}

// PrintThread prints the thread to the pager.
func PrintThread(settings Settings, thread model.Thread, myAccountID string) error {
	return renderTemplateToPager(settings, "thread", myAccountID, thread)
}

// PrintList prints the details of the list to the pager.
func PrintList(settings Settings, list model.List) error {
	return renderTemplateToPager(settings, "list", "", list)
}

// PrintLists prints the set of lists to the pager.
func PrintLists(settings Settings, lists []model.List) error {
	return renderTemplateToPager(settings, "listOflist", "", lists)
}

// PrintNotification prints the details of the notification to the pager.
func PrintNotification(settings Settings, notification model.Notification, myAccountID string) error {
	return renderTemplateToPager(settings, "notificationDoc", myAccountID, notification)
}

// PrintNotificationList prints the list of notifications to the pager.
func PrintNotificationList(settings Settings, list []model.Notification, myAccountID string) error {
	return renderTemplateToPager(settings, "notificationList", myAccountID, list)
}

// PrintTokenList prints the list of tokens to the pager.
func PrintTokenList(settings Settings, list model.TokenList) error {
	return renderTemplateToPager(settings, "tokenList", "", list)
}

// PrintToken prints the details of the token.
func PrintToken(settings Settings, token model.Token) error {
	return renderTemplateToPager(settings, "tokenDoc", "", token)
}

// PrintAliases prints the user's list of aliases.
func PrintAliases(settings Settings, aliases map[string]string) error {
	return renderTemplateToPager(settings, "aliases", "", aliases)
}

// PrintFilters prints the user's list of filters.
func PrintFilters(settings Settings, filters []model.FilterV2) error {
	return renderTemplateToPager(settings, "filterList", "", filters)
}

// PrintFilter prints the details of a filter.
func PrintFilter(settings Settings, filter model.FilterV2) error {
	return renderTemplateToPager(settings, "filter", "", filter)
}

// PrintFilterKeyword prints the details of a filter-keyword.
func PrintFilterKeyword(settings Settings, filterKeyword model.FilterKeyword) error {
	return renderTemplateToPager(settings, "filter-keyword", "", filterKeyword)
}

// PrintFilterStatus prints the details of a filter-status.
func PrintFilterStatus(settings Settings, filterStatus model.FilterStatus) error {
	return renderTemplateToPager(settings, "filter-status", "", filterStatus)
}

func renderTemplateToPager(settings Settings, templateName, myAccountID string, data any) error {
	if settings.pager == "" {
		return renderTemplateToStdout(
			settings,
			templateName,
			myAccountID,
			data,
		)
	}

	cmdSplit := strings.Split(settings.pager, " ")

	binary, err := exec.LookPath(cmdSplit[0])
	if err != nil {
		return fmt.Errorf("unable to perform the lookup for %q, %w", cmdSplit[0], err)
	}

	var pagerCmd *exec.Cmd

	if len(cmdSplit) == 1 {
		pagerCmd = exec.Command(binary) // #nosec G204 -- External command call defined in user's configuration file.
	} else {
		pagerCmd = exec.Command(binary, cmdSplit[1:]...) // #nosec G204 -- External command call defined in user's configuration file.
	}

	// Render the template to the pipe.
	reader, writer, err := os.Pipe()
	if err != nil {
		return fmt.Errorf("error creating the pipe: %w", err)
	}

	defer reader.Close()

	if err := renderTemplate(
		writer,
		settings,
		templateName,
		myAccountID,
		data,
	); err != nil {
		return fmt.Errorf("error rendering the template: %w", err)
	}

	_ = writer.Close()

	// Pipe the text data to the pager.
	pagerCmd.Stdin = reader
	pagerCmd.Stdout = os.Stdout

	_ = pagerCmd.Run()

	return nil
}

func renderTemplateToStdout(
	settings Settings,
	templateName string,
	myAccountID string,
	data any,
) error {
	return renderTemplate(
		os.Stdout,
		settings,
		templateName,
		myAccountID,
		data,
	)
}

func renderTemplate(
	writer io.Writer,
	settings Settings,
	templateName string,
	myAccountID string,
	data any,
) error {
	tmpl, err := template.New("").
		Funcs(funcMap(settings, myAccountID)).
		ParseFS(templatesFS, "templates/*")
	if err != nil {
		return fmt.Errorf("error parsing the templates: %w", err)
	}

	if err := tmpl.ExecuteTemplate(writer, templateName, data); err != nil {
		return fmt.Errorf("error executing the %q template: %w", templateName, err)
	}

	return nil
}

func printToStdout(text string) {
	_, _ = os.Stdout.WriteString(text)
}

func printToStderr(text string) {
	_, _ = os.Stderr.WriteString(text)
}
