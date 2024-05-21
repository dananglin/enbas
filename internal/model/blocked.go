package model

import (
	"fmt"

	"codeflow.dananglin.me.uk/apollo/enbas/internal/utilities"
)

type BlockedAccounts []Account

func (b BlockedAccounts) String() string {
	output := "\n"
	output += utilities.HeaderFormat("BLOCKED ACCOUNTS:")

	for i := range b {
		output += fmt.Sprintf(
			"\n  • %s (%s)",
			b[i].Acct,
			b[i].ID,
		)
	}

	return output
}
