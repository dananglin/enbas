// SPDX-FileCopyrightText: 2024 Dan Anglin <d.n.i.anglin@gmail.com>
//
// SPDX-License-Identifier: GPL-3.0-or-later

package printer

import (
	"strings"

	"codeflow.dananglin.me.uk/apollo/enbas/internal/model"
)

func (p Printer) PrintList(list model.List) {
	var builder strings.Builder

	builder.WriteString("\n" + p.headerFormat("LIST TITLE:") + "\n")
	builder.WriteString(list.Title + "\n\n")
	builder.WriteString(p.headerFormat("LIST ID:") + "\n")
	builder.WriteString(list.ID + "\n\n")
	builder.WriteString(p.headerFormat("REPLIES POLICY:") + "\n")
	builder.WriteString(list.RepliesPolicy.String() + "\n\n")
	builder.WriteString(p.headerFormat("ADDED ACCOUNTS:"))

	if len(list.Accounts) > 0 {
		for acct, name := range list.Accounts {
			builder.WriteString("\n" + p.bullet + " " + p.fullDisplayNameFormat(name, acct))
		}
	} else {
		builder.WriteString("\n" + "None")
	}

	builder.WriteString("\n")

	printToStdout(builder.String())
}

func (p Printer) PrintLists(lists []model.List) {
	var builder strings.Builder

	builder.WriteString("\n" + p.headerFormat("LISTS"))

	for i := range lists {
		builder.WriteString("\n" + p.bullet + " " + lists[i].Title + " (" + lists[i].ID + ")")
	}

	builder.WriteString("\n")

	printToStdout(builder.String())
}
