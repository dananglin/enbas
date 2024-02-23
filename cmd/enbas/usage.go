package main

import (
	"flag"
	"fmt"
	"slices"
	"strings"
)

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
				"\n  -%s, --%s\n        %s",
				f.Name,
				f.Name,
				f.Usage,
			)
		})

		builder.WriteString("\n")

		w := flag.CommandLine.Output()

		fmt.Fprint(w, builder.String())
	}
}

func enbasUsageFunc(summaries map[string]string) func() {
	cmds := make([]string, len(summaries))
	ind := 0

	for k := range summaries {
		cmds[ind] = k
		ind++
	}

	slices.Sort(cmds)

	return func() {
		var builder strings.Builder

		builder.WriteString("SUMMARY:\n    enbas - A GoToSocial client for the terminal.\n\n")

		if binaryVersion != "" {
			builder.WriteString("VERSION:\n  " + binaryVersion + "\n\n")
		}

		builder.WriteString("USAGE:\n    enbas [flags]\n    enbas [command]\n\nCOMMANDS:")

		for _, cmd := range cmds {
			fmt.Fprintf(&builder, "\n    %s\t%s", cmd, summaries[cmd])
		}

		builder.WriteString("\n\nFLAGS:\n    -help, --help\n        print the help message\n")
		flag.VisitAll(func(f *flag.Flag) {
			fmt.Fprintf(&builder, "\n    -%s, --%s\n        %s\n", f.Name, f.Name, f.Usage)
		})

		builder.WriteString("\nUse \"enbas [command] --help\" for more information about a command.\n")

		w := flag.CommandLine.Output()
		fmt.Fprint(w, builder.String())
	}
}
