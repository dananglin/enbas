package main

import (
	"encoding/json"
	"fmt"

	"codeflow.dananglin.me.uk/apollo/enbas/internal/utilities"
)

type definitions struct {
	TopLevelFlags  map[string]topLevelFlag `json:"topLevelFlags"`
	BuiltInAliases map[string][]string     `json:"builtInAliases"`
	Flags          map[string]string       `json:"flags"`
	Actions        map[string]string       `json:"actions"`
	Targets        map[string]target       `json:"targets"`
}

type topLevelFlag struct {
	Description string `json:"description"`
	Type        string `json:"type"`
	Default     string `json:"default"`
	Required    bool   `json:"bool"`
}

type target struct {
	Description string                  `json:"description"`
	Actions     map[string]targetAction `json:"actions"`
}

type targetAction struct {
	Flags          []targetActionFlag       `json:"flags"`
	Preposition    string                   `json:"preposition"`
	RelatedTargets map[string]relatedTarget `json:"relatedTargets"`
}

type relatedTarget struct {
	Flags []targetActionFlag `json:"flags"`
}

type targetActionFlag struct {
	Name     string   `json:"name"`
	Type     string   `json:"type"`
	Default  string   `json:"default"`
	Enum     []string `json:"enum"`
	Required bool     `json:"required"`
}

func loadDefinitionsFromFile(path string) (definitions, error) {
	file, err := utilities.OpenFile(path)
	if err != nil {
		return definitions{}, fmt.Errorf("unable to open the definitions file: %w", err)
	}
	defer file.Close()

	var defs definitions

	if err := json.NewDecoder(file).Decode(&defs); err != nil {
		return definitions{}, fmt.Errorf("unable to decode the JSON data: %w", err)
	}

	return defs, nil
}
