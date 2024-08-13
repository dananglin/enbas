package executor

import (
	"fmt"

	"codeflow.dananglin.me.uk/apollo/enbas/internal/config"
)

func (s *SwitchExecutor) Execute() error {
	funcMap := map[string]func() error{
		resourceAccount: s.switchToAccount,
	}

	doFunc, ok := funcMap[s.to]
	if !ok {
		return UnsupportedTypeError{resourceType: s.to}
	}

	return doFunc()
}

func (s *SwitchExecutor) switchToAccount() error {
	expectedNumAccountNames := 1
	if !s.accountName.ExpectedLength(expectedNumAccountNames) {
		return fmt.Errorf(
			"found an unexpected number of --account-name flags: expected %d",
			expectedNumAccountNames,
		)
	}

	if err := config.UpdateCurrentAccount(s.accountName[0], s.config.CredentialsFile); err != nil {
		return fmt.Errorf("unable to switch account to the account: %w", err)
	}

	s.printer.PrintSuccess("The current account is now set to '" + s.accountName[0] + "'.")

	return nil
}
