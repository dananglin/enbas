package main

import (
	"os"

	"codeflow.dananglin.me.uk/apollo/enbas/internal/executor"
)

func main() {
	if err := executor.Execute(); err != nil {
		os.Exit(1)
	}
}
