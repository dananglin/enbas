package utilities

import (
	"strings"
	"unicode"
)

func Titled(str string) string {
	if str == "" {
		return str
	}

	runes := []rune(str)

	runes[0] = unicode.ToUpper(runes[0])

	return string(runes)
}

func AllCaps(str string) string {
	return strings.ToUpper(str)
}

func SnakeToCamel(value string, pascal bool) string {
	var builder strings.Builder

	runes := []rune(value)
	numRunes := len(runes)
	cursor := 0

	for cursor < numRunes {
		switch {
		case cursor == 0 && pascal:
			builder.WriteRune(unicode.ToUpper(runes[cursor]))

			cursor++
		case runes[cursor] != '-':
			builder.WriteRune(runes[cursor])

			cursor++
		case cursor != numRunes-1 && unicode.IsLower(runes[cursor+1]):
			builder.WriteRune(unicode.ToUpper(runes[cursor+1]))

			cursor += 2
		default:
			cursor++
		}
	}

	return builder.String()
}

func Join(lines []string) string {
	return strings.Join(lines, " ")
}
