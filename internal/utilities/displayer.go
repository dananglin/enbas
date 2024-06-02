// SPDX-FileCopyrightText: 2024 Dan Anglin <d.n.i.anglin@gmail.com>
//
// SPDX-License-Identifier: GPL-3.0-or-later

package utilities

import "os"

type Displayer interface {
	Display(noColor bool) string
}

func Display(d Displayer, noColor bool) {
	os.Stdout.WriteString(d.Display(noColor) + "\n")
}
