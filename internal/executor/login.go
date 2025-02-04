package executor

import (
	"fmt"
	"strings"

	"codeflow.dananglin.me.uk/apollo/enbas/internal/config"
	"codeflow.dananglin.me.uk/apollo/enbas/internal/gtsclient"
	"codeflow.dananglin.me.uk/apollo/enbas/internal/model"
	"codeflow.dananglin.me.uk/apollo/enbas/internal/printer"
	"codeflow.dananglin.me.uk/apollo/enbas/internal/server"
	"codeflow.dananglin.me.uk/apollo/enbas/internal/utilities"
)

func (l *LoginExecutor) Execute() error {
	var err error

	if l.instance == "" {
		return Error{"please specify the instance that you want to log into"}
	}

	if !strings.HasPrefix(l.instance, "https") || !strings.HasPrefix(l.instance, "http") {
		l.instance = "https://" + l.instance
	}

	for strings.HasSuffix(l.instance, "/") {
		l.instance = l.instance[:len(l.instance)-1]
	}

	client, err := server.Connect(l.config.Server, l.configDir)
	if err != nil {
		return fmt.Errorf("error creating the client for the daemon process: %w", err)
	}
	defer client.Close()

	// Update the GTSClient's auth details for the registration process.
	auth := config.Credentials{
		Instance:     l.instance,
		ClientID:     "",
		ClientSecret: "",
		AccessToken:  "",
	}

	if err := client.Call("GTSClient.UpdateAuthentication", auth, nil); err != nil {
		return fmt.Errorf("error updating the GTSClient's authentication details: %w", err)
	}

	if err := client.Call("GTSClient.Register", gtsclient.NoRPCArgs{}, nil); err != nil {
		return fmt.Errorf("unable to register the application: %w", err)
	}

	var consentPageURL string

	if err := client.Call("GTSClient.AuthCodeURL", gtsclient.NoRPCArgs{}, &consentPageURL); err != nil {
		return fmt.Errorf("unable to get the URL of the consent page: %w", err)
	}

	_ = utilities.OpenLink(l.config.Integrations.Browser, consentPageURL)

	var builder strings.Builder

	builder.WriteString("\nYou'll need to sign into your GoToSocial's consent page in order to generate the out-of-band token to continue with the application's login process.")
	builder.WriteString("\nYour browser may have opened the link to the consent page already. If not, please copy and paste the link below to your browser:")
	builder.WriteString("\n\n" + consentPageURL)
	builder.WriteString("\n\n" + "Once you have the code please copy and paste it below.")
	builder.WriteString("\n" + "Out-of-band token: ")

	printer.PrintInfo(builder.String())

	var code string

	if _, err := fmt.Scanln(&code); err != nil {
		return fmt.Errorf("failed to read access code: %w", err)
	}

	if err := client.Call("GTSClient.UpdateToken", code, &auth); err != nil {
		return fmt.Errorf("unable to update the client's access token: %w", err)
	}

	var account model.Account
	if err := client.Call("GTSClient.VerifyCredentials", gtsclient.NoRPCArgs{}, &account); err != nil {
		return fmt.Errorf("unable to verify the credentials: %w", err)
	}

	loginName, err := config.SaveCredentials(l.config.CredentialsFile, account.Username, auth)
	if err != nil {
		return fmt.Errorf("unable to save the authentication details: %w", err)
	}

	printer.PrintSuccess(l.printSettings, "You have successfully logged in as "+loginName+".")

	return nil
}
