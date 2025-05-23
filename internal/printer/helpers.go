package printer

import (
	"regexp"
	"strings"
	"text/template"
	"time"

	"codeflow.dananglin.me.uk/apollo/enbas/internal/model"
)

func funcMap(settings Settings, myAccountID string) template.FuncMap {
	return template.FuncMap{
		"convertHTMLToText":     convertHTMLToText,
		"formatDate":            formatDate,
		"formatDateTime":        formatDateTime,
		"headerFormat":          headerFormat(settings.noColor),
		"fieldFormat":           fieldFormat(settings.noColor),
		"fullDisplayNameFormat": fullDisplayNameFormat(settings.noColor),
		"boldFormat":            boldFormat(settings.noColor),
		"drawCardSeparator":     drawCardSeparator(settings.lineWrapCharacterLimit),
		"drawBoostSymbol":       drawBoostSymbol(settings.noColor),
		"drawLikeSymbol":        drawLikeSymbol(settings.noColor),
		"drawBookmarkSymbol":    drawBookmarkSymbol(settings.noColor),
		"wrapLines":             wrapLines(settings.lineWrapCharacterLimit),
		"showPollResults":       showPollResults(myAccountID),
		"getPollOptionDetails":  getPollOptionDetails(settings.noColor, settings.lineWrapCharacterLimit),
		"notificationSummary":   notificationSummary,
		"statusFilterAction":    statusFilterAction,
		"statusFilteredTitle":   statusFilteredTitle(settings.noColor),
	}
}

func formatDate(date time.Time) string {
	return date.Local().Format("02 Jan 2006") //nolint:gosmopolitan
}

func formatDateTime(date time.Time) string {
	return date.Local().Format("02 Jan 2006, 15:04 (MST)") //nolint:gosmopolitan
}

func showPollResults(myAccountID string) func(string, bool, bool) bool {
	return func(statusOwnerID string, expired, voted bool) bool {
		return (myAccountID == statusOwnerID) || expired || voted
	}
}

func headerFormat(noColor bool) func(string) string {
	return func(text string) string {
		if noColor {
			return text
		}

		return boldblue + text + reset
	}
}

func fieldFormat(noColor bool) func(string) string {
	return func(text string) string {
		if noColor {
			return text + ":"
		}

		return green + text + reset + ":"
	}
}

func boldFormat(noColor bool) func(string) string {
	return func(text string) string {
		if noColor {
			return text
		}

		return bold + text + reset
	}
}

func statusFilteredTitle(noColor bool) func() string {
	return func() string {
		if noColor {
			return "Status Filtered"
		}

		return boldyellow + "Status Filtered" + reset
	}
}

func fullDisplayNameFormat(noColor bool) func(string, string) string {
	return func(displayName, acct string) string {
		// use this pattern to remove all emoji strings
		pattern := regexp.MustCompile(`\s:[A-Za-z0-9_]*:`)

		var builder strings.Builder

		if noColor {
			builder.WriteString(pattern.ReplaceAllString(displayName, ""))
		} else {
			builder.WriteString(boldmagenta + pattern.ReplaceAllString(displayName, "") + reset)
		}

		builder.WriteString(" (@" + acct + ")")

		return builder.String()
	}
}

func drawCardSeparator(charLimit int) func() string {
	separator := strings.Repeat("\u2501", charLimit)

	return func() string {
		return separator
	}
}

func drawBoostSymbol(noColor bool) func(bool) string {
	return func(boosted bool) string {
		if boosted && !noColor {
			return boldyellow + "\u2BAD" + reset
		}

		return "\u2BAD"
	}
}

func drawLikeSymbol(noColor bool) func(bool) string {
	return func(liked bool) string {
		if liked && !noColor {
			return boldyellow + "\uF51F" + reset
		} else if liked && noColor {
			return "\uF51F"
		}

		return "\uF41E"
	}
}

func drawBookmarkSymbol(noColor bool) func(bool) string {
	return func(bookmarked bool) string {
		if bookmarked && !noColor {
			return boldyellow + "\uF47A" + reset
		} else if bookmarked && noColor {
			return "\uF47A"
		}

		return "\uF461"
	}
}

type notificationSummaryDetails struct {
	Header  string
	Details string
}

func notificationSummary(notificationType string, fullDisplayName string) notificationSummaryDetails {
	switch notificationType {
	case "follow":
		return notificationSummaryDetails{
			Header:  "SOMEONE FOLLOWED YOU:",
			Details: fullDisplayName + " followed you.",
		}
	case "follow_request":
		return notificationSummaryDetails{
			Header:  "YOU'VE RECEIVED A FOLLOW REQUEST:",
			Details: fullDisplayName + " sent you a follow request.",
		}
	case "mention":
		return notificationSummaryDetails{
			Header:  "SOMEONE MENTIONED YOU IN A STATUS:",
			Details: fullDisplayName + " mentioned you in the below status.",
		}
	case "reblog":
		return notificationSummaryDetails{
			Header:  "SOMEONE BOOSTED YOUR STATUS:",
			Details: fullDisplayName + " boosted your status.",
		}
	case "favourite":
		return notificationSummaryDetails{
			Header:  "SOMEONE LIKED YOUR STATUS:",
			Details: fullDisplayName + " liked your status.",
		}
	case "poll":
		return notificationSummaryDetails{
			Header:  "POLL CLOSED:",
			Details: "The poll below has closed.",
		}
	case "status":
		return notificationSummaryDetails{
			Header:  "SOMEONE POSTED A STATUS:",
			Details: fullDisplayName + " posted the status below.",
		}
	default:
		return notificationSummaryDetails{
			Header:  "UNKNOWN NOTIFICATION TYPE:",
			Details: "Received a notification of an unknown type.",
		}
	}
}

func statusFilterAction(filterResults []model.FilterResult) string {
	if len(filterResults) == 0 {
		return ""
	}

	action := ""

	for idx := range filterResults {
		switch filterResults[idx].Filter.Action {
		case model.FilterActionHide:
			return model.FilterActionHide
		case model.FilterActionWarn:
			action = model.FilterActionWarn
		}
	}

	return action
}
