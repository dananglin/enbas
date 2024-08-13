package main

import (
	"embed"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"text/template"
	"unicode"
)

func main() {
	var (
		enbasCLISchemaFilepath string
		packageName            string
	)

	flag.StringVar(&enbasCLISchemaFilepath, "path-to-enbas-cli-schema", "", "The path to the Enbas CLI schema file")
	flag.StringVar(&packageName, "package", "", "The name of the internal package")
	flag.Parse()

	schema, err := newEnbasCLISchemaFromFile(enbasCLISchemaFilepath)
	if err != nil {
		fmt.Printf("ERROR: Unable to read the schema file: %v.\n", err)
	}

	if err := generateExecutors(schema, packageName); err != nil {
		fmt.Printf("ERROR: Unable to generate the executors: %v.\n", err)
	}
}

//go:embed templates/*
var executorTemplates embed.FS

var errNoPackageFlag = errors.New("the --package flag must be used")

func generateExecutors(schema enbasCLISchema, packageName string) error {
	if packageName == "" {
		return errNoPackageFlag
	}

	dirName := "templates/" + packageName

	fsDir, err := executorTemplates.ReadDir(dirName)
	if err != nil {
		return fmt.Errorf("unable to read the template directory in the file system (FS): %w", err)
	}

	funcMap := template.FuncMap{
		"capitalise":         capitalise,
		"flagFieldName":      flagFieldName,
		"getFlagType":        schema.Flags.getType,
		"getFlagDescription": schema.Flags.getDescription,
		"internalFlagValue":  internalFlagValue,
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

			file, err := os.Create(output)
			if err != nil {
				return fmt.Errorf("unable to create the output file: %w", err)
			}
			defer file.Close()

			if err := tmpl.Execute(file, schema.Commands); err != nil {
				return fmt.Errorf("unable to generate the code from the template: %w", err)
			}

			if err := runGoImports(output); err != nil {
				return fmt.Errorf("unable to run goimports: %w", err)
			}

			return nil
		}(); err != nil {
			return fmt.Errorf("received an error after attempting to generate the code for %q: %w", templateFilename, err)
		}
	}

	return nil
}

func runGoImports(path string) error {
	imports := exec.Command("goimports", "-w", path)

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

func flagFieldName(flagRef enbasCLISchemaFlagReference) string {
	if flagRef.FieldName != "" {
		return flagRef.FieldName
	}

	return convertFlagToMixedCaps(flagRef.Flag)
}

func convertFlagToMixedCaps(value string) string {
	var builder strings.Builder

	runes := []rune(value)
	numRunes := len(runes)
	cursor := 0

	for cursor < numRunes {
		if runes[cursor] != '-' {
			builder.WriteRune(runes[cursor])

			cursor++
		} else {
			if cursor != numRunes-1 && unicode.IsLower(runes[cursor+1]) {
				builder.WriteRune(unicode.ToUpper(runes[cursor+1]))
				cursor += 2
			} else {
				cursor++
			}
		}
	}

	return builder.String()
}

func internalFlagValue(flagType string) bool {
	internalFlagValues := map[string]struct{}{
		"StringSliceValue":  {},
		"IntSliceValue":     {},
		"TimeDurationValue": {},
		"BoolPtrValue":      {},
	}

	_, exists := internalFlagValues[flagType]

	return exists
}
