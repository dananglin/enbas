// SPDX-FileCopyrightText: 2024 Dan Anglin <d.n.i.anglin@gmail.com>
//
// SPDX-License-Identifier: GPL-3.0-or-later

package utilities

import (
	"fmt"
	"os/exec"
	"regexp"
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

	command := exec.Command(viewer, paths...)

	if err := command.Start(); err != nil {
		return fmt.Errorf("received an error after starting the image viewer: %w", err)
	}

	return nil
}
