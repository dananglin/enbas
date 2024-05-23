package executor

import (
	"flag"
	"fmt"
	"strings"
)

// commandUsageFunc returns the function used to print a command's help page.
func commandUsageFunc(name, summary string, flagset *flag.FlagSet) func() {
	return func() {
		var builder strings.Builder

		fmt.Fprintf(
			&builder,
			"SUMMARY:\n  %s - %s\n\nUSAGE:\n  enbas %s [flags]\n\nFLAGS:",
			name,
			summary,
			name,
		)

		flagset.VisitAll(func(f *flag.Flag) {
			fmt.Fprintf(
				&builder,
				"\n  --%s\n        %s",
				f.Name,
				f.Usage,
			)
		})

		builder.WriteString("\n")

		w := flag.CommandLine.Output()

		fmt.Fprint(w, builder.String())
	}
}
