package main

import (
	"flag"
	"fmt"
	"strings"

	"codeflow.dananglin.me.uk/apollo/enbas/internal/client"
	"codeflow.dananglin.me.uk/apollo/enbas/internal/config"
	"codeflow.dananglin.me.uk/apollo/enbas/internal/utilities"
)

type loginCommand struct {
	*flag.FlagSet

	topLevelFlags topLevelFlags
	instance string
}

func newLoginCommand(tlf topLevelFlags, name, summary string) *loginCommand {
	command := loginCommand{
		FlagSet:  flag.NewFlagSet(name, flag.ExitOnError),
		topLevelFlags: tlf,
		instance: "",
	}

	command.StringVar(&command.instance, instanceFlag, "", "specify the instance that you want to login to.")

	command.Usage = commandUsageFunc(name, summary, command.FlagSet)

	return &command
}

func (c *loginCommand) Execute() error {
	var err error

	if c.instance == "" {
		return flagNotSetError{flagText: instanceFlag}
	}

	instance := c.instance

	if !strings.HasPrefix(instance, "https") || !strings.HasPrefix(instance, "http") {
		instance = "https://" + instance
	}

	for strings.HasSuffix(instance, "/") {
		instance = instance[:len(instance)-1]
	}

	credentials := config.Credentials{
		Instance: instance,
	}

	gtsClient := client.NewClient(credentials)

	if err := gtsClient.Register(); err != nil {
		return fmt.Errorf("unable to register the application; %w", err)
	}

	consentPageURL := gtsClient.AuthCodeURL()

	utilities.OpenLink(consentPageURL)

	consentMessageFormat := `
You'll need to sign into your GoToSocial's consent page in order to generate the out-of-band token to continue with
the application's login process. Your browser may have opened the link to the consent page already. If not, please
copy and paste the link below to your browser:

%s

Once you have the code please copy and paste it below.

`

	fmt.Printf(consentMessageFormat, consentPageURL)

	var code string
	fmt.Print("Out-of-band token: ")

	if _, err := fmt.Scanln(&code); err != nil {
		return fmt.Errorf("failed to read access code; %w", err)
	}

	if err := gtsClient.UpdateToken(code); err != nil {
		return fmt.Errorf("unable to update the client's access token; %w", err)
	}

	account, err := gtsClient.VerifyCredentials()
	if err != nil {
		return fmt.Errorf("unable to verify the credentials; %w", err)
	}

	loginName, err := config.SaveCredentials(c.topLevelFlags.configDir, account.Username, gtsClient.Authentication)
	if err != nil {
		return fmt.Errorf("unable to save the authentication details; %w", err)
	}

	fmt.Printf("Successfully logged into %s\n", loginName)

	return nil
}
