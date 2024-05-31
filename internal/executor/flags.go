package executor

import "strings"

type AccountNames []string

func (a *AccountNames) String() string {
	return strings.Join(*a, ", ")
}

func (a *AccountNames) Set(value string) error {
	if len(value) > 0 {
		*a = append(*a, value)
	}

	return nil
}

type TopLevelFlags struct {
	ConfigDir string
	NoColor   *bool
}
