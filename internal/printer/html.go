package printer

import (
	"io"
	"strconv"
	"strings"

	"golang.org/x/net/html"
)

const (
	htmlNoList int = iota
	htmlOrderedList
	htmlUnorderedList
)

type htmlConvertState struct {
	htmlListType     int
	orderedListIndex int
}

func convertHTMLToText(text string) string {
	var builder strings.Builder

	state := htmlConvertState{
		htmlListType:     htmlNoList,
		orderedListIndex: 1,
	}

	token := html.NewTokenizer(strings.NewReader(text))

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
			processTagToken(&state, &builder, tag)
		}
	}
}

func processTagToken(state *htmlConvertState, writer io.StringWriter, tag string) {
	switch tag {
	case "<br>", "<p>", "</p>", "</li>":
		_, _ = writer.WriteString("\n")
	case "<ul>":
		state.htmlListType = htmlUnorderedList
		_, _ = writer.WriteString("\n")
	case "<ol>":
		state.htmlListType = htmlOrderedList
		_, _ = writer.WriteString("\n")
	case "</ul>":
		state.htmlListType = htmlNoList
	case "</ol>":
		state.htmlListType = htmlNoList
		state.orderedListIndex = 1
	case "<li>":
		switch state.htmlListType {
		case htmlUnorderedList:
			_, _ = writer.WriteString(symbolBullet + " ")
		case htmlOrderedList:
			_, _ = writer.WriteString(strconv.Itoa(state.orderedListIndex) + ". ")
			state.orderedListIndex++
		}
	}
}
