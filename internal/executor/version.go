// SPDX-FileCopyrightText: 2024 Dan Anglin <d.n.i.anglin@gmail.com>
//
// SPDX-License-Identifier: GPL-3.0-or-later

package executor

import (
	"flag"
	"fmt"
	"os"
	"strings"
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

func (c *VersionExecutor) Execute() error {
	var builder strings.Builder

	if c.showFullVersion {
		fmt.Fprintf(
			&builder,
			"Enbas\n  Version: %s\n  Git commit: %s\n  Go version: %s\n  Build date: %s\n",
			c.binaryVersion,
			c.gitCommit,
			c.goVersion,
			c.buildTime,
		)
	} else {
		fmt.Fprintln(&builder, c.binaryVersion)
	}

	fmt.Fprint(os.Stdout, builder.String())

	return nil
}
