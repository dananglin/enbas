package usage

import (
	"flag"
	"fmt"
	"slices"
	"strings"
	"text/tabwriter"

	"codeflow.dananglin.me.uk/apollo/enbas/internal/info"
)

func AppUsageFunc() func() {
	cmds := make([]string, len(summaries))
	ind := 0

	for k := range summaries {
		cmds[ind] = k
		ind++
	}

	slices.Sort(cmds)

	return func() {
		var builder strings.Builder

		builder.WriteString("SUMMARY:\n    " + info.ApplicationName + " - A GoToSocial client for the terminal.\n\n")

		if info.BinaryVersion != "" {
			builder.WriteString("VERSION:\n    " + info.BinaryVersion + "\n\n")
		}

		builder.WriteString("USAGE:\n    " + info.ApplicationName + " [flags]\n    " + info.ApplicationName + " [flags] [command]\n\nCOMMANDS:")

		tableWriter := tabwriter.NewWriter(&builder, 0, 8, 0, '\t', 0)

		for _, cmd := range cmds {
			fmt.Fprintf(tableWriter, "\n    %s\t%s", cmd, summaries[cmd])
		}

		tableWriter.Flush()

		builder.WriteString("\n\nFLAGS:\n    --help\n        print the help message")
		flag.VisitAll(func(f *flag.Flag) {
			fmt.Fprintf(&builder, "\n    --%s\n        %s", f.Name, f.Usage)
		})

		builder.WriteString("\n\nUse \"" + info.ApplicationName + " [command] --help\" for more information about a command.\n")

		w := flag.CommandLine.Output()
		fmt.Fprint(w, builder.String())
	}
}
