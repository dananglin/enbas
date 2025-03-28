package printer

import (
	"math"
	"strings"
)

type PollOptionDetails struct {
	Percentage int
	Meter      string
	Voted      bool
}

func getPollOptionDetails(
	noColor bool,
	charLimit int,
) func(int, int, int, []int) PollOptionDetails {
	return func(numVotesForOption, totalVotes, optionID int, ownVotes []int) PollOptionDetails {
		const symbolPollMeter = "\u2501"

		var (
			votage     float64
			percentage int
		)

		if totalVotes > 0 {
			votage = float64(numVotesForOption) / float64(totalVotes)
			percentage = int(math.Floor(100 * votage))
		}

		numVoteBlocks := int(math.Floor(float64(charLimit) * votage))
		numBackgroundBlocks := charLimit - numVoteBlocks

		voteBlockColour := boldgreen
		backgroundBlockColor := grey

		if noColor {
			voteBlockColour = reset

			if numVoteBlocks == 0 {
				numVoteBlocks = 1
			}
		}

		meter := voteBlockColour + strings.Repeat(symbolPollMeter, numVoteBlocks) + reset

		if !noColor {
			meter += backgroundBlockColor + strings.Repeat(symbolPollMeter, numBackgroundBlocks) + reset
		}

		voted := false

		for _, vote := range ownVotes {
			if vote == optionID {
				voted = true

				break
			}
		}

		return PollOptionDetails{
			Percentage: percentage,
			Meter:      meter,
			Voted:      voted,
		}
	}
}
