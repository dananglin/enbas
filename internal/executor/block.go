package executor

import (
	"fmt"

	"codeflow.dananglin.me.uk/apollo/enbas/internal/client"
)

func (b *BlockExecutor) Execute() error {
	funcMap := map[string]func(*client.Client) error{
		resourceAccount: b.blockAccount,
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

func (b *BlockExecutor) blockAccount(gtsClient *client.Client) error {
	accountID, err := getAccountID(gtsClient, false, b.accountName, b.config.CredentialsFile)
	if err != nil {
		return fmt.Errorf("received an error while getting the account ID: %w", err)
	}

	if err := gtsClient.BlockAccount(accountID); err != nil {
		return fmt.Errorf("unable to block the account: %w", err)
	}

	b.printer.PrintSuccess("Successfully blocked the account.")

	return nil
}
