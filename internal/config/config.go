// SPDX-FileCopyrightText: 2024 Dan Anglin <d.n.i.anglin@gmail.com>
//
// SPDX-License-Identifier: GPL-3.0-or-later

package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"codeflow.dananglin.me.uk/apollo/enbas/internal/utilities"
)

const (
	configFileName = "config.json"
)

type Config struct {
	CredentialsFile  string       `json:"credentialsFile"`
	CacheDirectory   string       `json:"cacheDirectory"`
	LineWrapMaxWidth int          `json:"lineWrapMaxWidth"`
	HTTP             HTTPConfig   `json:"http"`
	Integrations     Integrations `json:"integrations"`
}

type HTTPConfig struct {
	Timeout      int `json:"timeout"`
	MediaTimeout int `json:"mediaTimeout"`
}

type Integrations struct {
	Browser     string `json:"browser"`
	Editor      string `json:"editor"`
	Pager       string `json:"pager"`
	ImageViewer string `json:"imageViewer"`
	VideoPlayer string `json:"videoPlayer"`
}

func NewConfigFromFile(configDir string) (*Config, error) {
	path := configFile(configDir)

	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("unable to open %s: %w", path, err)
	}
	defer file.Close()

	var config Config

	if err := json.NewDecoder(file).Decode(&config); err != nil {
		return nil, fmt.Errorf("unable to decode the JSON data: %w", err)
	}

	return &config, nil
}

func FileExists(configDir string) (bool, error) {
	path := configFile(configDir)

	return utilities.FileExists(path)
}

func SaveDefaultConfigToFile(configDir string) error {
	path := configFile(configDir)

	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("unable to create the file at %s: %w", path, err)
	}
	defer file.Close()

	config := defaultConfig(configDir)
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "    ")

	if err := encoder.Encode(config); err != nil {
		return fmt.Errorf("unable to save the JSON data to the config file: %w", err)
	}

	return nil
}

func configFile(configDir string) string {
	return filepath.Join(utilities.CalculateConfigDir(configDir), configFileName)
}

func defaultConfig(configDir string) Config {
	credentialsFilePath := defaultCredentialsConfigFile(configDir)

	return Config{
		CredentialsFile: credentialsFilePath,
		CacheDirectory:  "",
		HTTP: HTTPConfig{
			Timeout:      5,
			MediaTimeout: 30,
		},
		LineWrapMaxWidth: 80,
		Integrations: Integrations{
			Browser:     "",
			Editor:      "",
			Pager:       "",
			ImageViewer: "",
			VideoPlayer: "",
		},
	}
}
