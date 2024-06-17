// SPDX-FileCopyrightText: 2024 Dan Anglin <d.n.i.anglin@gmail.com>
//
// SPDX-License-Identifier: GPL-3.0-or-later

package printer

import (
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"
)

const (
	minTerminalWidth = 40
)

type theme struct {
	reset       string
	boldblue    string
	boldmagenta string
	green       string
	boldgreen   string
	grey        string
	red         string
	boldred     string
}

type Printer struct {
	theme            theme
	noColor          bool
	maxTerminalWidth int
	pager            string
	statusSeparator  string
	bullet           string
	pollMeterSymbol  string
	successSymbol    string
	failureSymbol    string
	dateFormat       string
	dateTimeFormat   string
}

func NewPrinter(
	noColor bool,
	pager string,
	maxTerminalWidth int,
) *Printer {
	theme := theme{
		reset:       "\033[0m",
		boldblue:    "\033[34;1m",
		boldmagenta: "\033[35;1m",
		green:       "\033[32m",
		boldgreen:   "\033[32;1m",
		grey:        "\033[90m",
		red:         "\033[31m",
		boldred:     "\033[31;1m",
	}

	if maxTerminalWidth < minTerminalWidth {
		maxTerminalWidth = minTerminalWidth
	}

	return &Printer{
		noColor:          noColor,
		maxTerminalWidth: maxTerminalWidth,
		pager:            pager,
		statusSeparator:  strings.Repeat("\u2501", maxTerminalWidth),
		bullet:           "\u2022",
		pollMeterSymbol:  "\u2501",
		successSymbol:    "\u2714",
		failureSymbol:    "\u2717",
		dateFormat:       "02 Jan 2006",
		dateTimeFormat:   "02 Jan 2006, 15:04 (MST)",
		theme:            theme,
	}
}

func (p Printer) PrintSuccess(text string) {
	success := p.theme.boldgreen + p.successSymbol + p.theme.reset
	if p.noColor {
		success = p.successSymbol
	}

	printToStdout(success + " " + text + "\n")
}

func (p Printer) PrintFailure(text string) {
	failure := p.theme.boldred + p.failureSymbol + p.theme.reset
	if p.noColor {
		failure = p.failureSymbol
	}

	printToStderr(failure + " " + text + "\n")
}

func (p Printer) PrintInfo(text string) {
	printToStdout(text)
}

func (p Printer) headerFormat(text string) string {
	if p.noColor {
		return text
	}

	return p.theme.boldblue + text + p.theme.reset
}

func (p Printer) fieldFormat(text string) string {
	if p.noColor {
		return text
	}

	return p.theme.green + text + p.theme.reset
}

func (p Printer) fullDisplayNameFormat(displayName, acct string) string {
	// use this pattern to remove all emoji strings
	pattern := regexp.MustCompile(`\s:[A-Za-z0-9_]*:`)

	var builder strings.Builder

	if p.noColor {
		builder.WriteString(pattern.ReplaceAllString(displayName, ""))
	} else {
		builder.WriteString(p.theme.boldmagenta + pattern.ReplaceAllString(displayName, "") + p.theme.reset)
	}

	builder.WriteString(" (@" + acct + ")")

	return builder.String()
}

func (p Printer) formatDate(date time.Time) string {
	return date.Local().Format(p.dateFormat)
}

func (p Printer) formatDateTime(date time.Time) string {
	return date.Local().Format(p.dateTimeFormat)
}

func (p Printer) print(text string) {
	if p.pager == "" {
		printToStdout(text)

		return
	}

	cmdSplit := strings.Split(p.pager, " ")

	pager := new(exec.Cmd)

	if len(cmdSplit) == 1 {
		pager = exec.Command(cmdSplit[0]) //nolint:gosec
	} else {
		pager = exec.Command(cmdSplit[0], cmdSplit[1:]...) //nolint:gosec
	}

	pipe, err := pager.StdinPipe()
	if err != nil {
		printToStdout(text)

		return
	}

	pager.Stdout = os.Stdout
	pager.Stderr = os.Stderr

	_ = pager.Start()

	defer func() {
		_ = pipe.Close()
		_ = pager.Wait()
	}()

	_, _ = pipe.Write([]byte(text))
}

func printToStdout(text string) {
	os.Stdout.WriteString(text)
}

func printToStderr(text string) {
	os.Stderr.WriteString(text)
}
