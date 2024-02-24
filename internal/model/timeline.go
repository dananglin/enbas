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

	builder.WriteString(t.Name + "\n" + separator + "\n")

	for _, status := range t.Statuses {
		builder.WriteString(utilities.DisplayNameFormat(status.Account.DisplayName) + " (@" + status.Account.Username + ")\n\n")
		builder.WriteString(utilities.WrapLines(status.Text, "\n", 80) + "\n\n")
		builder.WriteString(utilities.FieldFormat("ID:") + " " + status.ID + "\t" + utilities.FieldFormat("Created at:") + " " + utilities.FormatTime(status.CreatedAt) + "\n")
		builder.WriteString(separator + "\n")
	}

	return builder.String()
}
