package main

import (
	"errors"
	"flag"
	"fmt"
	"strings"
	"unicode"

	"codeflow.dananglin.me.uk/apollo/enbas/internal/client"
	"codeflow.dananglin.me.uk/apollo/enbas/internal/config"
	"golang.org/x/net/html"
)

var instanceDetailsFormat = `INSTANCE:
  %s - %s

DOMAIN:
  %s

VERSION:
  Running GoToSocial %s

CONTACT:
  name: %s
  username: %s
  email: %s
`

var accountDetailsFormat = `
%s (@%s)

ACCOUNT ID:
  %s

JOINED ON:
  %s

STATS:
  Followers: %d
  Following: %d
  Statuses: %d

BIOGRAPHY:
  %s

METADATA: %s

ACCOUNT URL:
  %s
`

type showCommand struct {
	*flag.FlagSet
	targetType string
	account    string
	myAccount  bool
}

func newShowCommand(name, summary string) *showCommand {
	command := showCommand{
		FlagSet:    flag.NewFlagSet(name, flag.ExitOnError),
		targetType: "",
		myAccount:  false,
	}

	command.StringVar(&command.targetType, "type", "", "specify the type of resource to display")
	command.StringVar(&command.account, "account", "", "specify the account URI to lookup")
	command.BoolVar(&command.myAccount, "my-account", false, "set to true to lookup your account")
	command.Usage = commandUsageFunc(name, summary, command.FlagSet)

	return &command
}

func (c *showCommand) Execute() error {
	gtsClient, err := client.NewClientFromConfig()
	if err != nil {
		return fmt.Errorf("unable to create the GoToSocial client; %w", err)
	}

	funcMap := map[string]func(*client.Client) error{
		"instance": c.showInstance,
		"account":  c.showAccount,
	}

	doFunc, ok := funcMap[c.targetType]
	if !ok {
		return fmt.Errorf("unsupported type %q", c.targetType)
	}

	return doFunc(gtsClient)
}

func (c *showCommand) showInstance(gts *client.Client) error {
	instance, err := gts.GetInstance()
	if err != nil {
		return fmt.Errorf("unable to retrieve the instance details; %w", err)
	}

	fmt.Printf(
		instanceDetailsFormat,
		instance.Title,
		instance.Description,
		instance.Domain,
		instance.Version,
		instance.Contact.Account.DisplayName,
		instance.Contact.Account.Username,
		instance.Contact.Email,
	)

	return nil
}

func (c *showCommand) showAccount(gts *client.Client) error {
	var accountURI string

	if c.myAccount {
		authConfig, err := config.NewAuthenticationConfigFromFile()
		if err != nil {
			return fmt.Errorf("unable to retrieve the authentication configuration; %w", err)
		}

		accountURI = authConfig.CurrentAccount
	} else {
		if c.account == "" {
			return errors.New("the account flag is not set")
		}

		accountURI = c.account
	}

	account, err := gts.GetAccount(accountURI)
	if err != nil {
		return fmt.Errorf("unable to retrieve the account details; %w", err)
	}

	metadata := ""

	for _, field := range account.Fields {
		metadata += fmt.Sprintf("\n  %s: %s", field.Name, stripHTMLTags(field.Value))
	}

	fmt.Printf(
		accountDetailsFormat,
		account.DisplayName,
		account.Username,
		account.ID,
		account.CreatedAt.Format("02 Jan 2006"),
		account.FollowersCount,
		account.FollowingCount,
		account.StatusCount,
		wrapLine(stripHTMLTags(account.Note), "\n  ", 80),
		metadata,
		account.URL,
	)

	return nil
}

func stripHTMLTags(text string) string {
	token := html.NewTokenizer(strings.NewReader(text))

	var builder strings.Builder

	for {
		tt := token.Next()
		switch tt {
		case html.ErrorToken:
			return builder.String()
		case html.TextToken:
			builder.WriteString(token.Token().Data + " ")
		}
	}
}

func wrapLine(line, separator string, charLimit int) string {
	if len(line) <= charLimit {
		return line
	}

	leftcursor, rightcursor := 0, 0

	var builder strings.Builder

	for rightcursor < (len(line) - charLimit) {
		rightcursor += charLimit
		for !unicode.IsSpace(rune(line[rightcursor-1])) {
			rightcursor--
		}
		builder.WriteString(line[leftcursor:rightcursor] + separator)
		leftcursor = rightcursor
	}

	builder.WriteString(line[rightcursor:])

	return builder.String()
}
