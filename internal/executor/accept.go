package executor

import (
	"fmt"
	"net/rpc"

	"codeflow.dananglin.me.uk/apollo/enbas/internal/printer"
	"codeflow.dananglin.me.uk/apollo/enbas/internal/server"
)

func (a *AcceptExecutor) Execute() error {
	funcMap := map[string]func(*rpc.Client) error{
		resourceFollowRequest: a.acceptFollowRequest,
	}

	doFunc, ok := funcMap[a.resourceType]
	if !ok {
		return UnsupportedTypeError{resourceType: a.resourceType}
	}

	client, err := server.Connect(a.config.Server, a.configDir)
	if err != nil {
		return fmt.Errorf("error creating the client for the daemon process: %w", err)
	}
	defer client.Close()

	return doFunc(client)
}

func (a *AcceptExecutor) acceptFollowRequest(client *rpc.Client) error {
	accountID, err := getAccountID(client, a.accountName)
	if err != nil {
		return fmt.Errorf("received an error while getting the account ID: %w", err)
	}

	if err := client.Call("GTSClient.AcceptFollowRequest", accountID, nil); err != nil {
		return fmt.Errorf("unable to accept the follow request: %w", err)
	}

	printer.PrintSuccess(a.printSettings, "Successfully accepted the follow request.")

	return nil
}
