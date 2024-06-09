// SPDX-FileCopyrightText: 2024 Dan Anglin <d.n.i.anglin@gmail.com>
//
// SPDX-License-Identifier: GPL-3.0-or-later

package utilities

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type Displayer interface {
	Display(noColor bool) string
}

func Display(displayer Displayer, noColor bool, pagerCommand string) {
	if pagerCommand == "" {
		os.Stdout.WriteString(displayer.Display(noColor) + "\n")

		return
	}

	split := strings.Split(pagerCommand, " ")

	pager := new(exec.Cmd)

	if len(split) == 1 {
		pager = exec.Command(split[0]) //nolint:gosec
	} else {
		pager = exec.Command(split[0], split[1:]...) //nolint:gosec
	}

	pipe, err := pager.StdinPipe()
	if err != nil {
		os.Stdout.WriteString(displayer.Display(noColor) + "\n")

		return
	}

	pager.Stdout = os.Stdout
	pager.Stderr = os.Stderr

	_ = pager.Start()

	defer func() {
		_ = pipe.Close()
		_ = pager.Wait()
	}()

	fmt.Fprintln(pipe, displayer.Display(noColor)+"\n")
}
