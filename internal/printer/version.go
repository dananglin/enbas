// SPDX-FileCopyrightText: 2024 Dan Anglin <d.n.i.anglin@gmail.com>
//
// SPDX-License-Identifier: GPL-3.0-or-later

package printer

import (
	"strings"
	"text/tabwriter"
)

func (p Printer) PrintVersion(showFullVersion bool, binaryVersion, buildTime, goVersion, gitCommit string) {
	if !showFullVersion {
		printToStdout("Enbas " + binaryVersion + "\n")

		return
	}

	var builder strings.Builder

	builder.WriteString(p.headerFormat("Enbas") + "\n\n")

	tableWriter := tabwriter.NewWriter(&builder, 0, 4, 1, ' ', 0)

	_, _ = tableWriter.Write([]byte(p.fieldFormat("Version:") + "\t" + binaryVersion + "\n"))
	_, _ = tableWriter.Write([]byte(p.fieldFormat("Git commit:") + "\t" + gitCommit + "\n"))
	_, _ = tableWriter.Write([]byte(p.fieldFormat("Go version:") + "\t" + goVersion + "\n"))
	_, _ = tableWriter.Write([]byte(p.fieldFormat("Build date:") + "\t" + buildTime + "\n"))

	tableWriter.Flush()

	printToStdout(builder.String())
}
