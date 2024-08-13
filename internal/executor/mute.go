package executor

import (
	"fmt"

	"codeflow.dananglin.me.uk/apollo/enbas/internal/client"
)

func (m *MuteExecutor) Execute() error {
	funcMap := map[string]func(*client.Client) error{
		resourceAccount: m.muteAccount,
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
