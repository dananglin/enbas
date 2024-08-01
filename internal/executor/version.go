package executor

import (
	"flag"

	"codeflow.dananglin.me.uk/apollo/enbas/internal/printer"
)

type VersionExecutor struct {
	*flag.FlagSet

	printer         *printer.Printer
	showFullVersion bool
	binaryVersion   string
	buildTime       string
	goVersion       string
	gitCommit       string
}

func NewVersionExecutor(
	printer *printer.Printer,
	name,
	summary,
	binaryVersion,
	buildTime,
	goVersion,
	gitCommit string,
) *VersionExecutor {
	command := VersionExecutor{
		FlagSet: flag.NewFlagSet(name, flag.ExitOnError),

		printer:         printer,
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
	v.printer.PrintVersion(v.showFullVersion, v.binaryVersion, v.buildTime, v.goVersion, v.gitCommit)

	return nil
}
