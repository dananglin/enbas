package utilities

import (
	"strings"
	"unicode"

	"golang.org/x/net/html"
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
