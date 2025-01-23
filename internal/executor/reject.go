package executor

import (
	"fmt"
	"net/rpc"

	"codeflow.dananglin.me.uk/apollo/enbas/internal/server"
)

func (r *RejectExecutor) Execute() error {
	funcMap := map[string]func(*rpc.Client) error{
		resourceFollowRequest: r.rejectFollowRequest,
	}

	doFunc, ok := funcMap[r.resourceType]
	if !ok {
		return UnsupportedTypeError{resourceType: r.resourceType}
	}

	client, err := server.Connect(r.config.Server, r.configDir)
	if err != nil {
		return fmt.Errorf("error creating the client for the daemon process: %w", err)
	}
	defer client.Close()

	return doFunc(client)
}

func (r *RejectExecutor) rejectFollowRequest(client *rpc.Client) error {
	accountID, err := getAccountID(client, false, r.accountName)
	if err != nil {
		return fmt.Errorf("received an error while getting the account ID: %w", err)
	}

	if err := client.Call("GTSClient.RejectFollowRequest", accountID, nil); err != nil {
		return fmt.Errorf("unable to reject the follow request: %w", err)
	}

	r.printer.PrintSuccess("Successfully rejected the follow request.")

	return nil
}
