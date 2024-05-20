package model

import (
	"fmt"

	"codeflow.dananglin.me.uk/apollo/enbas/internal/utilities"
)

type Followers []Account

func (f Followers) String() string {
	output := "\n"
	output += utilities.HeaderFormat("FOLLOWED BY:")

	for i := range f {
		output += fmt.Sprintf(
			"\n  • %s (%s)",
			utilities.DisplayNameFormat(f[i].DisplayName),
			f[i].Acct,
		)
	}

	return output
}

type Following []Account

func (f Following) String() string {
	output := "\n"
	output += utilities.HeaderFormat("FOLLOWING:")

	for i := range f {
		output += fmt.Sprintf(
			"\n  • %s (%s)",
			utilities.DisplayNameFormat(f[i].DisplayName),
			f[i].Acct,
		)
	}

	return output
}
