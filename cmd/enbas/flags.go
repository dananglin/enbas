package main

import "strings"

type accountNames []string

func (a *accountNames) String() string {
	return strings.Join(*a, ", ")
}

func (a *accountNames) Set(value string) error {
	if len(value) > 0 {
		*a = append(*a, value)
	}

	return nil
}

type topLevelFlags struct {
	configDir string
}
