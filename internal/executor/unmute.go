package executor

import (
	"fmt"

	"codeflow.dananglin.me.uk/apollo/enbas/internal/client"
)

func (m *UnmuteExecutor) Execute() error {
	funcMap := map[string]func(*client.Client) error{
		resourceAccount: m.unmuteAccount,
	}

	doFunc, ok := funcMap[m.resourceType]
	if !ok {
		return UnsupportedTypeError{resourceType: m.resourceType}
	}

	gtsClient, err := client.NewClientFromFile(m.config.CredentialsFile)
	if err != nil {
		return fmt.Errorf("unable to create the GoToSocial client: %w", err)
	}

	return doFunc(gtsClient)
}

func (m *UnmuteExecutor) unmuteAccount(gtsClient *client.Client) error {
	accountID, err := getAccountID(gtsClient, false, m.accountName, m.config.CredentialsFile)
	if err != nil {
		return fmt.Errorf("received an error while getting the account ID: %w", err)
	}

	if err := gtsClient.UnmuteAccount(accountID); err != nil {
		return fmt.Errorf("unable to unmute the account: %w", err)
	}

	m.printer.PrintSuccess("Successfully unmuted the account.")

	return nil
}