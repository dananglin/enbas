package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

var (
	binaryVersion string
	buildTime     string
	goVersion     string
	gitCommit     string
)

type versionCommand struct {
	*flag.FlagSet
	showFullVersion bool
	binaryVersion   string
	buildTime       string
	goVersion       string
	gitCommit       string
}

func newVersionCommand(name, summary string) *versionCommand {
	command := versionCommand{
		FlagSet:         flag.NewFlagSet(name, flag.ExitOnError),
		binaryVersion:   binaryVersion,
		buildTime:       buildTime,
		goVersion:       goVersion,
		gitCommit:       gitCommit,
		showFullVersion: false,
	}

	command.BoolVar(&command.showFullVersion, "full", false, "prints the full build information")

	command.Usage = commandUsageFunc(name, summary, command.FlagSet)

	return &command
}

func (c *versionCommand) Execute() error {
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
