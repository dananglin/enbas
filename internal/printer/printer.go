package printer

import (
	"embed"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"text/template"
	"time"

	"codeflow.dananglin.me.uk/apollo/enbas/internal/model"
)

//go:embed templates/*
var templatesFS embed.FS

type Printer struct {
	noColor                bool
	lineWrapCharacterLimit int
	pager                  string
	statusSeparator        string
}

func NewPrinter(
	noColor bool,
	pager string,
	lineWrapCharacterLimit int,
) *Printer {
	if lineWrapCharacterLimit < minTerminalWidth {
		lineWrapCharacterLimit = minTerminalWidth
	}

	return &Printer{
		noColor:                noColor,
		lineWrapCharacterLimit: lineWrapCharacterLimit,
		pager:                  pager,
		statusSeparator:        strings.Repeat("\u2501", lineWrapCharacterLimit),
	}
}

func (p Printer) PrintSuccess(text string) {
	success := boldgreen + symbolCheckMark + reset
	if p.noColor {
		success = symbolCheckMark
	}

	printToStdout(success + " " + text + "\n")
}

func (p Printer) PrintFailure(text string) {
	failure := boldred + symbolFailure + reset
	if p.noColor {
		failure = symbolFailure
	}

	printToStderr(failure + " " + text + "\n")
}

func (p Printer) PrintInfo(text string) {
	printToStdout(text)
}

func (p Printer) PrintAccount(
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

	return p.renderTemplateToPager("account", myAccountID, data)
}

func (p Printer) PrintAccountList(list model.AccountList) error {
	if list.BlockedAccounts {
		return p.renderTemplateToPager("blockedAccounts", "", list)
	}

	return p.renderTemplateToPager("accountList", "", list)
}

func (p Printer) PrintStatus(
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

	return p.renderTemplateToPager("statusDoc", myAccountID, data)
}

// PrintStatusList prints a drawn list of statuses.
func (p Printer) PrintStatusList(list model.StatusList, myAccountID string) error {
	return p.renderTemplateToPager("statusList", myAccountID, list)
}

func (p Printer) PrintInstance(instance model.InstanceV2) error {
	return p.renderTemplateToPager("instance", "", instance)
}

func (p Printer) PrintMediaAttachment(attachement model.MediaAttachment) error {
	return p.renderTemplateToPager("mediaAttachmentDoc", "", attachement)
}

func (p Printer) PrintTag(tag model.Tag) error {
	return p.renderTemplateToPager("tag", "", tag)
}

func (p Printer) PrintTagList(list model.TagList) error {
	return p.renderTemplateToPager("tagList", "", list)
}

func (p Printer) PrintThread(thread model.Thread, myAccountID string) error {
	return p.renderTemplateToPager("thread", myAccountID, thread)
}

func (p Printer) PrintList(list model.List) error {
	return p.renderTemplateToPager("list", "", list)
}

func (p Printer) PrintLists(lists []model.List) error {
	return p.renderTemplateToPager("listOflist", "", lists)
}

func (p Printer) headerFormat(text string) string {
	if p.noColor {
		return text
	}

	return boldblue + text + reset
}

func (p Printer) fieldFormat(text string) string {
	if p.noColor {
		return text + ":"
	}

	return green + text + reset + ":"
}

func (p Printer) boldFormat(text string) string {
	if p.noColor {
		return text
	}

	return bold + text + reset
}

func (p Printer) drawBoostSymbol(boosted bool) string {
	if boosted && !p.noColor {
		return boldyellow + symbolBoosted + reset
	}

	return symbolBoosted
}

func (p Printer) drawLikeSymbol(liked bool) string {
	if liked && !p.noColor {
		return boldyellow + symbolLiked + reset
	}

	return symbolNotLiked
}

func (p Printer) drawBookmarkSymbol(bookmarked bool) string {
	if bookmarked && !p.noColor {
		return boldyellow + symbolBookmarked + reset
	}

	return symbolNotBookmarked
}

func (p Printer) fullDisplayNameFormat(displayName, acct string) string {
	// use this pattern to remove all emoji strings
	pattern := regexp.MustCompile(`\s:[A-Za-z0-9_]*:`)

	var builder strings.Builder

	if p.noColor {
		builder.WriteString(pattern.ReplaceAllString(displayName, ""))
	} else {
		builder.WriteString(boldmagenta + pattern.ReplaceAllString(displayName, "") + reset)
	}

	builder.WriteString(" (@" + acct + ")")

	return builder.String()
}

func (p Printer) drawStatusCardSeparator() string {
	return p.statusSeparator
}

func (p Printer) renderTemplateToPager(templateName, myAccountID string, data any) error {
	cmdSplit := strings.Split(p.pager, " ")

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

	funcMap := template.FuncMap{
		"convertHTMLToText":       convertHTMLToText,
		"formatDate":              formatDate,
		"formatDateTime":          formatDateTime,
		"headerFormat":            p.headerFormat,
		"fieldFormat":             p.fieldFormat,
		"fullDisplayNameFormat":   p.fullDisplayNameFormat,
		"boldFormat":              p.boldFormat,
		"drawStatusCardSeparator": p.drawStatusCardSeparator,
		"drawBoostSymbol":         p.drawBoostSymbol,
		"drawLikeSymbol":          p.drawLikeSymbol,
		"drawBookmarkSymbol":      p.drawBookmarkSymbol,
		"wrapLines":               p.wrapLines,
		"showPollResults":         showPollResults(myAccountID),
		"getPollOptionDetails":    p.getPollOptionDetails,
	}

	tmpl, err := template.New("").Funcs(funcMap).ParseFS(templatesFS, "templates/*")
	if err != nil {
		return fmt.Errorf("error parsing the templates: %w", err)
	}

	if err := tmpl.ExecuteTemplate(writer, templateName, data); err != nil {
		return fmt.Errorf("error executing the %q template: %w", templateName, err)
	}

	_ = writer.Close()

	// Pipe the text data to the pager.
	pagerCmd.Stdin = reader
	pagerCmd.Stdout = os.Stdout

	_ = pagerCmd.Run()

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
