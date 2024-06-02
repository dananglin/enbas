// SPDX-FileCopyrightText: 2024 Dan Anglin <d.n.i.anglin@gmail.com>
//
// SPDX-License-Identifier: GPL-3.0-or-later

package executor

import "fmt"

type Executor interface {
	Name() string
	Parse(args []string) error
	Execute() error
}

func Execute(executor Executor, args []string) error {
	if err := executor.Parse(args); err != nil {
		return fmt.Errorf("unable to parse the command line flags; %w", err)
	}

	if err := executor.Execute(); err != nil {
		return fmt.Errorf("unable to execute the command %q; %w", executor.Name(), err)
	}

	return nil
}
