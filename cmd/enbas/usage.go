// SPDX-FileCopyrightText: 2024 Dan Anglin <d.n.i.anglin@gmail.com>
//
// SPDX-License-Identifier: GPL-3.0-or-later

package main

import (
	"flag"
	"fmt"
	"slices"
	"strings"
	"text/tabwriter"
)

func usageFunc(summaries map[string]string) func() {
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
			builder.WriteString("VERSION:\n    " + binaryVersion + "\n\n")
		}

		builder.WriteString("USAGE:\n    enbas [flags]\n    enbas [flags] [command]\n\nCOMMANDS:")

		tableWriter := tabwriter.NewWriter(&builder, 0, 8, 0, '\t', 0)

		for _, cmd := range cmds {
			fmt.Fprintf(tableWriter, "\n    %s\t%s", cmd, summaries[cmd])
		}

		tableWriter.Flush()

		builder.WriteString("\n\nFLAGS:\n    --help\n        print the help message")
		flag.VisitAll(func(f *flag.Flag) {
			fmt.Fprintf(&builder, "\n    --%s\n        %s", f.Name, f.Usage)
		})

		builder.WriteString("\n\nUse \"enbas [command] --help\" for more information about a command.\n")

		w := flag.CommandLine.Output()
		fmt.Fprint(w, builder.String())
	}
}
