// SPDX-FileCopyrightText: 2024 Dan Anglin <d.n.i.anglin@gmail.com>
//
// SPDX-License-Identifier: GPL-3.0-or-later

package utilities

import (
	"os"
	"os/exec"
	"runtime"
)

func OpenLink(url string) {
	var open string

	envBrower := os.Getenv("BROWSER")

	switch {
	case len(envBrower) > 0:
		open = envBrower
	case runtime.GOOS == "linux":
		open = "xdg-open"
	default:
		return
	}

	command := exec.Command(open, url)

	_ = command.Start()
}
