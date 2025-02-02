package executor

import (
	"fmt"
	"net/rpc"

	"codeflow.dananglin.me.uk/apollo/enbas/internal/gtsclient"
	"codeflow.dananglin.me.uk/apollo/enbas/internal/model"
	"codeflow.dananglin.me.uk/apollo/enbas/internal/printer"
	"codeflow.dananglin.me.uk/apollo/enbas/internal/server"
)

func (m *MuteExecutor) Execute() error {
	funcMap := map[string]func(*rpc.Client) error{
		resourceAccount: m.muteAccount,
		resourceStatus:  m.muteStatus,
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

func (m *MuteExecutor) muteAccount(client *rpc.Client) error {
	accountID, err := getAccountID(client, false, m.accountName)
	if err != nil {
		return fmt.Errorf("received an error while getting the account ID: %w", err)
	}

	if err := client.Call(
		"GTSClient.MuteAccount",
		gtsclient.MuteAccountArgs{
			AccountID:     accountID,
			Notifications: m.muteNotifications,
			Duration:      int(m.muteDuration.Duration.Seconds()),
		},
		nil,
	); err != nil {
		return fmt.Errorf("error muting the account: %w", err)
	}

	printer.PrintSuccess(m.printSettings, "Successfully muted the account.")

	return nil
}

func (m *MuteExecutor) muteStatus(client *rpc.Client) error {
	if m.statusID == "" {
		return MissingIDError{
			resource: resourceStatus,
			action:   "mute",
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

	canMute := false

	if status.Account.ID == myAccountID {
		canMute = true
	} else {
		for _, mention := range status.Mentions {
			if mention.ID == myAccountID {
				canMute = true

				break
			}
		}
	}

	if !canMute {
		return Error{"unable to mute the status because the status does not belong to you nor are you mentioned in it"}
	}

	if err := client.Call("GTSClient.MuteStatus", m.statusID, nil); err != nil {
		return fmt.Errorf("error muting the status: %w", err)
	}

	printer.PrintSuccess(m.printSettings, "Successfully muted the status.")

	return nil
}
