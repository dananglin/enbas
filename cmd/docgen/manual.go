package main

import (
	"embed"
	"fmt"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"codeflow.dananglin.me.uk/apollo/enbas/internal/gen/definitions"
	genUtils "codeflow.dananglin.me.uk/apollo/enbas/internal/gen/utilities"
	"codeflow.dananglin.me.uk/apollo/enbas/internal/utilities"
)

type manpage struct {
	templateName string
	dir          string
	page         string
	navRef       string
	description  string
}

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

	manpages := []manpage{
		{
			templateName: "enbas.1",
			dir:          "man1",
			page:         applicationName + ".1",
			navRef:       applicationName + "(1)",
			description:  "General operations manual.",
		},
		{
			templateName: "enbas.5",
			dir:          "man5",
			page:         applicationName + ".5",
			navRef:       applicationName + "(5)",
			description:  "The configuration manual.",
		},
		{
			templateName: "enbas-topics.7",
			dir:          "man7",
			page:         applicationName + "-topics.7",
			navRef:       applicationName + "-topics(7)",
			description:  "Manual containing details of the features in " + applicationName + ".",
		},
	}

	funcMap := template.FuncMap{
		"titled":                     genUtils.Titled,
		"allCaps":                    genUtils.AllCaps,
		"dateNow":                    dateNow,
		"newOperation":               newOperation,
		"newTargetToTargetOperation": newTargetToTargetOperation,
		"seeAlso":                    seeAlsoFunc(manpages),
		"builtInAlias":               builtInAlias,
	}

	tmpl := template.Must(template.New("").
		Funcs(funcMap).
		ParseFS(manualTemplates, "templates/manpages/*"),
	)

	for idx := range manpages {
		if err := func() error {
			outputDir := filepath.Join(
				rootManDir,
				manpages[idx].dir,
			)

			outputPath := filepath.Join(
				outputDir,
				manpages[idx].page,
			)

			if err := utilities.EnsureDirectory(outputDir); err != nil {
				return fmt.Errorf("error ensuring the presence of directory %q: %w", outputDir, err)
			}

			file, err := utilities.CreateFile(outputPath)
			if err != nil {
				return fmt.Errorf("error creating the output file: %w", err)
			}
			defer file.Close()

			if err := tmpl.ExecuteTemplate(
				file,
				manpages[idx].templateName,
				data,
			); err != nil {
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

func seeAlsoFunc(manpages []manpage) func(string) string {
	return func(templateName string) string {
		var builder strings.Builder

		builder.WriteString(".SH SEE ALSO\n")

		for idx := range manpages {
			if manpages[idx].templateName == templateName {
				continue
			}

			builder.WriteString(".B " + manpages[idx].navRef)
			builder.WriteString("\n.RS\n")
			builder.WriteString(manpages[idx].description)
			builder.WriteString("\n.br\n")
			builder.WriteString("\n.RE\n")
		}

		return builder.String()
	}
}

func builtInAlias(alias string, operation definitions.BuiltInAlias) string {
	var builder strings.Builder

	builder.WriteString(".B " + alias)
	builder.WriteString("\n.RS\n")
	builder.WriteString(genUtils.Titled(operation.Description) + ".")
	builder.WriteString("\n.br\n")
	builder.WriteString("This is an alias for \\fB\"" + strings.Join(operation.Operation, " ") + "\"\\fR\\&.")
	builder.WriteString("\n.RE\n")

	return builder.String()
}
