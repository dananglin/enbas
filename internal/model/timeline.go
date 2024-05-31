package model

import (
	"strings"

	"codeflow.dananglin.me.uk/apollo/enbas/internal/utilities"
)

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

		builder.WriteString(utilities.WrapLines(utilities.StripHTMLTags(status.Content), "\n", 80) + "\n\n")
		builder.WriteString(utilities.FieldFormat(noColor, "ID:") + " " + statusID + "\t" + utilities.FieldFormat(noColor, "Created at:") + " " + utilities.FormatTime(createdAt) + "\n")
		builder.WriteString(separator + "\n")
	}

	return builder.String()
}
