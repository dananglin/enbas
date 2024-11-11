package main

import (
	"encoding/json"
	"fmt"

	"codeflow.dananglin.me.uk/apollo/enbas/internal/utilities"
)

type enbasCLISchema struct {
	Flags    enbasCLISchemaFlagMap            `json:"flags"`
	Commands map[string]enbasCLISchemaCommand `json:"commands"`
}

func newEnbasCLISchemaFromFile(path string) (enbasCLISchema, error) {
	file, err := utilities.OpenFile(path)
	if err != nil {
		return enbasCLISchema{}, fmt.Errorf("unable to open the schema file: %w", err)
	}
	defer file.Close()

	var schema enbasCLISchema

	if err := json.NewDecoder(file).Decode(&schema); err != nil {
		return enbasCLISchema{}, fmt.Errorf("unable to decode the JSON data: %w", err)
	}

	return schema, nil
}

type enbasCLISchemaFlag struct {
	Type        string `json:"type"`
	Description string `json:"description"`
}

type enbasCLISchemaFlagMap map[string]enbasCLISchemaFlag

func (e enbasCLISchemaFlagMap) getType(name string) string {
	flag, ok := e[name]
	if !ok {
		return "UNKNOWN TYPE"
	}

	return flag.Type
}

func (e enbasCLISchemaFlagMap) getDescription(name string) string {
	flag, ok := e[name]
	if !ok {
		return "UNKNOWN DESCRIPTION"
	}

	return flag.Description
}

type enbasCLISchemaCommand struct {
	AdditionalFields []enbasCLISchemaAdditionalFields `json:"additionalFields"`
	Flags            []enbasCLISchemaFlagReference    `json:"flags"`
	Summary          string                           `json:"summary"`
	UseConfig        bool                             `json:"useConfig"`
	UsePrinter       bool                             `json:"usePrinter"`
}

type enbasCLISchemaFlagReference struct {
	Flag      string `json:"flag"`
	FieldName string `json:"fieldName"`
	Default   string `json:"default"`
}

type enbasCLISchemaAdditionalFields struct {
	Name string `json:"name"`
	Type string `json:"type"`
}
