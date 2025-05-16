package definitions

import (
	"encoding/json"
	"fmt"

	"codeflow.dananglin.me.uk/apollo/enbas/internal/utilities"
)

type Definitions struct {
	TopLevelFlags  map[string]TopLevelFlag `json:"topLevelFlags"`
	BuiltInAliases map[string][]string     `json:"builtInAliases"`
	Flags          map[string]string       `json:"flags"`
	Actions        map[string]string       `json:"actions"`
	Targets        map[string]Target       `json:"targets"`
}

type TopLevelFlag struct {
	Description string `json:"description"`
	Type        string `json:"type"`
	Default     string `json:"default"`
	Required    bool   `json:"bool"`
}

type Target struct {
	Description string                  `json:"description"`
	Actions     map[string]TargetAction `json:"actions"`
}

type TargetAction struct {
	Description    string                   `json:"description"`
	ExtraDetails   []string                 `json:"extraDetails"`
	Flags          []TargetActionFlag       `json:"flags"`
	Preposition    string                   `json:"preposition"`
	RelatedTargets map[string]RelatedTarget `json:"relatedTargets"`
}

type RelatedTarget struct {
	Description  string             `json:"description"`
	ExtraDetails []string           `json:"extraDetails"`
	Flags        []TargetActionFlag `json:"flags"`
}

type TargetActionFlag struct {
	Name     string   `json:"name"`
	Type     string   `json:"type"`
	Default  string   `json:"default"`
	Enum     []string `json:"enum"`
	Required bool     `json:"required"`
}

func LoadFromFile(path string) (Definitions, error) {
	file, err := utilities.OpenFile(path)
	if err != nil {
		return Definitions{}, fmt.Errorf("unable to open the definitions file: %w", err)
	}
	defer file.Close()

	var defs Definitions

	if err := json.NewDecoder(file).Decode(&defs); err != nil {
		return Definitions{}, fmt.Errorf("unable to decode the JSON data: %w", err)
	}

	return defs, nil
}
