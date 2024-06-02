// SPDX-FileCopyrightText: 2024 Dan Anglin <d.n.i.anglin@gmail.com>
//
// SPDX-License-Identifier: GPL-3.0-or-later

package model

import (
	"strings"

	"codeflow.dananglin.me.uk/apollo/enbas/internal/utilities"
)

const (
	TimelineCategoryHome   = "home"
	TimelineCategoryPublic = "public"
	TimelineCategoryTag    = "tag"
	TimelineCategoryList   = "list"
)

type InvalidTimelineCategoryError struct {
	Value string
}

func (e InvalidTimelineCategoryError) Error() string {
	return "'" +
		e.Value +
		"' is not a valid timeline category (valid values are " +
		TimelineCategoryHome + ", " +
		TimelineCategoryPublic + ", " +
		TimelineCategoryTag + ", " +
		TimelineCategoryList + ")"
}

type Timeline struct {
	Name     string
	Statuses []Status
}

func (t Timeline) Display(noColor bool) string {
	var builder strings.Builder

	separator := "────────────────────────────────────────────────────────────────────────────────"

	builder.WriteString(utilities.HeaderFormat(noColor, t.Name) + "\n")

	for _, status := range t.Statuses {
		builder.WriteString("\n" + utilities.DisplayNameFormat(noColor, status.Account.DisplayName) + " (@" + status.Account.Acct + ")\n")

		statusID := status.ID
		createdAt := status.CreatedAt

		if status.Reblog != nil {
			builder.WriteString("reposted this status from " + utilities.DisplayNameFormat(noColor, status.Reblog.Account.DisplayName) + " (@" + status.Reblog.Account.Acct + ")\n")
			statusID = status.Reblog.ID
			createdAt = status.Reblog.CreatedAt
		}

		builder.WriteString(utilities.WrapLines(utilities.ConvertHTMLToText(status.Content), "\n", 80) + "\n\n")
		builder.WriteString(utilities.FieldFormat(noColor, "ID:") + " " + statusID + "\t" + utilities.FieldFormat(noColor, "Created at:") + " " + utilities.FormatTime(createdAt) + "\n")
		builder.WriteString(separator + "\n")
	}

	return builder.String()
}
