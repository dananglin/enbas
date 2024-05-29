package utilities

import (
	"regexp"
	"time"
)

const (
	reset       = "\033[0m"
	boldblue    = "\033[34;1m"
	boldmagenta = "\033[35;1m"
	green       = "\033[32m"
)

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
