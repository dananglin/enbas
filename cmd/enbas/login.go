package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os/exec"
	"runtime"
	"strings"

	"codeflow.dananglin.me.uk/apollo/enbas/internal"
	"codeflow.dananglin.me.uk/apollo/enbas/internal/client"
	"codeflow.dananglin.me.uk/apollo/enbas/internal/config"
	"golang.org/x/oauth2"
)

type loginCommand struct {
	*flag.FlagSet
	instance string
}

var (
	errEmptyAccessToken = errors.New("received an empty access token")
	errInstanceNotSet   = errors.New("the instance flag is not set")
)

var consentMessageFormat = `
You'll need to sign into your GoToSocial's consent page in order to generate the out-of-band token to continue with
the application's login process. Your browser may have opened the link to the consent page already. If not, please
copy and paste the link below to your browser:

%s

Once you have the code please copy and paste it below.

`

func newLoginCommand(name, summary string) *loginCommand {
	command := loginCommand{
		FlagSet:  flag.NewFlagSet(name, flag.ExitOnError),
		instance: "",
	}

	command.StringVar(&command.instance, "instance", "", "specify the instance that you want to login to.")

	command.Usage = commandUsageFunc(name, summary, command.FlagSet)

	return &command
}

func (c *loginCommand) Execute() error {
	var err error

	if c.instance == "" {
		return errInstanceNotSet
	}

	instance := c.instance

	if !strings.HasPrefix(instance, "https") || !strings.HasPrefix(instance, "http") {
		instance = "https://" + instance
	}

	for strings.HasSuffix(instance, "/") {
		instance = instance[:len(instance)-1]
	}

	authentication := config.Authentication{
		Instance: instance,
	}

	gtsClient := client.NewClient(authentication)

	if err := gtsClient.Register(); err != nil {
		return fmt.Errorf("unable to register the application; %w", err)
	}

	oauth2Conf := oauth2.Config{
		ClientID:     gtsClient.Authentication.ClientID,
		ClientSecret: gtsClient.Authentication.ClientSecret,
		Scopes:       []string{"read"},
		RedirectURL:  internal.RedirectUri,
		Endpoint: oauth2.Endpoint{
			AuthURL:  gtsClient.Authentication.Instance + "/oauth/authorize",
			TokenURL: gtsClient.Authentication.Instance + "/oauth/token",
		},
	}

	consentPageURL := authCodeURL(oauth2Conf)

	openLink(consentPageURL)

	fmt.Printf(consentMessageFormat, consentPageURL)

	var code string
	fmt.Print("Out-of-band token: ")

	if _, err := fmt.Scanln(&code); err != nil {
		return fmt.Errorf("failed to read access code; %w", err)
	}

	gtsClient.Authentication, err = addAccessToken(gtsClient.Authentication, oauth2Conf, code)
	if err != nil {
		return fmt.Errorf("unable to get the access token; %w", err)
	}

	account, err := gtsClient.VerifyCredentials()
	if err != nil {
		return fmt.Errorf("unable to verify the credentials; %w", err)
	}

	loginName, err := config.SaveAuthentication(account.Username, gtsClient.Authentication)
	if err != nil {
		return fmt.Errorf("unable to save the authentication details; %w", err)
	}

	fmt.Printf("Successfully logged into %s\n", loginName)

	return nil
}

func authCodeURL(oauth2Conf oauth2.Config) string {
	url := oauth2Conf.AuthCodeURL(
		"state",
		oauth2.AccessTypeOffline,
	) + "&client_name=" + internal.ApplicationName

	return url
}

func addAccessToken(authentication config.Authentication, oauth2Conf oauth2.Config, code string) (config.Authentication, error) {
	token, err := oauth2Conf.Exchange(context.Background(), code)
	if err != nil {
		return config.Authentication{}, fmt.Errorf("unable to exchange the code for an access token; %w", err)
	}

	if token == nil || token.AccessToken == "" {
		return config.Authentication{}, errEmptyAccessToken
	}

	authentication.AccessToken = token.AccessToken

	return authentication, nil
}

func openLink(url string) {
	var open string

	if runtime.GOOS == "linux" {
		open = "xdg-open"
	} else {
		return
	}

	command := exec.Command(open, url)

	_ = command.Start()
}
