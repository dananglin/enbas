package main

import (
	"flag"
	"os"

	"codeflow.dananglin.me.uk/apollo/enbas/internal/executor"
	internalFlag "codeflow.dananglin.me.uk/apollo/enbas/internal/flag"
	"codeflow.dananglin.me.uk/apollo/enbas/internal/usage"
)

func main() {
	if err := run(); err != nil {
		os.Exit(1)
	}
}

func run() error {
	var (
		configDir   string
		noColorFlag internalFlag.BoolPtrValue
	)

	flag.StringVar(&configDir, "config-dir", "", "Specify your config directory")
	flag.Var(&noColorFlag, "no-color", "Disable ANSI colour output when displaying text on screen")

	flag.Usage = usage.AppUsageFunc()

	flag.Parse()

	if flag.NArg() < 1 {
		flag.Usage()

		return nil
	}

	var noColor bool

	if noColorFlag.Value != nil {
		noColor = *noColorFlag.Value
	} else if os.Getenv("NO_COLOR") != "" {
		noColor = true
	}

	command := flag.Arg(0)
	args := flag.Args()[1:]

	return executor.Execute(command, args, noColor, configDir)
}
