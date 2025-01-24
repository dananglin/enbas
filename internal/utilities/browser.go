package utilities

import (
	"fmt"
	"os/exec"
	"strings"
)

func OpenLink(browser, url string) error {
	if browser == "" {
		return UnspecifiedBrowserError{}
	}

	cmd := strings.Split(browser, " ")
	cmd = append(cmd, url)

	command := exec.Command(cmd[0], cmd[1:]...) // #nosec G204 -- External command call defined in user's configuration file.

	if err := command.Start(); err != nil {
		return fmt.Errorf("received an error after starting the program to view the link: %w", err)
	}

	return nil
}
