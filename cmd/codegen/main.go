package main

import (
	"embed"
	"errors"
	"flag"
	"fmt"
	"maps"
	"os"
	"os/exec"
	"slices"
	"strings"
	"text/template"
	"unicode"

	"codeflow.dananglin.me.uk/apollo/enbas/internal/utilities"
)

func main() {
	var (
		pathToDefinitions  string
		packageName        string
		pathToToolsModfile string
	)

	flag.StringVar(&pathToDefinitions, "path-to-definitions", "", "The path to the definitions file")
	flag.StringVar(&packageName, "package", "", "The name of the internal package")
	flag.StringVar(&pathToToolsModfile, "path-to-tools-modfile", "", "The path to the tools modfile")
	flag.Parse()

	defs, err := loadDefinitionsFromFile(pathToDefinitions)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: Unable to read the definitions file: %v.\n", err)
		os.Exit(1)
	}

	if err := generateExecutors(defs, packageName, pathToToolsModfile); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: Unable to generate the executors: %v.\n", err)
		os.Exit(1)
	}
}

//go:embed templates/*
var executorTemplates embed.FS

var errUndefinedPackageName = errors.New("the package name is not defined")

func generateExecutors(defs definitions, packageName, toolsModfile string) error {
	if packageName == "" {
		return errUndefinedPackageName
	}

	dirName := "templates/" + packageName

	fsDir, err := executorTemplates.ReadDir(dirName)
	if err != nil {
		return fmt.Errorf("unable to read the template directory in the file system (FS): %w", err)
	}

	funcMap := template.FuncMap{
		"capitalise":          capitalise,
		"snakeToCamel":        snakeToCamel,
		"getTargetsForAction": getTargetsForAction,
	}

	for _, obj := range fsDir {
		templateFilename := obj.Name()

		if !strings.HasSuffix(templateFilename, ".go.gotmpl") {
			continue
		}

		if err := func() error {
			tmpl := template.Must(template.New(templateFilename).
				Funcs(funcMap).
				ParseFS(executorTemplates, dirName+"/"+templateFilename),
			)

			output := strings.TrimSuffix(templateFilename, ".gotmpl")

			file, err := utilities.CreateFile(output)
			if err != nil {
				return fmt.Errorf("error creating the output file: %w", err)
			}
			defer file.Close()

			if err := tmpl.Execute(file, defs); err != nil {
				return fmt.Errorf("error generating the code from the template: %w", err)
			}

			if err := runGoImports(output, toolsModfile); err != nil {
				return fmt.Errorf("error running goimports: %w", err)
			}

			return nil
		}(); err != nil {
			return fmt.Errorf("received an error after attempting to generate the code for %q: %w", templateFilename, err)
		}
	}

	return nil
}

func runGoImports(path, toolsModfile string) error {
	imports := exec.Command(
		"go",
		"tool",
		"-modfile",
		toolsModfile,
		"goimports",
		"-w",
		path,
	)

	if err := imports.Run(); err != nil {
		return fmt.Errorf("received an error after running goimports: %w", err)
	}

	return nil
}

func capitalise(str string) string {
	runes := []rune(str)

	runes[0] = unicode.ToUpper(runes[0])

	return string(runes)
}

func snakeToCamel(value string, pascal bool) string {
	var builder strings.Builder

	runes := []rune(value)
	numRunes := len(runes)
	cursor := 0

	for cursor < numRunes {
		switch {
		case cursor == 0 && pascal:
			builder.WriteRune(unicode.ToUpper(runes[cursor]))

			cursor++
		case runes[cursor] != '-':
			builder.WriteRune(runes[cursor])

			cursor++
		case cursor != numRunes-1 && unicode.IsLower(runes[cursor+1]):
			builder.WriteRune(unicode.ToUpper(runes[cursor+1]))

			cursor += 2
		default:
			cursor++
		}
	}

	return builder.String()
}

func getTargetsForAction(targetMap map[string]target, action string) []string {
	output := make([]string, 0)

	for name, target := range maps.All(targetMap) {
		if _, ok := target.Actions[action]; ok {
			output = append(output, name)
		}
	}

	slices.Sort(output)

	return output
}
