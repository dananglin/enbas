package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"codeflow.dananglin.me.uk/apollo/enbas/internal/utilities"
)

const (
	configFileName string = "config.json"

	defaultHTTPTimeout      int = 5
	defaultHTTPMediaTimeout int = 30
	defaultLineWrapMaxWidth int = 80
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
	path, err := configPath(configDir)
	if err != nil {
		return nil, fmt.Errorf("unable to calculate the path to your config file: %w", err)
	}

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
	path, err := configPath(configDir)
	if err != nil {
		return false, fmt.Errorf("unable to calculate the path to your config file: %w", err)
	}

	return utilities.FileExists(path)
}

func SaveDefaultConfigToFile(configDir string) error {
	path, err := configPath(configDir)
	if err != nil {
		return fmt.Errorf("unable to calculate the path to your config file: %w", err)
	}

	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("unable to create the file at %s: %w", path, err)
	}
	defer file.Close()

	config := defaultConfig()

	credentialsFilePath, err := defaultCredentialsConfigFile(configDir)
	if err != nil {
		return fmt.Errorf("unable to calculate the path to the credentials file: %w", err)
	}

	config.CredentialsFile = credentialsFilePath

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "    ")

	if err := encoder.Encode(config); err != nil {
		return fmt.Errorf("unable to save the JSON data to the config file: %w", err)
	}

	return nil
}

func configPath(configDir string) (string, error) {
	configDir, err := utilities.CalculateConfigDir(configDir)
	if err != nil {
		return "", fmt.Errorf("unable to get the config directory: %w", err)
	}

	return filepath.Join(configDir, configFileName), nil
}

func defaultConfig() Config {
	return Config{
		CredentialsFile: "",
		CacheDirectory:  "",
		HTTP: HTTPConfig{
			Timeout:      defaultHTTPTimeout,
			MediaTimeout: defaultHTTPMediaTimeout,
		},
		LineWrapMaxWidth: defaultLineWrapMaxWidth,
		Integrations: Integrations{
			Browser:     "",
			Editor:      "",
			Pager:       "",
			ImageViewer: "",
			VideoPlayer: "",
		},
	}
}
