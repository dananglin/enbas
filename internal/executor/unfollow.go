package executor

import (
	"fmt"
	"net/rpc"

	"codeflow.dananglin.me.uk/apollo/enbas/internal/server"
)

func (f *UnfollowExecutor) Execute() error {
	funcMap := map[string]func(*rpc.Client) error{
		resourceAccount: f.unfollowAccount,
	}

	doFunc, ok := funcMap[f.resourceType]
	if !ok {
		return UnsupportedTypeError{resourceType: f.resourceType}
	}

	client, err := server.Connect(f.config.Server, f.configDir)
	if err != nil {
		return fmt.Errorf("error creating the client for the daemon process: %w", err)
	}
	defer client.Close()

	return doFunc(client)
}

func (f *UnfollowExecutor) unfollowAccount(client *rpc.Client) error {
	accountID, err := getAccountID(client, false, f.accountName)
	if err != nil {
		return fmt.Errorf("received an error while getting the account ID: %w", err)
	}

	if err := client.Call("GTSClient.UnfollowAccount", accountID, nil); err != nil {
		return fmt.Errorf("unable to unfollow the account: %w", err)
	}

	f.printer.PrintSuccess("Successfully unfollowed the account.")

	return nil
}
