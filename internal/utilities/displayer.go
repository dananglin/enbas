package utilities

import "os"

type Displayer interface {
	Display(noColor bool) string
}

func Display(d Displayer, noColor bool) {
	os.Stdout.WriteString(d.Display(noColor) + "\n")
}
