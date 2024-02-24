package utilities

import (
	"regexp"
	"strings"
	"time"
	"unicode"

	"golang.org/x/net/html"
)

const (
	reset    = "\033[0m"
	boldblue = "\033[34;1m"
	boldmagenta = "\033[35;1m"
	green    = "\033[32m"
)

func StripHTMLTags(text string) string {
	token := html.NewTokenizer(strings.NewReader(text))

	var builder strings.Builder

	for {
		tt := token.Next()
		switch tt {
		case html.ErrorToken:
			return builder.String()
		case html.TextToken:
			builder.WriteString(token.Token().Data + " ")
		}
	}
}

func WrapLine(line, separator string, charLimit int) string {
	if len(line) <= charLimit {
		return line
	}

	leftcursor, rightcursor := 0, 0

	var builder strings.Builder

	for rightcursor < (len(line) - charLimit) {
		rightcursor += charLimit
		for !unicode.IsSpace(rune(line[rightcursor-1])) {
			rightcursor--
		}
		builder.WriteString(line[leftcursor:rightcursor] + separator)
		leftcursor = rightcursor
	}

	builder.WriteString(line[rightcursor:])

	return builder.String()
}

func HeaderFormat(text string) string {
	return boldblue + text + reset
}

func FieldFormat(text string) string {
	return green + text + reset
}

func DisplayNameFormat(text string) string {
	// use this pattern to remove all emoji strings
	pattern := regexp.MustCompile(`\s:[A-Za-z0-9]*:`)

	return boldmagenta + pattern.ReplaceAllString(text, "") + reset
}

func FormatDate(date time.Time) string {
	return date.Local().Format("02 Jan 2006")
}

func FormatTime(date time.Time) string {
	return date.Local().Format("02 Jan 2006, 15:04 (MST)")
}
