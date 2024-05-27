package model

import (
	"strings"

	"codeflow.dananglin.me.uk/apollo/enbas/internal/utilities"
)

type Timeline struct {
	Name     string
	Statuses []Status
}

func (t Timeline) String() string {
	var builder strings.Builder

	separator := "────────────────────────────────────────────────────────────────────────────────"

	builder.WriteString(utilities.HeaderFormat(t.Name) + "\n\n")

	for _, status := range t.Statuses {
		builder.WriteString(utilities.DisplayNameFormat(status.Account.DisplayName) + " (@" + status.Account.Acct + ")\n")

		statusID := status.ID
		createdAt := status.CreatedAt

		if status.Reblog != nil {
			builder.WriteString("reposted this status from " + utilities.DisplayNameFormat(status.Reblog.Account.DisplayName) + " (@" + status.Reblog.Account.Acct + ")\n")
			statusID = status.Reblog.ID
			createdAt = status.Reblog.CreatedAt
		}

		builder.WriteString(utilities.WrapLines(utilities.StripHTMLTags(status.Content), "\n", 80) + "\n\n")
		builder.WriteString(utilities.FieldFormat("ID:") + " " + statusID + "\t" + utilities.FieldFormat("Created at:") + " " + utilities.FormatTime(createdAt) + "\n")
		builder.WriteString(separator + "\n")
	}

	return builder.String()
}
