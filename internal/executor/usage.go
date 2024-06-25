// SPDX-FileCopyrightText: 2024 Dan Anglin <d.n.i.anglin@gmail.com>
//
// SPDX-License-Identifier: GPL-3.0-or-later

package executor

import (
	"flag"
	"slices"
	"strings"
)

// commandUsageFunc returns the function used to print a command's help page.
func commandUsageFunc(name, summary string, flagset *flag.FlagSet) func() {
	return func() {
		var builder strings.Builder

		builder.WriteString("SUMMARY:")
		builder.WriteString("\n  " + name + " - " + summary)
		builder.WriteString("\n\nUSAGE:")
		builder.WriteString("\n  enbas " + name)

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
