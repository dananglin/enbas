package main

import "strings"

type accountIDs []string

func (a *accountIDs) String() string {
	return strings.Join(*a, ", ")
}

func (a *accountIDs) Set(value string) error {
	if len(value) > 0 {
		*a = append(*a, value)
	}

	return nil
}
