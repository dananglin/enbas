package command

import "codeflow.dananglin.me.uk/apollo/enbas/internal/cli"

type Command struct {
	Action             string
	FocusedTarget      string
	FocusedTargetFlags []string
	Preposition        string
	RelatedTarget      string
	RelatedTargetFlags []string
}

func UsageCommand() Command {
	return Command{
		Action:             cli.ActionShow,
		FocusedTarget:      cli.TargetUsage,
		FocusedTargetFlags: []string{},
		Preposition:        "",
		RelatedTarget:      "",
		RelatedTargetFlags: []string{},
	}
}
