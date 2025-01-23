package executor

import (
	"fmt"

	"codeflow.dananglin.me.uk/apollo/enbas/internal/gtsclient"
	"codeflow.dananglin.me.uk/apollo/enbas/internal/server"
	"codeflow.dananglin.me.uk/apollo/enbas/internal/utilities"
)

func (w *WhoamiExecutor) Execute() error {
	client, err := server.Connect(w.config.Server, w.configDir)
	if err != nil {
		return fmt.Errorf("error creating the client for the daemon process: %w", err)
	}
	defer client.Close()

	account, err := getMyAccount(client)
	if err != nil {
		return fmt.Errorf("error getting your account information: %w", err)
	}

	var instanceURL string
	if err := client.Call("GTSClient.GetInstanceURL", gtsclient.NoRPCArgs{}, &instanceURL); err != nil {
		return fmt.Errorf("error getting the instance URL: %w", err)
	}

	w.printer.PrintInfo("You are logged in as '" + account.Username + "@" + utilities.GetFQDN(instanceURL) + "'.\n")

	return nil
}
