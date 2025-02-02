package executor

import (
	"fmt"
	"net/rpc"

	"codeflow.dananglin.me.uk/apollo/enbas/internal/model"
	"codeflow.dananglin.me.uk/apollo/enbas/internal/printer"
	"codeflow.dananglin.me.uk/apollo/enbas/internal/server"
)

func (m *UnmuteExecutor) Execute() error {
	funcMap := map[string]func(*rpc.Client) error{
		resourceAccount: m.unmuteAccount,
		resourceStatus:  m.unmuteStatus,
	}

	doFunc, ok := funcMap[m.resourceType]
	if !ok {
		return UnsupportedTypeError{resourceType: m.resourceType}
	}

	client, err := server.Connect(m.config.Server, m.configDir)
	if err != nil {
		return fmt.Errorf("error creating the client for the daemon process: %w", err)
	}
	defer client.Close()

	return doFunc(client)
}

func (m *UnmuteExecutor) unmuteAccount(client *rpc.Client) error {
	accountID, err := getAccountID(client, false, m.accountName)
	if err != nil {
		return fmt.Errorf("received an error while getting the account ID: %w", err)
	}

	if err := client.Call("GTSClient.UnmuteAccount", accountID, nil); err != nil {
		return fmt.Errorf("error unmuting the account: %w", err)
	}

	printer.PrintSuccess(m.printSettings, "Successfully unmuted the account.")

	return nil
}

func (m *UnmuteExecutor) unmuteStatus(client *rpc.Client) error {
	if m.statusID == "" {
		return MissingIDError{
			resource: resourceStatus,
			action:   "unmute",
		}
	}

	var status model.Status
	if err := client.Call("GTSClient.GetStatus", m.statusID, &status); err != nil {
		return fmt.Errorf("unable to retrieve the status: %w", err)
	}

	myAccountID, err := getAccountID(client, true, nil)
	if err != nil {
		return fmt.Errorf("unable to get your account ID: %w", err)
	}

	canUnmute := false

	if status.Account.ID == myAccountID {
		canUnmute = true
	} else {
		for _, mention := range status.Mentions {
			if mention.ID == myAccountID {
				canUnmute = true

				break
			}
		}
	}

	if !canUnmute {
		return Error{"unable to unmute the status because the status does not belong to you nor are you mentioned in it"}
	}

	if err := client.Call("GTSClient.UnmuteStatus", m.statusID, nil); err != nil {
		return fmt.Errorf("error unmuting the status: %w", err)
	}

	printer.PrintSuccess(m.printSettings, "Successfully unmuted the status.")

	return nil
}
