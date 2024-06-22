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
	"codeflow.dananglin.me.uk/apollo/enbas/internal/printer"
)

var (
	binaryVersion string
	buildTime     string
	goVersion     string
	gitCommit     string
)

func main() {
	if err := run(); err != nil {
		os.Exit(1)
	}
}

func run() error {
	var (
		configDir        string
		cacheDir         string
		pager            string
		imageViewer      string
		videoPlayer      string
		maxTerminalWidth int
		noColor          *bool
	)

	flag.StringVar(&configDir, "config-dir", "", "Specify your config directory")
	flag.StringVar(&cacheDir, "cache-dir", "", "Specify your cache directory")
	flag.StringVar(&pager, "pager", "", "Specify your preferred pager to page through long outputs. This is disabled by default.")
	flag.StringVar(&imageViewer, "image-viewer", "", "Specify your favourite image viewer.")
	flag.StringVar(&videoPlayer, "video-player", "", "Specify your favourite video player.")
	flag.IntVar(&maxTerminalWidth, "max-terminal-width", 80, "Specify the maximum terminal width when displaying resources on screen.")

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

	// If NoColor is still unspecified, check to see if the NO_COLOR environment variable is set
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

	printer := printer.NewPrinter(*noColor, pager, maxTerminalWidth)

	executorMap := map[string]executor.Executor{
		executor.CommandAccept: executor.NewAcceptOrRejectExecutor(
			printer,
			configDir,
			executor.CommandAccept,
			executor.CommandSummaryLookup(executor.CommandAccept),
		),
		executor.CommandAdd: executor.NewAddExecutor(
			printer,
			configDir,
			executor.CommandAdd,
			executor.CommandSummaryLookup(executor.CommandAdd),
		),
		executor.CommandBlock: executor.NewBlockOrUnblockExecutor(
			printer,
			configDir,
			executor.CommandBlock,
			executor.CommandSummaryLookup(executor.CommandBlock),
		),
		executor.CommandCreate: executor.NewCreateExecutor(
			printer,
			configDir,
			executor.CommandCreate,
			executor.CommandSummaryLookup(executor.CommandCreate),
		),
		executor.CommandDelete: executor.NewDeleteExecutor(
			printer,
			configDir,
			executor.CommandDelete,
			executor.CommandSummaryLookup(executor.CommandDelete),
		),
		executor.CommandEdit: executor.NewEditExecutor(
			printer,
			configDir,
			executor.CommandEdit,
			executor.CommandSummaryLookup(executor.CommandEdit),
		),
		executor.CommandFollow: executor.NewFollowOrUnfollowExecutor(
			printer,
			configDir,
			executor.CommandFollow,
			executor.CommandSummaryLookup(executor.CommandFollow),
		),
		executor.CommandLogin: executor.NewLoginExecutor(
			printer,
			configDir,
			executor.CommandLogin,
			executor.CommandSummaryLookup(executor.CommandLogin),
		),
		executor.CommandMute: executor.NewMuteOrUnmuteExecutor(
			printer,
			configDir,
			executor.CommandMute,
			executor.CommandSummaryLookup(executor.CommandMute),
		),
		executor.CommandReject: executor.NewAcceptOrRejectExecutor(
			printer,
			configDir,
			executor.CommandReject,
			executor.CommandSummaryLookup(executor.CommandReject),
		),
		executor.CommandRemove: executor.NewRemoveExecutor(
			printer,
			configDir,
			executor.CommandRemove,
			executor.CommandSummaryLookup(executor.CommandRemove),
		),
		executor.CommandSwitch: executor.NewSwitchExecutor(
			printer,
			configDir,
			executor.CommandSwitch,
			executor.CommandSummaryLookup(executor.CommandSwitch),
		),
		executor.CommandUnfollow: executor.NewFollowOrUnfollowExecutor(
			printer,
			configDir,
			executor.CommandUnfollow,
			executor.CommandSummaryLookup(executor.CommandUnfollow),
		),
		executor.CommandUnmute: executor.NewMuteOrUnmuteExecutor(
			printer,
			configDir,
			executor.CommandUnmute,
			executor.CommandSummaryLookup(executor.CommandUnmute),
		),
		executor.CommandUnblock: executor.NewBlockOrUnblockExecutor(
			printer,
			configDir,
			executor.CommandUnblock,
			executor.CommandSummaryLookup(executor.CommandUnblock),
		),
		executor.CommandShow: executor.NewShowExecutor(
			printer,
			configDir,
			cacheDir,
			imageViewer,
			videoPlayer,
			executor.CommandShow,
			executor.CommandSummaryLookup(executor.CommandShow),
		),
		executor.CommandVersion: executor.NewVersionExecutor(
			printer,
			executor.CommandVersion,
			executor.CommandSummaryLookup(executor.CommandVersion),
			binaryVersion,
			buildTime,
			goVersion,
			gitCommit,
		),
		executor.CommandWhoami: executor.NewWhoAmIExecutor(
			printer,
			configDir,
			executor.CommandWhoami,
			executor.CommandSummaryLookup(executor.CommandWhoami),
		),
	}

	exe, ok := executorMap[command]
	if !ok {
		err := executor.UnknownCommandError{Command: command}

		printer.PrintFailure(err.Error() + ".")
		flag.Usage()

		return err
	}

	if err := executor.Execute(exe, args); err != nil {
		printer.PrintFailure("(" + command + ") " + err.Error() + ".")

		return err
	}

	return nil
}
