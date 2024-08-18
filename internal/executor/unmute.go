package executor

import (
	"fmt"

	"codeflow.dananglin.me.uk/apollo/enbas/internal/client"
)

func (m *UnmuteExecutor) Execute() error {
	funcMap := map[string]func(*client.Client) error{
		resourceAccount: m.unmuteAccount,
		resourceStatus:  m.unmuteStatus,
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

func (m *UnmuteExecutor) unmuteAccount(gtsClient *client.Client) error {
	accountID, err := getAccountID(gtsClient, false, m.accountName)
	if err != nil {
		return fmt.Errorf("received an error while getting the account ID: %w", err)
	}

	if err := gtsClient.UnmuteAccount(accountID); err != nil {
		return fmt.Errorf("unable to unmute the account: %w", err)
	}

	m.printer.PrintSuccess("Successfully unmuted the account.")

	return nil
}

func (m *UnmuteExecutor) unmuteStatus(gtsClient *client.Client) error {
	if m.statusID == "" {
		return MissingIDError{
			resource: resourceStatus,
			action:   "unmute",
		}
	}

	status, err := gtsClient.GetStatus(m.statusID)
	if err != nil {
		return fmt.Errorf("unable to retrieve the status: %w", err)
	}

	myAccountID, err := getAccountID(gtsClient, true, nil)
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

	if err := gtsClient.UnmuteStatus(m.statusID); err != nil {
		return fmt.Errorf("unable to unmute the status: %w", err)
	}

	m.printer.PrintSuccess("Successfully unmuted the status.")

	return nil
}
