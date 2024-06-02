// SPDX-FileCopyrightText: 2024 Dan Anglin <d.n.i.anglin@gmail.com>
//
// SPDX-License-Identifier: GPL-3.0-or-later

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
