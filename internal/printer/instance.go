// SPDX-FileCopyrightText: 2024 Dan Anglin <d.n.i.anglin@gmail.com>
//
// SPDX-License-Identifier: GPL-3.0-or-later

package printer

import (
	"strings"

	"codeflow.dananglin.me.uk/apollo/enbas/internal/model"
)

func (p Printer) PrintInstance(instance model.InstanceV2) {
	var builder strings.Builder

	builder.WriteString("\n" + p.headerFormat("INSTANCE TITLE:"))
	builder.WriteString("\n" + instance.Title)

	builder.WriteString("\n\n" + p.headerFormat("INSTANCE DESCRIPTION:"))
	builder.WriteString("\n" + p.wrapLines(instance.DescriptionText, 0))

	builder.WriteString("\n\n" + p.headerFormat("DOMAIN:"))
	builder.WriteString("\n" + instance.Domain)

	builder.WriteString("\n\n" + p.headerFormat("TERMS AND CONDITIONS:"))
	builder.WriteString("\n" + p.wrapLines(instance.TermsText, 2))

	builder.WriteString("\n\n" + p.headerFormat("VERSION:"))
	builder.WriteString("\nRunning GoToSocial " + instance.Version)

	builder.WriteString("\n\n" + p.headerFormat("CONTACT:"))
	builder.WriteString("\n" + p.fieldFormat("Name:"))
	builder.WriteString(" " + instance.Contact.Account.DisplayName)
	builder.WriteString("\n" + p.fieldFormat("Username:"))
	builder.WriteString(" " + instance.Contact.Account.Acct)
	builder.WriteString("\n" + p.fieldFormat("Email:"))
	builder.WriteString(" " + instance.Contact.Email)

	builder.WriteString("\n\n")

	p.print(builder.String())
}
