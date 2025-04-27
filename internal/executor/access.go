package executor

import (
	"fmt"
	"strings"

	"codeflow.dananglin.me.uk/apollo/enbas/internal/cli"
	"codeflow.dananglin.me.uk/apollo/enbas/internal/command"
	"codeflow.dananglin.me.uk/apollo/enbas/internal/config"
	internalFlag "codeflow.dananglin.me.uk/apollo/enbas/internal/flag"
	"codeflow.dananglin.me.uk/apollo/enbas/internal/gtsclient"
	"codeflow.dananglin.me.uk/apollo/enbas/internal/model"
	"codeflow.dananglin.me.uk/apollo/enbas/internal/printer"
	"codeflow.dananglin.me.uk/apollo/enbas/internal/server"
	"codeflow.dananglin.me.uk/apollo/enbas/internal/utilities"
)

// accessFunc is the function for the access target for
// managing the access to the GoToSocial instance.
func accessFunc(
	opts topLevelOpts,
	cmd command.Command,
) error {
	// Load the configuration from file.
	cfg, err := config.NewConfigFromFile(opts.configDir)
	if err != nil {
		return fmt.Errorf("unable to load configuration: %w", err)
	}

	// Create the print settings.
	printSettings := printer.NewSettings(
		opts.noColor,
		"",
		cfg.LineWrapMaxWidth,
	)

	switch cmd.Action {
	case cli.ActionCreate:
		return accessCreate(
			cfg,
			printSettings,
			opts.configDir,
			cmd.FocusedTargetFlags,
		)
	case cli.ActionSwitch:
		return accessSwitch(
			cfg,
			printSettings,
			opts.configDir,
			cmd.RelatedTarget,
			cmd.RelatedTargetFlags,
		)
	case cli.ActionVerify:
		return accessVerify(
			cfg,
			printSettings,
			opts.configDir,
		)
	default:
		return unsupportedActionError{action: cmd.Action, target: cli.TargetAccess}
	}
}

func accessCreate(
	cfg config.Config,
	printSettings printer.Settings,
	configDir string,
	flags []string,
) error {
	var (
		url    string
		scopes = internalFlag.NewStringSliceValue()
		err    error
	)

	// Parse the remaining flags.
	if err := cli.ParseAccessCreateFlags(
		&scopes,
		&url,
		flags,
	); err != nil {
		return err
	}

	if url == "" {
		return loginNoInstanceError{}
	}

	if !strings.HasPrefix(url, "https") || !strings.HasPrefix(url, "http") {
		url = "https://" + url
	}

	for strings.HasSuffix(url, "/") {
		url = url[:len(url)-1]
	}

	client, err := server.Connect(cfg.Server, configDir)
	if err != nil {
		return fmt.Errorf("error creating the client for the daemon process: %w", err)
	}
	defer client.Close()

	// Update the GTSClient's auth details for the registration process.
	auth := config.Credentials{
		Instance:     url,
		ClientID:     "",
		ClientSecret: "",
		AccessToken:  "",
	}

	if err := client.Call(
		"GTSClient.UpdateAuthentication",
		auth,
		nil,
	); err != nil {
		return fmt.Errorf("error updating the GTSClient's authentication details: %w", err)
	}

	if err := client.Call(
		"GTSClient.RegisterApp",
		scopes,
		nil,
	); err != nil {
		return fmt.Errorf("error registering the application: %w", err)
	}

	var consentPageURL string

	if err := client.Call(
		"GTSClient.AuthCodeURL",
		scopes,
		&consentPageURL,
	); err != nil {
		return fmt.Errorf("error retrieving the URL of the consent page: %w", err)
	}

	_ = utilities.OpenLink(cfg.Integrations.Browser, consentPageURL)

	messageFmt := `
You'll need to sign into your GoToSocial's consent page in order to generate the out-of-band token to continue with the application's login process.
Your browser may have opened the link to the consent page already. If not, please copy and paste the link below to your browser:

%s

Once you have the code, please copy and paste it below.
Out-of-band token: `

	printer.PrintInfo(fmt.Sprintf(messageFmt, consentPageURL))

	var code string

	if _, err := fmt.Scanln(&code); err != nil {
		return fmt.Errorf("failed to read access code: %w", err)
	}

	if err := client.Call(
		"GTSClient.UpdateAccessToken",
		code,
		&auth,
	); err != nil {
		return fmt.Errorf("error updating the client's access token: %w", err)
	}

	var account model.Account
	if err := client.Call(
		"GTSClient.GetMyAccount",
		gtsclient.NoRPCArgs{},
		&account,
	); err != nil {
		return fmt.Errorf("error verifying the credentials: %w", err)
	}

	loginName, err := config.SaveCredentials(
		cfg.CredentialsFile,
		account.Username,
		auth,
	)
	if err != nil {
		return fmt.Errorf("error saving the authentication details: %w", err)
	}

	printer.PrintSuccess(printSettings, "You have successfully logged in as "+loginName+".")

	return nil
}

func accessVerify(
	cfg config.Config,
	printSettings printer.Settings,
	configDir string,
) error {
	client, err := server.Connect(cfg.Server, configDir)
	if err != nil {
		return fmt.Errorf("error creating the client for the daemon process: %w", err)
	}
	defer client.Close()

	var account model.Account
	if err := client.Call("GTSClient.GetMyAccount", gtsclient.NoRPCArgs{}, &account); err != nil {
		return fmt.Errorf("error getting your account information: %w", err)
	}

	var instanceURL string
	if err := client.Call("GTSClient.GetInstanceURL", gtsclient.NoRPCArgs{}, &instanceURL); err != nil {
		return fmt.Errorf("error getting the instance URL: %w", err)
	}

	printer.PrintSuccess(
		printSettings,
		"You are logged in as '"+account.Username+"@"+utilities.GetFQDN(instanceURL)+"'.",
	)

	return nil
}

func accessSwitch(
	cfg config.Config,
	printSettings printer.Settings,
	configDir string,
	relatedTarget string,
	relatedTargetFlags []string,
) error {
	switch relatedTarget {
	case cli.TargetAccount:
		return accessSwitchToAccount(
			cfg,
			printSettings,
			configDir,
			relatedTargetFlags,
		)
	default:
		return unsupportedTargetToTargetError{
			action:        cli.ActionSwitch,
			focusedTarget: cli.TargetAccess,
			preposition: cli.TargetActionPreposition(
				cli.TargetAccess,
				cli.ActionSwitch,
			),
			relatedTarget: relatedTarget,
		}
	}
}

func accessSwitchToAccount(
	cfg config.Config,
	printSettings printer.Settings,
	configDir string,
	flags []string,
) error {
	var accountName string

	// Parse the remaining flags
	if err := cli.ParseAccessSwitchToAccountFlags(
		&accountName,
		flags,
	); err != nil {
		return err
	}

	if accountName == "" {
		return missingAccountNameError{action: "switch the access to"}
	}

	// Create the client to the backend enbas server
	client, err := server.Connect(cfg.Server, configDir)
	if err != nil {
		return fmt.Errorf("error creating the client for the daemon process: %w", err)
	}
	defer client.Close()

	creds, err := config.NewCredentialsConfigFromFile(cfg.CredentialsFile)
	if err != nil {
		return fmt.Errorf("error retrieving the credentials: %w", err)
	}

	auth, ok := creds.Credentials[accountName]
	if !ok {
		return missingAccountInCredentialsError{}
	}

	if err := client.Call(
		"GTSClient.UpdateAuthentication",
		auth,
		nil,
	); err != nil {
		return fmt.Errorf("error updating the authentication details: %w", err)
	}

	if err := config.UpdateCurrentAccount(accountName, cfg.CredentialsFile); err != nil {
		return fmt.Errorf("error updating the credentials config file: %w", err)
	}

	printer.PrintSuccess(printSettings, "The current account is now set to '"+accountName+"'.")

	return nil
}
