package executor

import (
	"fmt"
	"net/rpc"

	"codeflow.dananglin.me.uk/apollo/enbas/internal/printer"
	"codeflow.dananglin.me.uk/apollo/enbas/internal/server"
)

func (b *UnblockExecutor) Execute() error {
	funcMap := map[string]func(*rpc.Client) error{
		resourceAccount: b.unblockAccount,
	}

	doFunc, ok := funcMap[b.resourceType]
	if !ok {
		return UnsupportedTypeError{resourceType: b.resourceType}
	}

	client, err := server.Connect(b.config.Server, b.configDir)
	if err != nil {
		return fmt.Errorf("error creating the client for the daemon process: %w", err)
	}
	defer client.Close()

	return doFunc(client)
}

func (b *UnblockExecutor) unblockAccount(client *rpc.Client) error {
	accountID, err := getAccountID(client, b.accountName)
	if err != nil {
		return fmt.Errorf("received an error while getting the account ID: %w", err)
	}

	if err := client.Call("GTSClient.UnblockAccount", accountID, nil); err != nil {
		return fmt.Errorf("unable to unblock the account: %w", err)
	}

	printer.PrintSuccess(b.printSettings, "Successfully unblocked the account.")

	return nil
}
