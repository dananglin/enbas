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

	// The content of the status.
	builder.WriteString("\n\n" + p.headerFormat("CONTENT:"))
	builder.WriteString(utilities.WrapLines(utilities.ConvertHTMLToText(status.Content), "\n", p.maxTerminalWidth))

	// If a poll exists in a status, write the contents to the builder.
	if status.Poll != nil {
		builder.WriteString(p.pollOptions(*status.Poll))
	}

	// The ID of the status
	builder.WriteString("\n\n" + p.headerFormat("STATUS ID:"))
	builder.WriteString("\n" + status.ID)

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

		if status.Reblog != nil {
			builder.WriteString(
				"\n" + utilities.WrapLines("reposted this status from "+p.fullDisplayNameFormat(status.Reblog.Account.DisplayName, status.Reblog.Account.Acct), "\n", p.maxTerminalWidth),
			)

			statusID = status.Reblog.ID
			createdAt = status.Reblog.CreatedAt
		}

		builder.WriteString("\n" + utilities.WrapLines(utilities.ConvertHTMLToText(status.Content), "\n", p.maxTerminalWidth))

		if status.Poll != nil {
			builder.WriteString(p.pollOptions(*status.Poll))
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
