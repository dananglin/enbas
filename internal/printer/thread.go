package printer

import (
	"strings"

	"codeflow.dananglin.me.uk/apollo/enbas/internal/model"
)

func (p Printer) PrintThread(thread model.Thread, userAccountID string) {
	var builder strings.Builder

	if len(thread.Ancestors.Statuses) > 0 {
		builder.WriteString(p.statusList(thread.Ancestors, userAccountID))
	}

	builder.WriteString(p.headerFormat("Context") + "\n")
	builder.WriteString(p.statusCard(thread.Context, userAccountID))

	if len(thread.Descendants.Statuses) > 0 {
		builder.WriteString(p.statusList(thread.Descendants, userAccountID))
	}

	p.print(builder.String())
}
