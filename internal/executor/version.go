// SPDX-FileCopyrightText: 2024 Dan Anglin <d.n.i.anglin@gmail.com>
//
// SPDX-License-Identifier: GPL-3.0-or-later

package executor

import (
	"flag"
	"os"
	"strings"
	"text/tabwriter"
)

type VersionExecutor struct {
	*flag.FlagSet
	showFullVersion bool
	binaryVersion   string
	buildTime       string
	goVersion       string
	gitCommit       string
}

func NewVersionExecutor(name, summary, binaryVersion, buildTime, goVersion, gitCommit string) *VersionExecutor {
	command := VersionExecutor{
		FlagSet:         flag.NewFlagSet(name, flag.ExitOnError),
		binaryVersion:   binaryVersion,
		buildTime:       buildTime,
		goVersion:       goVersion,
		gitCommit:       gitCommit,
		showFullVersion: false,
	}

	command.BoolVar(&command.showFullVersion, flagFull, false, "prints the full build information")

	command.Usage = commandUsageFunc(name, summary, command.FlagSet)

	return &command
}

func (v *VersionExecutor) Execute() error {
	var builder strings.Builder

	if v.showFullVersion {
		builder.WriteString("Enbas\n")

		tableWriter := tabwriter.NewWriter(&builder, 0, 8, 0, '\t', 0)

		tableWriter.Write([]byte("    Version:\t" + v.binaryVersion + "\n"))
		tableWriter.Write([]byte("    Git commit:\t" + v.gitCommit + "\n"))
		tableWriter.Write([]byte("    Go version:\t" + v.goVersion + "\n"))
		tableWriter.Write([]byte("    Build date:\t" + v.buildTime + "\n"))

		tableWriter.Flush()
	} else {
		builder.WriteString("Enbas " + v.binaryVersion + "\n")
	}

	os.Stdout.WriteString(builder.String())

	return nil
}
