package executor

import (
	"fmt"

	"codeflow.dananglin.me.uk/apollo/enbas/internal/client"
)

func (a *AcceptExecutor) Execute() error {
	funcMap := map[string]func(*client.Client) error{
		resourceFollowRequest: a.acceptFollowRequest,
	}

	doFunc, ok := funcMap[a.resourceType]
	if !ok {
		return UnsupportedTypeError{resourceType: a.resourceType}
	}

	gtsClient, err := client.NewClientFromFile(a.config.CredentialsFile)
	if err != nil {
		return fmt.Errorf("unable to create the GoToSocial client: %w", err)
	}

	return doFunc(gtsClient)
}

func (a *AcceptExecutor) acceptFollowRequest(gtsClient *client.Client) error {
	accountID, err := getAccountID(gtsClient, false, a.accountName)
	if err != nil {
		return fmt.Errorf("received an error while getting the account ID: %w", err)
	}

	if err := gtsClient.AcceptFollowRequest(accountID); err != nil {
		return fmt.Errorf("unable to accept the follow request: %w", err)
	}

	a.printer.PrintSuccess("Successfully accepted the follow request.")

	return nil
}
