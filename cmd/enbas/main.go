// SPDX-FileCopyrightText: 2024 Dan Anglin <d.n.i.anglin@gmail.com>
//
// SPDX-License-Identifier: GPL-3.0-or-later

package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"

	"codeflow.dananglin.me.uk/apollo/enbas/internal/config"
	"codeflow.dananglin.me.uk/apollo/enbas/internal/executor"
	"codeflow.dananglin.me.uk/apollo/enbas/internal/printer"
)

var (
	binaryVersion string //nolint:gochecknoglobals
	buildTime     string //nolint:gochecknoglobals
	goVersion     string //nolint:gochecknoglobals
	gitCommit     string //nolint:gochecknoglobals
)

func main() {
	if err := run(); err != nil {
		os.Exit(1)
	}
}

func run() error {
	var (
		configDir string
		noColor   *bool
	)

	flag.StringVar(&configDir, "config-dir", "", "Specify your config directory")
	flag.BoolFunc("no-color", "Disable ANSI colour output when displaying text on screen", func(value string) error {
		boolVal, err := strconv.ParseBool(value)
		if err != nil {
			return fmt.Errorf("unable to parse %q as a boolean: %w", value, err)
		}

		noColor = new(bool)
		*noColor = boolVal

		return nil
	})

	flag.Usage = usageFunc(executor.CommandSummaryMap())

	flag.Parse()

	if flag.NArg() < 1 {
		flag.Usage()

		return nil
	}

	// If NoColor is still unspecified,
	// check to see if the NO_COLOR environment variable is set
	if noColor == nil {
		noColor = new(bool)
		if os.Getenv("NO_COLOR") != "" {
			*noColor = true
		} else {
			*noColor = false
		}
	}

	command := flag.Arg(0)
	args := flag.Args()[1:]

	var (
		enbasConfig  *config.Config
		enbasPrinter *printer.Printer
		err          error
	)

	switch command {
	case executor.CommandInit, executor.CommandVersion:
		enbasPrinter = printer.NewPrinter(*noColor, "", 0)
	default:
		enbasConfig, err = config.NewConfigFromFile(configDir)
		if err != nil {
			enbasPrinter = printer.NewPrinter(*noColor, "", 0)
			enbasPrinter.PrintFailure("unable to load the configuration: " + err.Error() + ".")

			return err
		}

		enbasPrinter = printer.NewPrinter(*noColor, enbasConfig.Integrations.Pager, enbasConfig.LineWrapMaxWidth)
	}

	executorMap := map[string]executor.Executor{
		executor.CommandAccept: executor.NewAcceptOrRejectExecutor(
			enbasPrinter,
			enbasConfig,
			executor.CommandAccept,
			executor.CommandSummaryLookup(executor.CommandAccept),
		),
		executor.CommandAdd: executor.NewAddExecutor(
			enbasPrinter,
			enbasConfig,
			executor.CommandAdd,
			executor.CommandSummaryLookup(executor.CommandAdd),
		),
		executor.CommandBlock: executor.NewBlockOrUnblockExecutor(
			enbasPrinter,
			enbasConfig,
			executor.CommandBlock,
			executor.CommandSummaryLookup(executor.CommandBlock),
		),
		executor.CommandCreate: executor.NewCreateExecutor(
			enbasPrinter,
			enbasConfig,
			executor.CommandCreate,
			executor.CommandSummaryLookup(executor.CommandCreate),
		),
		executor.CommandDelete: executor.NewDeleteExecutor(
			enbasPrinter,
			enbasConfig,
			executor.CommandDelete,
			executor.CommandSummaryLookup(executor.CommandDelete),
		),
		executor.CommandEdit: executor.NewEditExecutor(
			enbasPrinter,
			enbasConfig,
			executor.CommandEdit,
			executor.CommandSummaryLookup(executor.CommandEdit),
		),
		executor.CommandFollow: executor.NewFollowOrUnfollowExecutor(
			enbasPrinter,
			enbasConfig,
			executor.CommandFollow,
			executor.CommandSummaryLookup(executor.CommandFollow),
		),
		executor.CommandInit: executor.NewInitExecutor(
			enbasPrinter,
			configDir,
			executor.CommandInit,
			executor.CommandSummaryLookup(executor.CommandInit),
		),
		executor.CommandLogin: executor.NewLoginExecutor(
			enbasPrinter,
			enbasConfig,
			executor.CommandLogin,
			executor.CommandSummaryLookup(executor.CommandLogin),
		),
		executor.CommandMute: executor.NewMuteOrUnmuteExecutor(
			enbasPrinter,
			enbasConfig,
			executor.CommandMute,
			executor.CommandSummaryLookup(executor.CommandMute),
		),
		executor.CommandReject: executor.NewAcceptOrRejectExecutor(
			enbasPrinter,
			enbasConfig,
			executor.CommandReject,
			executor.CommandSummaryLookup(executor.CommandReject),
		),
		executor.CommandRemove: executor.NewRemoveExecutor(
			enbasPrinter,
			enbasConfig,
			executor.CommandRemove,
			executor.CommandSummaryLookup(executor.CommandRemove),
		),
		executor.CommandSwitch: executor.NewSwitchExecutor(
			enbasPrinter,
			enbasConfig,
			executor.CommandSwitch,
			executor.CommandSummaryLookup(executor.CommandSwitch),
		),
		executor.CommandUnfollow: executor.NewFollowOrUnfollowExecutor(
			enbasPrinter,
			enbasConfig,
			executor.CommandUnfollow,
			executor.CommandSummaryLookup(executor.CommandUnfollow),
		),
		executor.CommandUnmute: executor.NewMuteOrUnmuteExecutor(
			enbasPrinter,
			enbasConfig,
			executor.CommandUnmute,
			executor.CommandSummaryLookup(executor.CommandUnmute),
		),
		executor.CommandUnblock: executor.NewBlockOrUnblockExecutor(
			enbasPrinter,
			enbasConfig,
			executor.CommandUnblock,
			executor.CommandSummaryLookup(executor.CommandUnblock),
		),
		executor.CommandShow: executor.NewShowExecutor(
			enbasPrinter,
			enbasConfig,
			executor.CommandShow,
			executor.CommandSummaryLookup(executor.CommandShow),
		),
		executor.CommandVersion: executor.NewVersionExecutor(
			enbasPrinter,
			executor.CommandVersion,
			executor.CommandSummaryLookup(executor.CommandVersion),
			binaryVersion,
			buildTime,
			goVersion,
			gitCommit,
		),
		executor.CommandWhoami: executor.NewWhoAmIExecutor(
			enbasPrinter,
			enbasConfig,
			executor.CommandWhoami,
			executor.CommandSummaryLookup(executor.CommandWhoami),
		),
	}

	exe, ok := executorMap[command]
	if !ok {
		err = executor.UnknownCommandError{Command: command}

		enbasPrinter.PrintFailure(err.Error() + ".")
		flag.Usage()

		return err
	}

	if err = executor.Execute(exe, args); err != nil {
		enbasPrinter.PrintFailure("(" + command + ") " + err.Error() + ".")

		return err
	}

	return nil
}
