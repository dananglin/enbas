package utilities

import (
	"strings"
	"unicode"
)

func WrapLines(text, separator string, charLimit int) string {
	lines := strings.Split(text, "\n")

	if len(lines) == 1 {
		return wrapLine(lines[0], separator, charLimit)
	}

	var builder strings.Builder

	for i, line := range lines {
		builder.WriteString(wrapLine(line, separator, charLimit))

		if i < len(lines)-1 {
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
		rightcursor += charLimit

		for !unicode.IsSpace(rune(line[rightcursor-1])) && (rightcursor > leftcursor) {
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