package main

import (
	"embed"
	"fmt"
	"path/filepath"
	"slices"
	"text/template"
	"time"

	"codeflow.dananglin.me.uk/apollo/enbas/internal/gen/definitions"
	genUtils "codeflow.dananglin.me.uk/apollo/enbas/internal/gen/utilities"
	"codeflow.dananglin.me.uk/apollo/enbas/internal/utilities"
)

//go:embed templates/manpages/*
var manualTemplates embed.FS

func generateManual(
	defs definitions.Definitions,
	applicationName string,
	binaryVersion string,
	rootManDir string,
) error {
	data := struct {
		ApplicationName string
		BinaryVersion   string
		Definitions     definitions.Definitions
	}{
		ApplicationName: applicationName,
		BinaryVersion:   binaryVersion,
		Definitions:     defs,
	}

	funcMap := template.FuncMap{
		"titled":                     genUtils.Titled,
		"allCaps":                    genUtils.AllCaps,
		"dateNow":                    dateNow,
		"newOperation":               newOperation,
		"newTargetToTargetOperation": newTargetToTargetOperation,
	}

	tmpl := template.Must(template.New("").
		Funcs(funcMap).
		ParseFS(manualTemplates, "templates/manpages/*"),
	)

	for _, section := range slices.All([]string{"1", "5"}) {
		if err := func() error {
			var (
				templateName = "enbas." + section
				outputPath   = filepath.Join(
					rootManDir,
					"man"+section,
					applicationName+"."+section,
				)
			)

			file, err := utilities.CreateFile(outputPath)
			if err != nil {
				return fmt.Errorf("error creating the output file: %w", err)
			}
			defer file.Close()

			if err := tmpl.ExecuteTemplate(file, templateName, data); err != nil {
				return fmt.Errorf(
					"error generating %q: %w",
					outputPath,
					err,
				)
			}

			return nil
		}(); err != nil {
			return fmt.Errorf(
				"received an error after attempting to generate the manual pages: %w",
				err,
			)
		}
	}

	return nil
}

func dateNow() string {
	return time.Now().Format(time.DateOnly)
}
