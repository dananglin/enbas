package printer

import (
	"embed"
	"fmt"
	"io"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"text/template"
	"time"

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

	success := boldgreen + icon + reset
	if settings.noColor {
		success = icon
	}

	printToStdout(success + " " + text + "\n")
}

// PrintFailure prints the failure message to standard error.
func PrintFailure(settings Settings, text string) {
	const icon = "\u2717"

	failure := boldred + icon + reset
	if settings.noColor {
		failure = icon
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
		GoVersion:     info.GoVersion,
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
	funcMap := template.FuncMap{
		"convertHTMLToText":     convertHTMLToText,
		"formatDate":            formatDate,
		"formatDateTime":        formatDateTime,
		"headerFormat":          headerFormat(settings.noColor),
		"fieldFormat":           fieldFormat(settings.noColor),
		"fullDisplayNameFormat": fullDisplayNameFormat(settings.noColor),
		"boldFormat":            boldFormat(settings.noColor),
		"drawCardSeparator":     drawCardSeparator(settings.lineWrapCharacterLimit),
		"drawBoostSymbol":       drawBoostSymbol(settings.noColor),
		"drawLikeSymbol":        drawLikeSymbol(settings.noColor),
		"drawBookmarkSymbol":    drawBookmarkSymbol(settings.noColor),
		"wrapLines":             wrapLines(settings.lineWrapCharacterLimit),
		"showPollResults":       showPollResults(myAccountID),
		"getPollOptionDetails":  getPollOptionDetails(settings.noColor, settings.lineWrapCharacterLimit),
		"notificationSummary":   notificationSummary,
	}

	tmpl, err := template.New("").Funcs(funcMap).ParseFS(templatesFS, "templates/*")
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

func formatDate(date time.Time) string {
	return date.Local().Format("02 Jan 2006") //nolint:gosmopolitan
}

func formatDateTime(date time.Time) string {
	return date.Local().Format("02 Jan 2006, 15:04 (MST)") //nolint:gosmopolitan
}

func showPollResults(myAccountID string) func(string, bool, bool) bool {
	return func(statusOwnerID string, expired, voted bool) bool {
		return (myAccountID == statusOwnerID) || expired || voted
	}
}

func headerFormat(noColor bool) func(string) string {
	return func(text string) string {
		if noColor {
			return text
		}

		return boldblue + text + reset
	}
}

func fieldFormat(noColor bool) func(string) string {
	return func(text string) string {
		if noColor {
			return text + ":"
		}

		return green + text + reset + ":"
	}
}

func boldFormat(noColor bool) func(string) string {
	return func(text string) string {
		if noColor {
			return text
		}

		return bold + text + reset
	}
}

func fullDisplayNameFormat(noColor bool) func(string, string) string {
	return func(displayName, acct string) string {
		// use this pattern to remove all emoji strings
		pattern := regexp.MustCompile(`\s:[A-Za-z0-9_]*:`)

		var builder strings.Builder

		if noColor {
			builder.WriteString(pattern.ReplaceAllString(displayName, ""))
		} else {
			builder.WriteString(boldmagenta + pattern.ReplaceAllString(displayName, "") + reset)
		}

		builder.WriteString(" (@" + acct + ")")

		return builder.String()
	}
}

func drawCardSeparator(charLimit int) func() string {
	separator := strings.Repeat("\u2501", charLimit)

	return func() string {
		return separator
	}
}

func drawBoostSymbol(noColor bool) func(bool) string {
	return func(boosted bool) string {
		if boosted && !noColor {
			return boldyellow + "\u2BAD" + reset
		}

		return "\u2BAD"
	}
}

func drawLikeSymbol(noColor bool) func(bool) string {
	return func(liked bool) string {
		if liked && !noColor {
			return boldyellow + "\uf51f" + reset
		} else if liked && noColor {
			return "\uf51f"
		}

		return "\uf41e"
	}
}

func drawBookmarkSymbol(noColor bool) func(bool) string {
	return func(bookmarked bool) string {
		if bookmarked && !noColor {
			return boldyellow + "\uf47a" + reset
		} else if bookmarked && noColor {
			return "\uf47a"
		}

		return "\uf461"
	}
}

type notificationSummaryDetails struct {
	Header  string
	Details string
}

func notificationSummary(notificationType model.NotificationType, fullDisplayName string) notificationSummaryDetails {
	switch notificationType {
	case model.NotificationTypeFollow:
		return notificationSummaryDetails{
			Header:  "SOMEONE FOLLOWED YOU:",
			Details: fullDisplayName + " followed you.",
		}
	case model.NotificationTypeFollowRequest:
		return notificationSummaryDetails{
			Header:  "YOU'VE RECEIVED A FOLLOW REQUEST:",
			Details: fullDisplayName + " sent you a follow request.",
		}
	case model.NotificationTypeMention:
		return notificationSummaryDetails{
			Header:  "SOMEONE MENTIONED YOU IN A STATUS:",
			Details: fullDisplayName + " mentioned you in the below status.",
		}
	case model.NotificationTypeReblog:
		return notificationSummaryDetails{
			Header:  "SOMEONE BOOSTED YOUR STATUS:",
			Details: fullDisplayName + " boosted your status.",
		}
	case model.NotificationTypeFavourite:
		return notificationSummaryDetails{
			Header:  "SOMEONE LIKED YOUR STATUS:",
			Details: fullDisplayName + " liked your status.",
		}
	case model.NotificationTypePoll:
		return notificationSummaryDetails{
			Header:  "POLL CLOSED:",
			Details: "The poll below has closed.",
		}
	case model.NotificationTypeStatus:
		return notificationSummaryDetails{
			Header:  "SOMEONE POSTED A STATUS:",
			Details: fullDisplayName + " posted the status below.",
		}
	default:
		return notificationSummaryDetails{
			Header:  "UNKNOWN NOTIFICATION TYPE:",
			Details: "Received a notification of an unknown type.",
		}
	}
}
