// SPDX-FileCopyrightText: 2024 Dan Anglin <d.n.i.anglin@gmail.com>
//
// SPDX-License-Identifier: GPL-3.0-or-later

package printer

import (
	"regexp"
	"strings"
	"unicode"
)

type extraIndentConditiion struct {
	pattern *regexp.Regexp
	indent  string
}

func (p Printer) wrapLines(text string, nIndent int) string {
	if nIndent >= p.lineWrapCharacterLimit {
		nIndent = 0
	}

	separator := "\n" + strings.Repeat(" ", nIndent)

	lines := strings.Split(text, "\n")

	if len(lines) == 1 {
		return wrapLine(lines[0], separator, p.lineWrapCharacterLimit-nIndent)
	}

	var builder strings.Builder

	extraIndentConditions := []extraIndentConditiion{
		{
			pattern: regexp.MustCompile(`^[-*` + symbolBullet + `]\s.*$`),
			indent:  "  ",
		},
		{
			pattern: regexp.MustCompile(`^[0-9]{1}\.\s.*$`),
			indent:  "   ",
		},
		{
			pattern: regexp.MustCompile(`^[0-9]{2}\.\s.*$`),
			indent:  "    ",
		},
	}

	for ind, line := range lines {
		builder.WriteString(wrapLine(line, separator+extraIndent(line, extraIndentConditions), p.lineWrapCharacterLimit-nIndent))

		if ind < len(lines)-1 {
			builder.WriteString(separator)
		}
	}

	return builder.String()
}

func wrapLine(line, separator string, charLimit int) string {
	if len(line) <= charLimit {
		return line
	}

	leftcursor, rightcursor := 0, 0

	var builder strings.Builder

	for rightcursor < (len(line) - charLimit) {
		rightcursor += (charLimit - 1)

		for (rightcursor > leftcursor) && !unicode.IsSpace(rune(line[rightcursor-1])) {
			rightcursor--
		}

		if rightcursor == leftcursor {
			rightcursor = leftcursor + charLimit
		}

		builder.WriteString(line[leftcursor:rightcursor] + separator)
		leftcursor = rightcursor
	}

	builder.WriteString(line[rightcursor:])

	return builder.String()
}

func extraIndent(line string, conditions []extraIndentConditiion) string {
	for ind := range conditions {
		if conditions[ind].pattern.MatchString(line) {
			return conditions[ind].indent
		}
	}

	return ""
}
