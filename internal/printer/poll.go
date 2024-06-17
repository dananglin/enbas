// SPDX-FileCopyrightText: 2024 Dan Anglin <d.n.i.anglin@gmail.com>
//
// SPDX-License-Identifier: GPL-3.0-or-later

package printer

import (
	"math"
	"strconv"
	"strings"

	"codeflow.dananglin.me.uk/apollo/enbas/internal/model"
)

func (p Printer) PrintPoll(poll model.Poll) {
	var builder strings.Builder

	builder.WriteString("\n" + p.headerFormat("POLL ID:"))
	builder.WriteString("\n" + poll.ID)

	builder.WriteString("\n\n" + p.headerFormat("OPTIONS:"))
	builder.WriteString(p.pollOptions(poll))

	builder.WriteString("\n\n" + p.headerFormat("MULTIPLE CHOICES ALLOWED:"))
	builder.WriteString("\n" + strconv.FormatBool(poll.Multiple))

	builder.WriteString("\n\n" + p.headerFormat("YOU VOTED:"))
	builder.WriteString("\n" + strconv.FormatBool(poll.Voted))

	if len(poll.OwnVotes) > 0 {
		builder.WriteString("\n\n" + p.headerFormat("YOUR VOTES:"))

		for _, vote := range poll.OwnVotes {
			builder.WriteString("\n" + "[" + strconv.Itoa(vote) + "] " + poll.Options[vote].Title)
		}
	}

	builder.WriteString("\n\n" + p.headerFormat("EXPIRED:"))
	builder.WriteString("\n" + strconv.FormatBool(poll.Expired))
	builder.WriteString("\n\n")

	p.print(builder.String())
}

func (p Printer) pollOptions(poll model.Poll) string {
	var builder strings.Builder

	for ind, option := range poll.Options {
		var (
			votage     float64
			percentage int
		)

		if poll.VotesCount == 0 {
			percentage = 0
		} else {
			votage = float64(option.VotesCount) / float64(poll.VotesCount)
			percentage = int(math.Floor(100 * votage))
		}

		builder.WriteString("\n\n" + "[" + strconv.Itoa(ind) + "] " + option.Title)
		builder.WriteString(p.pollMeter(votage))
		builder.WriteString("\n" + strconv.Itoa(option.VotesCount) + " votes " + "(" + strconv.Itoa(percentage) + "%)")
	}

	builder.WriteString("\n\n" + p.fieldFormat("Total votes:") + " " + strconv.Itoa(poll.VotesCount))
	builder.WriteString("\n" + p.fieldFormat("Poll ID:") + " " + poll.ID)
	builder.WriteString("\n" + p.fieldFormat("Poll is open until:") + " " + p.formatDateTime(poll.ExpiredAt))

	return builder.String()
}

func (p Printer) pollMeter(votage float64) string {
	numVoteBlocks := int(math.Floor(float64(p.maxTerminalWidth) * votage))
	numBackgroundBlocks := p.maxTerminalWidth - numVoteBlocks

	voteBlockColor := p.theme.boldgreen
	backgroundBlockColor := p.theme.grey

	if p.noColor {
		voteBlockColor = p.theme.reset

		if numVoteBlocks == 0 {
			numVoteBlocks = 1
		}
	}

	meter := "\n" + voteBlockColor + strings.Repeat(p.pollMeterSymbol, numVoteBlocks) + p.theme.reset

	if !p.noColor {
		meter += backgroundBlockColor + strings.Repeat(p.pollMeterSymbol, numBackgroundBlocks) + p.theme.reset
	}

	return meter
}
