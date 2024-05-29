package utilities

import (
	"strings"

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
			text := token.Token().Data
			builder.WriteString(text)
		case html.StartTagToken, html.EndTagToken:
			tag := token.Token().String()
			builder.WriteString(transformTag(tag))
		}
	}
}

func transformTag(tag string) string {
	switch tag {
	case "<br>":
		return "\n"
	case "<p>", "</p>":
		return "\n"
	}

	return ""
}
