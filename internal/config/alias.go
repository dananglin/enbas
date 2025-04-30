package config

import "fmt"

func CreateAlias(configFilepath, name, arguments string) error {
	cfg, err := newConfigFromFile(configFilepath)
	if err != nil {
		return fmt.Errorf("error loading the configuration from file: %w", err)
	}

	if cfg.Aliases == nil {
		cfg.Aliases = make(map[string]string)
		cfg.Aliases[name] = arguments
	} else {
		if _, exists := cfg.Aliases[name]; exists {
			return NewAliasAlreadyPresentError(name)
		}

		cfg.Aliases[name] = arguments
	}

	if err := saveConfig(configFilepath, cfg); err != nil {
		return fmt.Errorf("error saving the configuration: %w", err)
	}

	return nil
}

func EditAlias(configFilepath, name, arguments string) error {
	cfg, err := newConfigFromFile(configFilepath)
	if err != nil {
		return fmt.Errorf("error loading the configuration from file: %w", err)
	}

	if _, exists := cfg.Aliases[name]; !exists {
		return NewAliasNotPresentError(name)
	}

	cfg.Aliases[name] = arguments

	if err := saveConfig(configFilepath, cfg); err != nil {
		return fmt.Errorf("error saving the configuration: %w", err)
	}

	return nil
}

func DeleteAlias(configFilepath, name string) error {
	cfg, err := newConfigFromFile(configFilepath)
	if err != nil {
		return fmt.Errorf("error loading the configuration from file: %w", err)
	}

	if _, exists := cfg.Aliases[name]; !exists {
		return NewAliasNotPresentError(name)
	}

	delete(cfg.Aliases, name)

	if err := saveConfig(configFilepath, cfg); err != nil {
		return fmt.Errorf("error saving the configuration: %w", err)
	}

	return nil
}

func RenameAlias(configFilepath, oldName, newName string) error {
	cfg, err := newConfigFromFile(configFilepath)
	if err != nil {
		return fmt.Errorf("error loading the configuration from file: %w", err)
	}

	args, exists := cfg.Aliases[oldName]
	if !exists {
		return NewAliasNotPresentError(oldName)
	}

	if _, exists := cfg.Aliases[newName]; exists {
		return NewAliasAlreadyPresentError(newName)
	}

	cfg.Aliases[newName] = args

	delete(cfg.Aliases, oldName)

	if err := saveConfig(configFilepath, cfg); err != nil {
		return fmt.Errorf("error saving the configuration: %w", err)
	}

	return nil
}
