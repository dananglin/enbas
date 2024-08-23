package usage

import (
	"flag"
	"slices"
	"strings"

	"codeflow.dananglin.me.uk/apollo/enbas/internal/info"
)

// ExecutorUsageFunc returns the function used to print a command's help page.
func ExecutorUsageFunc(name, summary string, flagset *flag.FlagSet) func() {
	return func() {
		var builder strings.Builder

		builder.WriteString("SUMMARY:")
		builder.WriteString("\n  " + name + " - " + summary)
		builder.WriteString("\n\nUSAGE:")
		builder.WriteString("\n  " + info.ApplicationName + " " + name)

		flagMap := make(map[string]string)

		flagset.VisitAll(func(f *flag.Flag) {
			flagMap[f.Name] = f.Usage
		})

		if len(flagMap) > 0 {
			flags := make([]string, len(flagMap))
			ind := 0

			for f := range flagMap {
				flags[ind] = f
				ind++
			}

			slices.Sort(flags)

			builder.WriteString(" [flags]")
			builder.WriteString("\n\nFLAGS:")

			for _, value := range flags {
				builder.WriteString("\n  --" + value)
				builder.WriteString("\n        " + flagMap[value])
			}
		}

		builder.WriteString("\n")

		w := flag.CommandLine.Output()

		_, _ = w.Write([]byte(builder.String()))
	}
}
