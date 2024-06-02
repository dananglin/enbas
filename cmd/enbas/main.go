// SPDX-FileCopyrightText: 2024 Dan Anglin <d.n.i.anglin@gmail.com>
//
// SPDX-License-Identifier: GPL-3.0-or-later

package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"

	"codeflow.dananglin.me.uk/apollo/enbas/internal/executor"
)

const (
	commandLogin    string = "login"
	commandVersion  string = "version"
	commandShow     string = "show"
	commandSwitch   string = "switch"
	commandCreate   string = "create"
	commandDelete   string = "delete"
	commandEdit     string = "edit"
	commandWhoami   string = "whoami"
	commandAdd      string = "add"
	commandRemove   string = "remove"
	commandFollow   string = "follow"
	commandUnfollow string = "unfollow"
	commandBlock    string = "block"
	commandUnblock  string = "unblock"
)

var (
	binaryVersion string
	buildTime     string
	goVersion     string
	gitCommit     string
)

func main() {
	if err := run(); err != nil {
		fmt.Printf("ERROR: %v.\n", err)
		os.Exit(1)
	}
}

func run() error {
	commandSummaries := map[string]string{
		commandLogin:    "login to an account on GoToSocial",
		commandVersion:  "print the application's version and build information",
		commandShow:     "print details about a specified resource",
		commandSwitch:   "perform a switch operation (e.g. switch logged in accounts)",
		commandCreate:   "create a specific resource",
		commandDelete:   "delete a specific resource",
		commandEdit:     "edit a specific resource",
		commandWhoami:   "print the account that you are currently logged in to",
		commandAdd:      "add a resource to another resource",
		commandRemove:   "remove a resource from another resource",
		commandFollow:   "follow a resource (e.g. an account)",
		commandUnfollow: "unfollow a resource (e.g. an account)",
		commandBlock:    "block a resource (e.g. an account)",
		commandUnblock:  "unblock a resource (e.g. an account)",
	}

	topLevelFlags := executor.TopLevelFlags{
		ConfigDir: "",
		NoColor:   nil,
	}

	flag.StringVar(&topLevelFlags.ConfigDir, "config-dir", "", "specify your config directory")
	flag.BoolFunc("no-color", "disable ANSI colour output when displaying text on screen", func(value string) error {
		boolVal, err := strconv.ParseBool(value)
		if err != nil {
			return fmt.Errorf("unable to parse %q as a boolean; %w", value, err)
		}

		topLevelFlags.NoColor = new(bool)
		*topLevelFlags.NoColor = boolVal

		return nil
	})

	flag.Usage = usageFunc(commandSummaries)

	flag.Parse()

	if flag.NArg() < 1 {
		flag.Usage()

		return nil
	}

	// If NoColor is still unspecified, check to see if the NO_COLOR environment variable is set
	if topLevelFlags.NoColor == nil {
		topLevelFlags.NoColor = new(bool)
		if os.Getenv("NO_COLOR") != "" {
			*topLevelFlags.NoColor = true
		} else {
			*topLevelFlags.NoColor = false
		}
	}

	command := flag.Arg(0)
	args := flag.Args()[1:]

	var err error

	switch command {
	case commandAdd:
		exe := executor.NewAddExecutor(
			topLevelFlags,
			commandAdd,
			commandSummaries[commandAdd],
		)
		err = executor.Execute(exe, args)
	case commandBlock:
		exe := executor.NewBlockExecutor(
			topLevelFlags,
			commandBlock,
			commandSummaries[commandBlock],
			false,
		)
		err = executor.Execute(exe, args)
	case commandCreate:
		exe := executor.NewCreateExecutor(
			topLevelFlags,
			commandCreate,
			commandSummaries[commandCreate],
		)
		err = executor.Execute(exe, args)
	case commandDelete:
		exe := executor.NewDeleteExecutor(
			topLevelFlags,
			commandDelete,
			commandSummaries[commandDelete],
		)
		err = executor.Execute(exe, args)
	case commandEdit:
		exe := executor.NewEditExecutor(
			topLevelFlags,
			commandEdit,
			commandSummaries[commandEdit],
		)
		err = executor.Execute(exe, args)
	case commandFollow:
		exe := executor.NewFollowExecutor(
			topLevelFlags,
			commandFollow,
			commandSummaries[commandFollow],
			false,
		)
		err = executor.Execute(exe, args)
	case commandLogin:
		exe := executor.NewLoginExecutor(
			topLevelFlags,
			commandLogin,
			commandSummaries[commandLogin],
		)
		err = executor.Execute(exe, args)
	case commandRemove:
		exe := executor.NewRemoveExecutor(
			topLevelFlags,
			commandRemove,
			commandSummaries[commandRemove],
		)
		err = executor.Execute(exe, args)
	case commandSwitch:
		exe := executor.NewSwitchExecutor(
			topLevelFlags,
			commandSwitch,
			commandSummaries[commandSwitch],
		)
		err = executor.Execute(exe, args)
	case commandUnfollow:
		exe := executor.NewFollowExecutor(topLevelFlags, commandUnfollow, commandSummaries[commandUnfollow], true)
		err = executor.Execute(exe, args)
	case commandUnblock:
		exe := executor.NewBlockExecutor(topLevelFlags, commandUnblock, commandSummaries[commandUnblock], true)
		err = executor.Execute(exe, args)
	case commandShow:
		exe := executor.NewShowExecutor(topLevelFlags, commandShow, commandSummaries[commandShow])
		err = executor.Execute(exe, args)
	case commandVersion:
		exe := executor.NewVersionExecutor(
			commandVersion,
			commandSummaries[commandVersion],
			binaryVersion,
			buildTime,
			goVersion,
			gitCommit,
		)
		err = executor.Execute(exe, args)
	case commandWhoami:
		exe := executor.NewWhoAmIExecutor(topLevelFlags, commandWhoami, commandSummaries[commandWhoami])
		err = executor.Execute(exe, args)
	default:
		flag.Usage()

		return unknownCommandError{command}
	}

	if err != nil {
		return fmt.Errorf("received an error executing the command; %w", err)
	}

	return nil
}
