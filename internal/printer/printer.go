package printer

import (
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"
)

const (
	minTerminalWidth    = 40
	noMediaDescription  = "This media attachment has no description."
	symbolBullet        = "\u2022"
	symbolPollMeter     = "\u2501"
	symbolCheckMark     = "\u2714"
	symbolFailure       = "\u2717"
	symbolImage         = "\uf03e"
	symbolLiked         = "\uf51f"
	symbolNotLiked      = "\uf41e"
	symbolBookmarked    = "\uf47a"
	symbolNotBookmarked = "\uf461"
	symbolBoosted       = "\u2BAD"
	dateFormat          = "02 Jan 2006"
	dateTimeFormat      = "02 Jan 2006, 15:04 (MST)"
)

type theme struct {
	reset       string
	bold        string
	boldblue    string
	boldmagenta string
	green       string
	boldgreen   string
	grey        string
	red         string
	boldred     string
	yellow      string
	boldyellow  string
}

type Printer struct {
	theme                  theme
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
	theme := theme{
		reset:       "\033[0m",
		bold:        "\033[1m",
		boldblue:    "\033[34;1m",
		boldmagenta: "\033[35;1m",
		green:       "\033[32m",
		boldgreen:   "\033[32;1m",
		grey:        "\033[90m",
		red:         "\033[31m",
		boldred:     "\033[31;1m",
		yellow:      "\033[33m",
		boldyellow:  "\033[33;1m",
	}

	if lineWrapCharacterLimit < minTerminalWidth {
		lineWrapCharacterLimit = minTerminalWidth
	}

	return &Printer{
		theme:                  theme,
		noColor:                noColor,
		lineWrapCharacterLimit: lineWrapCharacterLimit,
		pager:                  pager,
		statusSeparator:        strings.Repeat("\u2501", lineWrapCharacterLimit),
	}
}

func (p Printer) PrintSuccess(text string) {
	success := p.theme.boldgreen + symbolCheckMark + p.theme.reset
	if p.noColor {
		success = symbolCheckMark
	}

	printToStdout(success + " " + text + "\n")
}

func (p Printer) PrintFailure(text string) {
	failure := p.theme.boldred + symbolFailure + p.theme.reset
	if p.noColor {
		failure = symbolFailure
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

func (p Printer) bold(text string) string {
	if p.noColor {
		return text
	}

	return p.theme.bold + text + p.theme.reset
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

func (p Printer) print(text string) {
	if p.pager == "" {
		printToStdout(text)

		return
	}

	cmdSplit := strings.Split(p.pager, " ")

	binary, err := exec.LookPath(cmdSplit[0])
	if err != nil {
		printToStdout(text)

		return
	}

	var pager *exec.Cmd

	if len(cmdSplit) == 1 {
		pager = exec.Command(binary) // #nosec G204 -- External command call defined in user's configuration file.
	} else {
		pager = exec.Command(binary, cmdSplit[1:]...) // #nosec G204 -- External command call defined in user's configuration file.
	}

	// Write the text data to the pipe.
	reader, writer, err := os.Pipe()
	if err != nil {
		printToStdout(text)

		return
	}
	defer reader.Close()

	if _, err = writer.WriteString(text); err != nil {
		printToStdout(text)

		return
	}

	if err := writer.Close(); err != nil {
		printToStdout(text)

		return
	}

	// Pipe the text data to the pager.
	pager.Stdin = reader
	pager.Stdout = os.Stdout

	_ = pager.Run()
}

func printToStdout(text string) {
	_, _ = os.Stdout.WriteString(text)
}

func printToStderr(text string) {
	_, _ = os.Stderr.WriteString(text)
}

func formatDate(date time.Time) string {
	return date.Local().Format(dateFormat) //nolint:gosmopolitan
}

func formatDateTime(date time.Time) string {
	return date.Local().Format(dateTimeFormat) //nolint:gosmopolitan
}
