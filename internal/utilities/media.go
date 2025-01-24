package utilities

import (
	"fmt"
	"os/exec"
	"slices"
	"strings"
)

func OpenMedia(viewer string, paths []string) error {
	if viewer == "" {
		return UnspecifiedProgramError{}
	}

	cmd := slices.Concat(strings.Split(viewer, " "), paths)

	command := exec.Command(cmd[0], cmd[1:]...) // #nosec G204 -- External command call defined in user's configuration file.

	if err := command.Start(); err != nil {
		return fmt.Errorf("received an error after starting the program: %w", err)
	}

	return nil
}
