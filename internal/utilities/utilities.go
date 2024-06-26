// SPDX-FileCopyrightText: 2024 Dan Anglin <d.n.i.anglin@gmail.com>
//
// SPDX-License-Identifier: GPL-3.0-or-later

package utilities

import (
	"fmt"
	"os/exec"
	"regexp"
	"slices"
	"strings"
)

func GetFQDN(url string) string {
	r := regexp.MustCompile(`http(s)?:\/\/`)

	return r.ReplaceAllString(url, "")
}

type UnspecifiedProgramError struct{}

func (e UnspecifiedProgramError) Error() string {
	return "the program to view these files is unspecified"
}

func OpenMedia(viewer string, paths []string) error {
	if viewer == "" {
		return UnspecifiedProgramError{}
	}

	cmd := slices.Concat(strings.Split(viewer, " "), paths)

	command := exec.Command(cmd[0], cmd[1:]...) //nolint:gosec

	if err := command.Start(); err != nil {
		return fmt.Errorf("received an error after starting the program: %w", err)
	}

	return nil
}

type UnspecifiedBrowserError struct{}

func (e UnspecifiedBrowserError) Error() string {
	return "the browser to view this link is not specified"
}

func OpenLink(browser, url string) error {
	if browser == "" {
		return UnspecifiedBrowserError{}
	}

	cmd := strings.Split(browser, " ")
	cmd = append(cmd, url)

	command := exec.Command(cmd[0], cmd[1:]...)

	if err := command.Start(); err != nil {
		return fmt.Errorf("received an error after starting the program to view the link: %w", err)
	}

	return nil
}
