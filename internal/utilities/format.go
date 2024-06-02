// SPDX-FileCopyrightText: 2024 Dan Anglin <d.n.i.anglin@gmail.com>
//
// SPDX-License-Identifier: GPL-3.0-or-later

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

func HeaderFormat(noColor bool, text string) string {
	if noColor {
		return text
	}

	return boldblue + text + reset
}

func FieldFormat(noColor bool, text string) string {
	if noColor {
		return text
	}

	return green + text + reset
}

func DisplayNameFormat(noColor bool, text string) string {
	// use this pattern to remove all emoji strings
	pattern := regexp.MustCompile(`\s:[A-Za-z0-9]*:`)

	if noColor {
		return pattern.ReplaceAllString(text, "")
	}

	return boldmagenta + pattern.ReplaceAllString(text, "") + reset
}

func FormatDate(date time.Time) string {
	return date.Local().Format("02 Jan 2006")
}

func FormatTime(date time.Time) string {
	return date.Local().Format("02 Jan 2006, 15:04 (MST)")
}
