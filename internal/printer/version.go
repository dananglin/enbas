package printer

import (
	"strings"
	"text/tabwriter"

	"codeflow.dananglin.me.uk/apollo/enbas/internal/info"
)

func (p Printer) PrintVersion(showFullVersion bool) {
	if !showFullVersion {
		printToStdout("Enbas " + info.BinaryVersion + "\n")

		return
	}

	var builder strings.Builder

	builder.WriteString(p.headerFormat("Enbas") + "\n\n")

	tableWriter := tabwriter.NewWriter(&builder, 0, 4, 1, ' ', 0)

	_, _ = tableWriter.Write([]byte(p.fieldFormat("Version:") + "\t" + info.BinaryVersion + "\n"))
	_, _ = tableWriter.Write([]byte(p.fieldFormat("Git commit:") + "\t" + info.GitCommit + "\n"))
	_, _ = tableWriter.Write([]byte(p.fieldFormat("Go version:") + "\t" + info.GoVersion + "\n"))
	_, _ = tableWriter.Write([]byte(p.fieldFormat("Build date:") + "\t" + info.BuildTime + "\n"))

	_ = tableWriter.Flush()

	printToStdout(builder.String())
}
