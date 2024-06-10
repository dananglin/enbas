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

var (
	binaryVersion string
	buildTime     string
	goVersion     string
	gitCommit     string
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %v.\n", err)
		os.Exit(1)
	}
}

func run() error {
	topLevelFlags := executor.TopLevelFlags{
		ConfigDir: "",
		NoColor:   nil,
		Pager:     "",
	}

	flag.StringVar(
		&topLevelFlags.ConfigDir,
		"config-dir",
		"",
		"Specify your config directory",
	)

	flag.BoolFunc(
		"no-color",
		"Disable ANSI colour output when displaying text on screen",
		func(value string) error {
			boolVal, err := strconv.ParseBool(value)
			if err != nil {
				return fmt.Errorf("unable to parse %q as a boolean: %w", value, err)
			}

			topLevelFlags.NoColor = new(bool)
			*topLevelFlags.NoColor = boolVal

			return nil
		},
	)

	flag.StringVar(
		&topLevelFlags.Pager,
		"pager",
		"",
		"Specify your preferred pager to page through long outputs. This is disabled by default.",
	)

	flag.Usage = usageFunc(executor.CommandSummaryMap())

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

	executorMap := map[string]executor.Executor{
		executor.CommandAccept: executor.NewAcceptOrRejectExecutor(
			topLevelFlags,
			executor.CommandAccept,
			executor.CommandSummaryLookup(executor.CommandAccept),
		),
		executor.CommandAdd: executor.NewAddExecutor(
			topLevelFlags,
			executor.CommandAdd,
			executor.CommandSummaryLookup(executor.CommandAdd),
		),
		executor.CommandBlock: executor.NewBlockOrUnblockExecutor(
			topLevelFlags,
			executor.CommandBlock,
			executor.CommandSummaryLookup(executor.CommandBlock),
		),
		executor.CommandCreate: executor.NewCreateExecutor(
			topLevelFlags,
			executor.CommandCreate,
			executor.CommandSummaryLookup(executor.CommandCreate),
		),
		executor.CommandDelete: executor.NewDeleteExecutor(
			topLevelFlags,
			executor.CommandDelete,
			executor.CommandSummaryLookup(executor.CommandDelete),
		),
		executor.CommandEdit: executor.NewEditExecutor(
			topLevelFlags,
			executor.CommandEdit,
			executor.CommandSummaryLookup(executor.CommandEdit),
		),
		executor.CommandFollow: executor.NewFollowOrUnfollowExecutor(
			topLevelFlags,
			executor.CommandFollow,
			executor.CommandSummaryLookup(executor.CommandFollow),
		),
		executor.CommandLogin: executor.NewLoginExecutor(
			topLevelFlags,
			executor.CommandLogin,
			executor.CommandSummaryLookup(executor.CommandLogin),
		),
		executor.CommandReject: executor.NewAcceptOrRejectExecutor(
			topLevelFlags,
			executor.CommandReject,
			executor.CommandSummaryLookup(executor.CommandReject),
		),
		executor.CommandRemove: executor.NewRemoveExecutor(
			topLevelFlags,
			executor.CommandRemove,
			executor.CommandSummaryLookup(executor.CommandRemove),
		),
		executor.CommandSwitch: executor.NewSwitchExecutor(
			topLevelFlags,
			executor.CommandSwitch,
			executor.CommandSummaryLookup(executor.CommandSwitch),
		),
		executor.CommandUnfollow: executor.NewFollowOrUnfollowExecutor(
			topLevelFlags,
			executor.CommandUnfollow,
			executor.CommandSummaryLookup(executor.CommandUnfollow),
		),
		executor.CommandUnblock: executor.NewBlockOrUnblockExecutor(
			topLevelFlags,
			executor.CommandUnblock,
			executor.CommandSummaryLookup(executor.CommandUnblock),
		),
		executor.CommandShow: executor.NewShowExecutor(
			topLevelFlags,
			executor.CommandShow,
			executor.CommandSummaryLookup(executor.CommandShow),
		),
		executor.CommandVersion: executor.NewVersionExecutor(
			executor.CommandVersion,
			executor.CommandSummaryLookup(executor.CommandVersion),
			binaryVersion,
			buildTime,
			goVersion,
			gitCommit,
		),
		executor.CommandWhoami: executor.NewWhoAmIExecutor(
			topLevelFlags,
			executor.CommandWhoami,
			executor.CommandSummaryLookup(executor.CommandWhoami),
		),
	}

	exe, ok := executorMap[command]
	if !ok {
		flag.Usage()

		return executor.UnknownCommandError{Command: command}
	}

	if err := executor.Execute(exe, args); err != nil {
		return fmt.Errorf("(%s) %w", command, err)
	}

	return nil
}
