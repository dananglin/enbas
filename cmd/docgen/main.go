package main

import (
	"flag"
	"fmt"
	"os"

	"codeflow.dananglin.me.uk/apollo/enbas/internal/gen/definitions"
)

func main() {
	const definitionsPath = "./definitions/definitions.json"

	var (
		applicationName string
		binaryVersion   string
		rootManDir      string
		examplesDir     string
	)

	flag.StringVar(&applicationName, "application-name", "enbas", "The name of the application")
	flag.StringVar(&binaryVersion, "binary-version", "", "The application's version")
	flag.StringVar(&rootManDir, "man-dir", ".", "The root directory for the manpages")
	flag.StringVar(&examplesDir, "examples-dir", ".", "The directory where the example configuration is generated")
	flag.Parse()

	defs, err := definitions.LoadFromFile(definitionsPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: Unable to read the definitions file: %v.\n", err)

		os.Exit(1)
	}

	if err := generateManual(
		defs,
		applicationName,
		binaryVersion,
		rootManDir,
	); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: Unable to generate the manual: %v.\n", err)

		os.Exit(1)
	}

	if err := generateExampleConfig(
		applicationName,
		examplesDir,
	); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: Unable to generate the example configuration: %v.\n", err)

		os.Exit(1)
	}
}
