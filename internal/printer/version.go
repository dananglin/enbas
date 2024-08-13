package printer

import (
	"strings"
	"text/tabwriter"

	"codeflow.dananglin.me.uk/apollo/enbas/internal/version"
)

func (p Printer) PrintVersion(showFullVersion bool) {
	if !showFullVersion {
		printToStdout("Enbas " + version.BinaryVersion + "\n")

		return
	}

	var builder strings.Builder

	builder.WriteString(p.headerFormat("Enbas") + "\n\n")

	tableWriter := tabwriter.NewWriter(&builder, 0, 4, 1, ' ', 0)

	_, _ = tableWriter.Write([]byte(p.fieldFormat("Version:") + "\t" + version.BinaryVersion + "\n"))
	_, _ = tableWriter.Write([]byte(p.fieldFormat("Git commit:") + "\t" + version.GitCommit + "\n"))
	_, _ = tableWriter.Write([]byte(p.fieldFormat("Go version:") + "\t" + version.GoVersion + "\n"))
	_, _ = tableWriter.Write([]byte(p.fieldFormat("Build date:") + "\t" + version.BuildTime + "\n"))

	tableWriter.Flush()

	printToStdout(builder.String())
}
