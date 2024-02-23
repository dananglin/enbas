package main

import (
	"flag"
	"fmt"
	"os"
)

type Executor interface {
	Name() string
	Parse([]string) error
	Execute() error
}

func main() {
	if err := run(); err != nil {
		fmt.Printf("ERROR: %v.\n", err)
		os.Exit(1)
	}
}

func run() error {
	const (
		login         string = "login"
		version       string = "version"
		show          string = "show"
		switchAccount string = "switch"
	)

	summaries := map[string]string{
		login:         "login to an account on GoToSocial",
		version:       "print the application's version and build information",
		show:          "print details about a specified resource",
		switchAccount: "switch to an account",
	}

	flag.Usage = enbasUsageFunc(summaries)

	flag.Parse()

	if flag.NArg() < 1 {
		flag.Usage()

		return nil
	}

	subcommand := flag.Arg(0)
	args := flag.Args()[1:]

	var executor Executor

	switch subcommand {
	case login:
		executor = newLoginCommand(login, summaries[login])
	case version:
		executor = newVersionCommand(version, summaries[version])
	case show:
		executor = newShowCommand(show, summaries[show])
	case switchAccount:
		executor = newSwitchCommand(switchAccount, summaries[switchAccount])
	default:
		flag.Usage()
		return fmt.Errorf("unknown subcommand %q", subcommand)
	}

	if err := executor.Parse(args); err != nil {
		return fmt.Errorf("unable to parse the command line flags; %w", err)
	}

	if err := executor.Execute(); err != nil {
		return fmt.Errorf("received an error after executing %q; %w", executor.Name(), err)
	}

	return nil
}
