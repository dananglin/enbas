package executor

import (
	"fmt"

	"codeflow.dananglin.me.uk/apollo/enbas/internal/client"
)

func (f *UnfollowExecutor) Execute() error {
	funcMap := map[string]func(*client.Client) error{
		resourceAccount: f.unfollowAccount,
	}

	doFunc, ok := funcMap[f.resourceType]
	if !ok {
		return UnsupportedTypeError{resourceType: f.resourceType}
	}

	gtsClient, err := client.NewClientFromFile(f.config.CredentialsFile)
	if err != nil {
		return fmt.Errorf("unable to create the GoToSocial client: %w", err)
	}

	return doFunc(gtsClient)
}

func (f *UnfollowExecutor) unfollowAccount(gtsClient *client.Client) error {
	accountID, err := getAccountID(gtsClient, false, f.accountName, f.config.CredentialsFile)
	if err != nil {
		return fmt.Errorf("received an error while getting the account ID: %w", err)
	}

	if err := gtsClient.UnfollowAccount(accountID); err != nil {
		return fmt.Errorf("unable to unfollow the account: %w", err)
	}

	f.printer.PrintSuccess("Successfully unfollowed the account.")

	return nil
}
