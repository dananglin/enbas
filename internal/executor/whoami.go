package executor

import (
	"fmt"

	"codeflow.dananglin.me.uk/apollo/enbas/internal/config"
)

func (e *WhoamiExecutor) Execute() error {
	config, err := config.NewCredentialsConfigFromFile(e.config.CredentialsFile)
	if err != nil {
		return fmt.Errorf("unable to load the credential config: %w", err)
	}

	e.printer.PrintInfo("You are logged in as '" + config.CurrentAccount + "'.\n")

	return nil
}
