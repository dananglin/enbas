package executor

import (
	"fmt"

	"codeflow.dananglin.me.uk/apollo/enbas/internal/client"
)

func (b *UnblockExecutor) Execute() error {
	funcMap := map[string]func(*client.Client) error{
		resourceAccount: b.unblockAccount,
	}

	doFunc, ok := funcMap[b.resourceType]
	if !ok {
		return UnsupportedTypeError{resourceType: b.resourceType}
	}

	gtsClient, err := client.NewClientFromFile(b.config.CredentialsFile)
	if err != nil {
		return fmt.Errorf("unable to create the GoToSocial client: %w", err)
	}

	return doFunc(gtsClient)
}

func (b *UnblockExecutor) unblockAccount(gtsClient *client.Client) error {
	accountID, err := getAccountID(gtsClient, false, b.accountName)
	if err != nil {
		return fmt.Errorf("received an error while getting the account ID: %w", err)
	}

	if err := gtsClient.UnblockAccount(accountID); err != nil {
		return fmt.Errorf("unable to unblock the account: %w", err)
	}

	b.printer.PrintSuccess("Successfully unblocked the account.")

	return nil
}
