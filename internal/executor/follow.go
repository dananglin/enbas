package executor

import (
	"fmt"

	"codeflow.dananglin.me.uk/apollo/enbas/internal/client"
)

func (f *FollowExecutor) Execute() error {
	funcMap := map[string]func(*client.Client) error{
		resourceAccount: f.followAccount,
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

func (f *FollowExecutor) followAccount(gtsClient *client.Client) error {
	accountID, err := getAccountID(gtsClient, false, f.accountName, f.config.CredentialsFile)
	if err != nil {
		return fmt.Errorf("received an error while getting the account ID: %w", err)
	}

	form := client.FollowAccountForm{
		AccountID:   accountID,
		ShowReposts: f.showReposts,
		Notify:      f.notify,
	}

	if err := gtsClient.FollowAccount(form); err != nil {
		return fmt.Errorf("unable to follow the account: %w", err)
	}

	f.printer.PrintSuccess("Successfully sent the follow request.")

	return nil
}
