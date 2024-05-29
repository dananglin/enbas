package utilities

import (
	"os/exec"
	"runtime"
)

func OpenLink(url string) {
	var open string
	//envBrower := os.Getenv("BROWSER")

	switch {
	//case len(envBrower) > 0:
	//	open = envBrower
	case runtime.GOOS == "linux":
		open = "xdg-open"
	default:
		return
	}

	command := exec.Command(open, url)

	_ = command.Start()
}
