// SPDX-FileCopyrightText: 2024 Dan Anglin <d.n.i.anglin@gmail.com>
//
// SPDX-License-Identifier: GPL-3.0-or-later

package printer

import (
	"strconv"
	"strings"

	"codeflow.dananglin.me.uk/apollo/enbas/internal/model"
	"codeflow.dananglin.me.uk/apollo/enbas/internal/utilities"
)

func (p Printer) PrintStatus(status model.Status) {
	var builder strings.Builder

	// The account information
	builder.WriteString("\n" + p.fullDisplayNameFormat(status.Account.DisplayName, status.Account.Acct))

	// The ID of the status
	builder.WriteString("\n\n" + p.headerFormat("STATUS ID:"))
	builder.WriteString("\n" + status.ID)

	// The content of the status.
	builder.WriteString("\n\n" + p.headerFormat("CONTENT:"))
	builder.WriteString(utilities.WrapLines(utilities.ConvertHTMLToText(status.Content), "\n", p.maxTerminalWidth))

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
		builder.WriteString(p.pollOptions(*status.Poll))
	}

	// Status creation time
	builder.WriteString("\n\n" + p.headerFormat("CREATED AT:"))
	builder.WriteString("\n" + p.formatDateTime(status.CreatedAt))

	// Status stats
	builder.WriteString("\n\n" + p.headerFormat("STATS:"))
	builder.WriteString("\n" + p.fieldFormat("Boosts: ") + strconv.Itoa(status.ReblogsCount))
	builder.WriteString("\n" + p.fieldFormat("Likes: ") + strconv.Itoa(status.FavouritesCount))
	builder.WriteString("\n" + p.fieldFormat("Replies: ") + strconv.Itoa(status.RepliesCount))

	// Status visibility
	builder.WriteString("\n\n" + p.headerFormat("VISIBILITY:"))
	builder.WriteString("\n" + status.Visibility.String())

	// Status URL
	builder.WriteString("\n\n" + p.headerFormat("URL:"))
	builder.WriteString("\n" + status.URL)
	builder.WriteString("\n\n")

	p.print(builder.String())
}

func (p Printer) PrintStatusList(list model.StatusList) {
	var builder strings.Builder

	builder.WriteString(p.headerFormat(list.Name) + "\n")

	for _, status := range list.Statuses {
		builder.WriteString("\n" + p.fullDisplayNameFormat(status.Account.DisplayName, status.Account.Acct))

		statusID := status.ID
		createdAt := status.CreatedAt
		content := status.Content
		poll := status.Poll
		mediaAttachments := status.MediaAttachments

		if status.Reblog != nil {
			builder.WriteString(
				"\n" + utilities.WrapLines("reposted this status from "+p.fullDisplayNameFormat(status.Reblog.Account.DisplayName, status.Reblog.Account.Acct), "\n", p.maxTerminalWidth),
			)

			statusID = status.Reblog.ID
			createdAt = status.Reblog.CreatedAt
			content = status.Reblog.Content
			poll = status.Reblog.Poll
			mediaAttachments = status.Reblog.MediaAttachments
		}

		builder.WriteString("\n" + utilities.WrapLines(utilities.ConvertHTMLToText(content), "\n", p.maxTerminalWidth))

		if poll != nil {
			builder.WriteString(p.pollOptions(*poll))
		}

		for _, media := range mediaAttachments {
			builder.WriteString("\n\n" + p.imageIcon + "  Media attachment: " + media.ID)
			builder.WriteString("\n   Media type: " + media.Type + "\n   ")

			description := media.Description
			if description == "" {
				description = noMediaDescription
			}

			builder.WriteString(utilities.WrapLines(description, "\n   ", p.maxTerminalWidth-3))
		}

		builder.WriteString(
			"\n\n" +
				p.fieldFormat("Status ID:") + " " + statusID + "\t" +
				p.fieldFormat("Created at:") + " " + p.formatDateTime(createdAt) +
				"\n",
		)

		builder.WriteString(p.statusSeparator + "\n")
	}

	p.print(builder.String())
}
