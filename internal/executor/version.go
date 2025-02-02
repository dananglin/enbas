package executor

import (
	"fmt"

	"codeflow.dananglin.me.uk/apollo/enbas/internal/printer"
)

func (v *VersionExecutor) Execute() error {
	if err := printer.PrintVersion(v.printSettings, v.full); err != nil {
		return fmt.Errorf("error printing the version: %w", err)
	}

	return nil
}
