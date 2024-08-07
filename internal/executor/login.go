package executor

import (
	"flag"
	"fmt"
	"strings"

	"codeflow.dananglin.me.uk/apollo/enbas/internal/client"
	"codeflow.dananglin.me.uk/apollo/enbas/internal/config"
	"codeflow.dananglin.me.uk/apollo/enbas/internal/printer"
	"codeflow.dananglin.me.uk/apollo/enbas/internal/utilities"
)

type LoginExecutor struct {
	*flag.FlagSet

	printer  *printer.Printer
	config   *config.Config
	instance string
}

func NewLoginExecutor(printer *printer.Printer, config *config.Config, name, summary string) *LoginExecutor {
	loginExe := LoginExecutor{
		FlagSet: flag.NewFlagSet(name, flag.ExitOnError),

		printer:  printer,
		config:   config,
		instance: "",
	}

	loginExe.StringVar(&loginExe.instance, flagInstance, "", "Specify the instance that you want to login to.")

	loginExe.Usage = commandUsageFunc(name, summary, loginExe.FlagSet)

	return &loginExe
}

func (l *LoginExecutor) Execute() error {
	var err error

	if l.instance == "" {
		return FlagNotSetError{flagText: flagInstance}
	}

	instance := l.instance

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
		return fmt.Errorf("unable to register the application: %w", err)
	}

	consentPageURL := gtsClient.AuthCodeURL()

	_ = utilities.OpenLink(l.config.Integrations.Browser, consentPageURL)

	var builder strings.Builder

	builder.WriteString("\nYou'll need to sign into your GoToSocial's consent page in order to generate the out-of-band token to continue with the application's login process.")
	builder.WriteString("\nYour browser may have opened the link to the consent page already. If not, please copy and paste the link below to your browser:")
	builder.WriteString("\n\n" + consentPageURL)
	builder.WriteString("\n\n" + "Once you have the code please copy and paste it below.")
	builder.WriteString("\n" + "Out-of-band token: ")

	l.printer.PrintInfo(builder.String())

	var code string

	if _, err := fmt.Scanln(&code); err != nil {
		return fmt.Errorf("failed to read access code: %w", err)
	}

	if err := gtsClient.UpdateToken(code); err != nil {
		return fmt.Errorf("unable to update the client's access token: %w", err)
	}

	account, err := gtsClient.VerifyCredentials()
	if err != nil {
		return fmt.Errorf("unable to verify the credentials: %w", err)
	}

	loginName, err := config.SaveCredentials(l.config.CredentialsFile, account.Username, gtsClient.Authentication)
	if err != nil {
		return fmt.Errorf("unable to save the authentication details: %w", err)
	}

	l.printer.PrintSuccess("You have successfully logged as " + loginName + ".")

	return nil
}
