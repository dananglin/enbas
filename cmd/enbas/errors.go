// SPDX-FileCopyrightText: 2024 Dan Anglin <d.n.i.anglin@gmail.com>
//
// SPDX-License-Identifier: GPL-3.0-or-later

package main

type unknownCommandError struct {
	subcommand string
}

func (e unknownCommandError) Error() string {
	return "unknown command '" + e.subcommand + "'"
}
