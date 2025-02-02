package executor

import (
	"fmt"
	"net/rpc"

	"codeflow.dananglin.me.uk/apollo/enbas/internal/printer"
	"codeflow.dananglin.me.uk/apollo/enbas/internal/server"
)

func (b *BlockExecutor) Execute() error {
	funcMap := map[string]func(*rpc.Client) error{
		resourceAccount: b.blockAccount,
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

func (b *BlockExecutor) blockAccount(client *rpc.Client) error {
	accountID, err := getAccountID(client, false, b.accountName)
	if err != nil {
		return fmt.Errorf("received an error while getting the account ID: %w", err)
	}

	if err := client.Call("GTSClient.BlockAccount", accountID, nil); err != nil {
		return fmt.Errorf("unable to block the account: %w", err)
	}

	printer.PrintSuccess(b.printSettings, "Successfully blocked the account.")

	return nil
}
