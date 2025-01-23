package executor

import (
	"fmt"
	"net/rpc"

	"codeflow.dananglin.me.uk/apollo/enbas/internal/config"
	"codeflow.dananglin.me.uk/apollo/enbas/internal/server"
)

func (s *SwitchExecutor) Execute() error {
	funcMap := map[string]func(*rpc.Client) error{
		resourceAccount: s.switchToAccount,
	}

	doFunc, ok := funcMap[s.to]
	if !ok {
		return UnsupportedTypeError{resourceType: s.to}
	}

	client, err := server.Connect(s.config.Server, s.configDir)
	if err != nil {
		return fmt.Errorf("error creating the client for the daemon process: %w", err)
	}
	defer client.Close()

	return doFunc(client)
}

func (s *SwitchExecutor) switchToAccount(client *rpc.Client) error {
	if !s.accountName.ExpectedLength(1) {
		return Error{"found an unexpected number of account names: expected 1"}
	}

	creds, err := config.NewCredentialsConfigFromFile(s.config.CredentialsFile)
	if err != nil {
		return fmt.Errorf("error retrieving the credentials: %w", err)
	}

	auth, ok := creds.Credentials[s.accountName[0]]
	if !ok {
		return Error{"the account is not present in the credentials file"}
	}

	if err := client.Call("GTSClient.UpdateAuthentication", auth, nil); err != nil {
		return fmt.Errorf("error updating the authentication details: %w", err)
	}

	if err := config.UpdateCurrentAccount(s.accountName[0], s.config.CredentialsFile); err != nil {
		return fmt.Errorf("error updating the credentials config file: %w", err)
	}

	s.printer.PrintSuccess("The current account is now set to '" + s.accountName[0] + "'.")

	return nil
}
