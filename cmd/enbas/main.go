package main

import (
	"flag"
	"fmt"
	"os"
)

const (
	accountFlag           = "account"
	instanceFlag          = "instance"
	listIDFlag            = "list-id"
	listTitleFlag         = "list-title"
	listRepliesPolicyFlag = "list-replies-policy"
	myAccountFlag         = "my-account"
	resourceTypeFlag      = "type"
	statusIDFlag          = "status-id"
	tagFlag               = "tag"
	timelineCategoryFlag  = "timeline-category"
	timelineLimitFlag     = "timeline-limit"
	toAccountFlag         = "to-account"
)

const (
	accountResource  = "account"
	instanceResource = "instance"
	listResource     = "list"
	statusResource   = "status"
	timelineResource = "timeline"
)

type Executor interface {
	Name() string
	Parse(args []string) error
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
		login          string = "login"
		version        string = "version"
		showResource   string = "show"
		switchAccount  string = "switch"
		createResource string = "create"
		deleteResource string = "delete"
		updateResource string = "update"
		whoami         string = "whoami"
	)

	summaries := map[string]string{
		login:          "login to an account on GoToSocial",
		version:        "print the application's version and build information",
		showResource:   "print details about a specified resource",
		switchAccount:  "switch to an account",
		createResource: "create a specific resource",
		deleteResource: "delete a specific resource",
		updateResource: "update a specific resource",
		whoami:         "print the account that you are currently logged in to",
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
	case showResource:
		executor = newShowCommand(showResource, summaries[showResource])
	case switchAccount:
		executor = newSwitchCommand(switchAccount, summaries[switchAccount])
	case createResource:
		executor = newCreateCommand(createResource, summaries[createResource])
	case deleteResource:
		executor = newDeleteCommand(deleteResource, summaries[deleteResource])
	case updateResource:
		executor = newUpdateCommand(updateResource, summaries[updateResource])
	case whoami:
		executor = newWhoAmICommand(whoami, summaries[whoami])
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
