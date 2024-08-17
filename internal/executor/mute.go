package executor

import (
	"fmt"

	"codeflow.dananglin.me.uk/apollo/enbas/internal/client"
)

func (m *MuteExecutor) Execute() error {
	funcMap := map[string]func(*client.Client) error{
		resourceAccount: m.muteAccount,
		resourceStatus:  m.muteStatus,
	}

	doFunc, ok := funcMap[m.resourceType]
	if !ok {
		return UnsupportedTypeError{resourceType: m.resourceType}
	}

	gtsClient, err := client.NewClientFromFile(m.config.CredentialsFile)
	if err != nil {
		return fmt.Errorf("unable to create the GoToSocial client: %w", err)
	}

	return doFunc(gtsClient)
}

func (m *MuteExecutor) muteAccount(gtsClient *client.Client) error {
	accountID, err := getAccountID(gtsClient, false, m.accountName)
	if err != nil {
		return fmt.Errorf("received an error while getting the account ID: %w", err)
	}

	form := client.MuteAccountForm{
		Notifications: m.muteNotifications,
		Duration:      int(m.muteDuration.Duration.Seconds()),
	}

	if err := gtsClient.MuteAccount(accountID, form); err != nil {
		return fmt.Errorf("unable to mute the account: %w", err)
	}

	m.printer.PrintSuccess("Successfully muted the account.")

	return nil
}

func (m *MuteExecutor) muteStatus(gtsClient *client.Client) error {
	if m.statusID == "" {
		return FlagNotSetError{flagText: flagStatusID}
	}

	status, err := gtsClient.GetStatus(m.statusID)
	if err != nil {
		return fmt.Errorf("unable to retrieve the status: %w", err)
	}

	myAccountID, err := getAccountID(gtsClient, true, nil)
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

	if err := gtsClient.MuteStatus(m.statusID); err != nil {
		return fmt.Errorf("unable to mute the status: %w", err)
	}

	m.printer.PrintSuccess("Successfully muted the status.")

	return nil
}
