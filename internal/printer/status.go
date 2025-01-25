package printer

import (
	"math"
	"strconv"
	"strings"

	"codeflow.dananglin.me.uk/apollo/enbas/internal/model"
)

func (p Printer) PrintStatus(
	status model.Status,
	userAccountID string,
	boostedBy model.AccountList,
	likedBy model.AccountList,
) {
	var builder strings.Builder

	// The account information
	builder.WriteString("\n" + p.fullDisplayNameFormat(status.Account.DisplayName, status.Account.Acct))

	// The ID of the status
	builder.WriteString("\n\n" + p.headerFormat("STATUS ID:"))
	builder.WriteString("\n" + status.ID)

	// The subject, summary of content warning of the status
	if status.SpoilerText != "" {
		builder.WriteString("\n\n" + p.headerFormat("SUMMARY:"))
		builder.WriteString("\n" + status.SpoilerText)
	}

	// The content of the status.
	builder.WriteString("\n\n" + p.headerFormat("CONTENT:"))
	builder.WriteString(p.convertHTMLToText(status.Content, true))

	// Details of media attachments (if any).
	if len(status.MediaAttachments) > 0 {
		builder.WriteString("\n\n" + p.headerFormat("MEDIA ATTACHMENTS:"))

		for ind, media := range status.MediaAttachments {
			builder.WriteString("\n\n[" + strconv.Itoa(ind) + "] " + p.fieldFormat("ID:") + " " + media.ID)
			builder.WriteString("\n    " + p.fieldFormat("Type:") + " " + media.Type)

			description := media.Description
			if description == "" {
				description = noMediaDescription
			}

			builder.WriteString("\n    " + p.fieldFormat("Description:") + " " + description)
			builder.WriteString("\n    " + p.fieldFormat("Media URL:") + " " + media.URL)
		}
	}

	// If a poll exists in a status, write the contents to the builder.
	if status.Poll != nil {
		pollOwner := false
		if status.Account.ID == userAccountID {
			pollOwner = true
		}

		builder.WriteString("\n\n" + p.headerFormat("POLL DETAILS:"))
		builder.WriteString(p.pollDetails(*status.Poll, pollOwner))
	}

	// Status creation time
	builder.WriteString("\n\n" + p.headerFormat("CREATED AT:"))
	builder.WriteString("\n" + p.formatDateTime(status.CreatedAt))

	// Status stats
	builder.WriteString("\n\n" + p.headerFormat("STATS:"))
	builder.WriteString("\n" + p.fieldFormat("Boosts: ") + strconv.Itoa(status.ReblogsCount))
	builder.WriteString("\n" + p.fieldFormat("Likes: ") + strconv.Itoa(status.FavouritesCount))
	builder.WriteString("\n" + p.fieldFormat("Replies: ") + strconv.Itoa(status.RepliesCount))

	// The user's actions on the status
	builder.WriteString("\n\n" + p.headerFormat("YOUR ACTIONS:"))
	builder.WriteString("\n" + p.fieldFormat("Boosted: ") + strconv.FormatBool(status.Reblogged))
	builder.WriteString("\n" + p.fieldFormat("Liked: ") + strconv.FormatBool(status.Favourited))
	builder.WriteString("\n" + p.fieldFormat("Bookmarked: ") + strconv.FormatBool(status.Bookmarked))
	builder.WriteString("\n" + p.fieldFormat("Muted: ") + strconv.FormatBool(status.Muted))

	// Status visibility
	builder.WriteString("\n\n" + p.headerFormat("VISIBILITY:"))
	builder.WriteString("\n" + status.Visibility.String())

	// Status URL
	builder.WriteString("\n\n" + p.headerFormat("URL:"))
	builder.WriteString("\n" + status.URL)

	// List of accounts that has boosted the status
	if boostedBy.Accounts != nil {
		builder.WriteString("\n\n" + p.accountList(boostedBy))
	}

	// List of accounts that have liked/starred the status
	if likedBy.Accounts != nil {
		builder.WriteString("\n\n" + p.accountList(likedBy))
	}

	builder.WriteString("\n\n")

	p.print(builder.String())
}

func (p Printer) PrintStatusList(list model.StatusList, userAccountID string) {
	p.print(p.statusList(list, userAccountID))
}

func (p Printer) statusList(list model.StatusList, userAccountID string) string {
	var builder strings.Builder

	builder.WriteString(p.headerFormat(list.Name) + "\n")

	for _, status := range list.Statuses {
		statusID := status.ID
		statusOwnerID := status.Account.ID
		createdAt := p.formatDateTime(status.CreatedAt)
		boostedAt := ""
		content := status.Content
		poll := status.Poll
		mediaAttachments := status.MediaAttachments
		summary := status.SpoilerText

		switch {
		case status.Reblog != nil:
			builder.WriteString("\n" + p.wrapLines(
				p.fullDisplayNameFormat(status.Account.DisplayName, status.Account.Acct)+
					" boosted this status from "+
					p.fullDisplayNameFormat(status.Reblog.Account.DisplayName, status.Reblog.Account.Acct)+
					":",
				0,
			))

			statusID = status.Reblog.ID
			statusOwnerID = status.Reblog.Account.ID
			createdAt = p.formatDateTime(status.Reblog.CreatedAt)
			boostedAt = p.formatDateTime(status.CreatedAt)
			content = status.Reblog.Content
			poll = status.Reblog.Poll
			mediaAttachments = status.Reblog.MediaAttachments
			summary = status.Reblog.SpoilerText

		case status.InReplyToID != "":
			builder.WriteString("\n" + p.wrapLines(
				p.fullDisplayNameFormat(status.Account.DisplayName, status.Account.Acct)+
					" posted in reply to "+
					status.InReplyToID+
					":",
				0,
			))
		default:
			builder.WriteString("\n" + p.fullDisplayNameFormat(status.Account.DisplayName, status.Account.Acct) + " posted:")
		}

		if summary != "" {
			builder.WriteString("\n\n" + p.bold(p.wrapLines(summary, 0)))
		}

		builder.WriteString("\n" + p.convertHTMLToText(content, true))

		if poll != nil {
			pollOwner := false
			if statusOwnerID == userAccountID {
				pollOwner = true
			}

			builder.WriteString(p.pollDetails(*poll, pollOwner))
		}

		for _, media := range mediaAttachments {
			builder.WriteString("\n\n" + symbolImage + " " + p.fieldFormat("Media attachment: ") + media.ID)
			builder.WriteString("\n  " + p.fieldFormat("Media type: ") + media.Type + "\n")

			description := "  " + p.fieldFormat("Description: ")

			if media.Description == "" {
				description += noMediaDescription
			} else {
				description += media.Description
			}

			builder.WriteString(p.wrapLines(description, 2))
		}

		boosted := symbolBoosted
		if status.Reblogged {
			boosted = p.theme.boldyellow + symbolBoosted + p.theme.reset
		}

		liked := symbolNotLiked
		if status.Favourited {
			liked = p.theme.boldyellow + symbolLiked + p.theme.reset
		}

		bookmarked := symbolNotBookmarked
		if status.Bookmarked {
			bookmarked = p.theme.boldyellow + symbolBookmarked + p.theme.reset
		}

		builder.WriteString("\n\n" + boosted + " " + p.fieldFormat("boosted: ") + strconv.FormatBool(status.Reblogged))
		builder.WriteString("\n" + liked + " " + p.fieldFormat("liked: ") + strconv.FormatBool(status.Favourited))
		builder.WriteString("\n" + bookmarked + " " + p.fieldFormat("bookmarked: ") + strconv.FormatBool(status.Bookmarked))

		builder.WriteString(
			"\n\n" +
				p.fieldFormat("Status ID: ") + statusID +
				"\n" + p.fieldFormat("Created at: ") + createdAt,
		)

		if boostedAt != "" {
			builder.WriteString("\n" + p.fieldFormat("Boosted at: ") + boostedAt)
		}

		builder.WriteString("\n" + p.statusSeparator + "\n")
	}

	return builder.String()
}

func (p Printer) pollDetails(poll model.Poll, owner bool) string {
	var builder strings.Builder

	for ind, option := range poll.Options {
		var (
			votage     float64
			percentage int
		)

		// Show the poll results under any of the following conditions:
		//     - the user is the owner of the poll
		//     - the poll has expired
		//     - the user has voted in the poll
		if owner || poll.Expired || poll.Voted {
			if poll.VotesCount == 0 {
				percentage = 0
			} else {
				votage = float64(option.VotesCount) / float64(poll.VotesCount)
				percentage = int(math.Floor(100 * votage))
			}

			optionTitle := "\n\n" + "[" + strconv.Itoa(ind) + "] " + option.Title

			for _, vote := range poll.OwnVotes {
				if ind == vote {
					optionTitle += " " + symbolCheckMark

					break
				}
			}

			builder.WriteString(optionTitle)
			builder.WriteString(p.pollMeter(votage))
			builder.WriteString("\n" + strconv.Itoa(option.VotesCount) + " votes " + "(" + strconv.Itoa(percentage) + "%)")
		} else {
			builder.WriteString("\n" + "[" + strconv.Itoa(ind) + "] " + option.Title)
		}
	}

	pollStatusField := "Poll is open until: "
	if poll.Expired {
		pollStatusField = "Poll was closed on: "
	}

	builder.WriteString("\n\n" + p.fieldFormat(pollStatusField) + p.formatDateTime(poll.ExpiredAt))
	builder.WriteString("\n" + p.fieldFormat("Total votes: ") + strconv.Itoa(poll.VotesCount))
	builder.WriteString("\n" + p.fieldFormat("Multiple choices allowed: ") + strconv.FormatBool(poll.Multiple))

	return builder.String()
}

func (p Printer) pollMeter(votage float64) string {
	numVoteBlocks := int(math.Floor(float64(p.lineWrapCharacterLimit) * votage))
	numBackgroundBlocks := p.lineWrapCharacterLimit - numVoteBlocks

	voteBlockColour := p.theme.boldgreen
	backgroundBlockColor := p.theme.grey

	if p.noColor {
		voteBlockColour = p.theme.reset

		if numVoteBlocks == 0 {
			numVoteBlocks = 1
		}
	}

	meter := "\n" + voteBlockColour + strings.Repeat(symbolPollMeter, numVoteBlocks) + p.theme.reset

	if !p.noColor {
		meter += backgroundBlockColor + strings.Repeat(symbolPollMeter, numBackgroundBlocks) + p.theme.reset
	}

	return meter
}
