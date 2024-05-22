package main

import (
	"flag"
	"fmt"
	"os"
)

const (
	accountNameFlag             = "account-name"
	addToFlag                   = "to"
	contentFlag                 = "content"
	instanceFlag                = "instance"
	limitFlag                   = "limit"
	listIDFlag                  = "list-id"
	listTitleFlag               = "list-title"
	listRepliesPolicyFlag       = "list-replies-policy"
	myAccountFlag               = "my-account"
	notifyFlag                  = "notify"
	removeFromFlag              = "from"
	resourceTypeFlag            = "type"
	showAccountRelationshipFlag = "show-account-relationship"
	showUserPreferencesFlag     = "show-preferences"
	showRepostsFlag             = "show-reposts"
	statusIDFlag                = "status-id"
	tagFlag                     = "tag"
	timelineCategoryFlag        = "timeline-category"
	toAccountFlag               = "to-account"
)

const (
	accountResource   = "account"
	blockedResource   = "blocked"
	followersResource = "followers"
	followingResource = "following"
	instanceResource  = "instance"
	listResource      = "list"
	noteResource      = "note"
	statusResource    = "status"
	timelineResource  = "timeline"
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
		add            string = "add"
		remove         string = "remove"
		follow         string = "follow"
		unfollow       string = "unfollow"
		block          string = "block"
		unblock        string = "unblock"
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
		add:            "add a resource to another resource",
		remove:         "remove a resource from another resource",
		follow:         "follow a resource (e.g. an account)",
		unfollow:       "unfollow a resource (e.g. an account)",
		block:          "block a resource (e.g. an account)",
		unblock:        "unblock a resource (e.g. an account)",
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
	case add:
		executor = newAddCommand(add, summaries[add])
	case remove:
		executor = newRemoveCommand(remove, summaries[remove])
	case follow:
		executor = newFollowCommand(follow, summaries[follow], false)
	case unfollow:
		executor = newFollowCommand(unfollow, summaries[unfollow], true)
	case block:
		executor = newBlockCommand(block, summaries[block], false)
	case unblock:
		executor = newBlockCommand(unblock, summaries[unblock], true)
	default:
		flag.Usage()

		return unknownSubcommandError{subcommand}
	}

	if err := executor.Parse(args); err != nil {
		return fmt.Errorf("unable to parse the command line flags; %w", err)
	}

	if err := executor.Execute(); err != nil {
		return fmt.Errorf("received an error after executing %q; %w", executor.Name(), err)
	}

	return nil
}
