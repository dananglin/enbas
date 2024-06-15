package model

import (
	"io"
	"math"
	"strconv"
	"strings"
	"time"

	"codeflow.dananglin.me.uk/apollo/enbas/internal/utilities"
)

type Poll struct {
	Emojis      []Emoji      `json:"emojis"`
	Expired     bool         `json:"expired"`
	Voted       bool         `json:"voted"`
	Multiple    bool         `json:"multiple"`
	ExpiredAt   time.Time    `json:"expires_at"`
	ID          string       `json:"id"`
	OwnVotes    []int        `json:"own_votes"`
	VotersCount int          `json:"voters_count"`
	VotesCount  int          `json:"votes_count"`
	Options     []PollOption `json:"options"`
}

type PollOption struct {
	Title      string `json:"title"`
	VotesCount int    `json:"votes_count"`
}

func (p Poll) Display(noColor bool) string {
	var builder strings.Builder

	indent := "  "

	builder.WriteString(
		utilities.HeaderFormat(noColor, "POLL ID:") +
			"\n" + indent + p.ID +
			"\n\n" + utilities.HeaderFormat(noColor, "OPTIONS:"),
	)

	displayPollContent(&builder, p, noColor, indent)

	builder.WriteString(
		"\n\n" +
			utilities.HeaderFormat(noColor, "MULTIPLE CHOICES ALLOWED:") +
			"\n" + indent + strconv.FormatBool(p.Multiple) +
			"\n\n" +
			utilities.HeaderFormat(noColor, "YOU VOTED:") +
			"\n" + indent + strconv.FormatBool(p.Voted),
	)

	if len(p.OwnVotes) > 0 {
		builder.WriteString("\n\n" + utilities.HeaderFormat(noColor, "YOUR VOTES:"))

		for _, vote := range p.OwnVotes {
			builder.WriteString("\n" + indent + "[" + strconv.Itoa(vote) + "] " + p.Options[vote].Title)
		}
	}

	builder.WriteString(
		"\n\n" +
			utilities.HeaderFormat(noColor, "EXPIRED:") +
			"\n" + indent + strconv.FormatBool(p.Expired),
	)

	return builder.String()
}

func displayPollContent(writer io.StringWriter, poll Poll, noColor bool, indent string) {
	for ind, option := range poll.Options {
		var percentage int
		var calculate float64

		if poll.VotesCount == 0 {
			percentage = 0
		} else {
			calculate = float64(option.VotesCount) / float64(poll.VotesCount)
			percentage = int(math.Floor(100 * calculate))
		}

		writer.WriteString("\n\n" + indent + "[" + strconv.Itoa(ind) + "] " + option.Title)
		drawPollMeter(writer, noColor, calculate, 80, indent)

		writer.WriteString(
			"\n" + indent + strconv.Itoa(option.VotesCount) + " votes " +
				"(" + strconv.Itoa(percentage) + "%)",
		)
	}

	writer.WriteString(
		"\n\n" +
			indent + utilities.FieldFormat(noColor, "Total votes:") + " " + strconv.Itoa(poll.VotesCount) +
			"\n" + indent + utilities.FieldFormat(noColor, "Poll ID:") + " " + poll.ID +
			"\n" + indent + utilities.FieldFormat(noColor, "Poll is open until:") + " " + utilities.FormatTime(poll.ExpiredAt),
	)
}

func drawPollMeter(writer io.StringWriter, noColor bool, calculated float64, limit int, indent string) {
	numVoteBlocks := int(math.Floor(float64(limit) * calculated))
	numBackgroundBlocks := limit - numVoteBlocks
	blockChar := "\u2501"
	voteBlockColor := "\033[32;1m"
	backgroundBlockColor := "\033[90m"

	if noColor {
		voteBlockColor = "\033[0m"

		if numVoteBlocks == 0 {
			numVoteBlocks = 1
		}
	}

	writer.WriteString("\n" + indent + voteBlockColor + strings.Repeat(blockChar, numVoteBlocks) + "\033[0m")

	if !noColor {
		writer.WriteString(backgroundBlockColor + strings.Repeat(blockChar, numBackgroundBlocks) + "\033[0m")
	}
}
