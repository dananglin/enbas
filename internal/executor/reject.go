package executor

import (
	"fmt"

	"codeflow.dananglin.me.uk/apollo/enbas/internal/client"
)

func (r *RejectExecutor) Execute() error {
	funcMap := map[string]func(*client.Client) error{
		resourceFollowRequest: r.rejectFollowRequest,
	}

	doFunc, ok := funcMap[r.resourceType]
	if !ok {
		return UnsupportedTypeError{resourceType: r.resourceType}
	}

	gtsClient, err := client.NewClientFromFile(r.config.CredentialsFile)
	if err != nil {
		return fmt.Errorf("unable to create the GoToSocial client: %w", err)
	}

	return doFunc(gtsClient)
}

func (r *RejectExecutor) rejectFollowRequest(gtsClient *client.Client) error {
	accountID, err := getAccountID(gtsClient, false, r.accountName, r.config.CredentialsFile)
	if err != nil {
		return fmt.Errorf("received an error while getting the account ID: %w", err)
	}

	if err := gtsClient.RejectFollowRequest(accountID); err != nil {
		return fmt.Errorf("unable to reject the follow request: %w", err)
	}

	r.printer.PrintSuccess("Successfully rejected the follow request.")

	return nil
}
