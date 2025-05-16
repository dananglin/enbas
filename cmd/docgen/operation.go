package main

import (
	"strings"

	"codeflow.dananglin.me.uk/apollo/enbas/internal/gen/definitions"
	genUtils "codeflow.dananglin.me.uk/apollo/enbas/internal/gen/utilities"
)

type operation struct {
	ApplicationName string
	OperationName   string
	FullDescription string
	Flags           []operationFlag
}

type operationFlag struct {
	Name            string
	Required        bool
	Default         string
	FullDescription string
}

func newOperation(
	application string,
	action string,
	target string,
	description string,
	extraDetails []string,
	flags []definitions.TargetActionFlag,
	flagDefinitions map[string]string,
) operation {
	fullDesc := genUtils.Titled(description) + "."
	if len(extraDetails) > 0 {
		fullDesc += " " + strings.Join(extraDetails, " ")
	}

	operationFlags := make([]operationFlag, len(flags))

	for idx := range flags {
		desc := genUtils.Titled(
			strings.ReplaceAll(
				strings.ReplaceAll(
					flagDefinitions[flags[idx].Name]+".",
					"{target}",
					target,
				),
				"{action}",
				action,
			),
		)

		if len(flags[idx].Enum) > 0 {
			desc += listValidValues(flags[idx].Enum)
		}

		if strings.HasPrefix(flags[idx].Type, "internalFlag.Multi") {
			desc += " This flag can be used multiple times to specify multiple values."
		}

		operationFlags[idx] = operationFlag{
			Name:            flags[idx].Name,
			Required:        flags[idx].Required,
			Default:         flags[idx].Default,
			FullDescription: desc,
		}
	}

	return operation{
		ApplicationName: application,
		OperationName:   action + " " + target,
		FullDescription: fullDesc,
		Flags:           operationFlags,
	}
}

func newTargetToTargetOperation(
	application string,
	action string,
	focusedTarget string,
	preposition string,
	relatedTarget string,
	description string,
	extraDetails []string,
	flags []definitions.TargetActionFlag,
	flagDefinitions map[string]string,
) operation {
	fullDesc := genUtils.Titled(description) + "."
	if len(extraDetails) > 0 {
		fullDesc += " " + strings.Join(extraDetails, " ")
	}

	operationFlags := make([]operationFlag, len(flags))

	for idx := range flags {
		desc := genUtils.Titled(
			strings.ReplaceAll(
				strings.ReplaceAll(
					flagDefinitions[flags[idx].Name]+".",
					"{target}",
					focusedTarget,
				),
				"{action}",
				action,
			),
		)

		if len(flags[idx].Enum) > 0 {
			desc += listValidValues(flags[idx].Enum)
		}

		if strings.HasPrefix(flags[idx].Type, "internalFlag.Multi") {
			desc += " This flag can be used multiple times to specify multiple values."
		}

		operationFlags[idx] = operationFlag{
			Name:            flags[idx].Name,
			Required:        flags[idx].Required,
			Default:         flags[idx].Default,
			FullDescription: desc,
		}
	}

	return operation{
		ApplicationName: application,
		OperationName:   action + " " + focusedTarget + " " + preposition + " " + relatedTarget,
		FullDescription: fullDesc,
		Flags:           operationFlags,
	}
}

func listValidValues(values []string) string {
	output := " The valid values for this flag are "
	numValues := len(values)

	for idx := range numValues {
		if idx == numValues-1 {
			output += "\"" + values[idx] + "\"."

			continue
		}

		if idx == numValues-2 {
			output += "\"" + values[idx] + "\", and "

			continue
		}

		output += "\"" + values[idx] + "\", "
	}

	return output
}
