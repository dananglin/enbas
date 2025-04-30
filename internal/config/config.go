package config

import (
	"encoding/json"
	"fmt"
	"path/filepath"

	"codeflow.dananglin.me.uk/apollo/enbas/internal/utilities"
)

const (
	defaultHTTPTimeout       int = 5
	defaultHTTPMediaTimeout  int = 30
	defaultLineWrapMaxWidth  int = 80
	defaultServerIdleTimeout int = 300
)

type Config struct {
	populated        bool
	Path             string            `json:"-"`
	Aliases          map[string]string `json:"aliases"`
	CredentialsFile  string            `json:"credentialsFile"`
	CacheDirectory   string            `json:"cacheDirectory"`
	LineWrapMaxWidth int               `json:"lineWrapMaxWidth"`
	GTSClient        GTSClient         `json:"gtsClient"`
	Server           Server            `json:"server"`
	Integrations     Integrations      `json:"integrations"`
}

func NewConfigFromFile(configFilepath string) (Config, error) {
	return newConfigFromFile(configFilepath)
}

func (c Config) IsZero() bool {
	return !c.populated
}

type GTSClient struct {
	Timeout      int `json:"timeout"`
	MediaTimeout int `json:"mediaTimeout"`
}

type Server struct {
	SocketPath  string `json:"socketPath"`
	IdleTimeout int    `json:"idleTimeout"`
}

type Integrations struct {
	Browser     string `json:"browser"`
	Editor      string `json:"editor"`
	Pager       string `json:"pager"`
	ImageViewer string `json:"imageViewer"`
	VideoPlayer string `json:"videoPlayer"`
	AudioPlayer string `json:"audioPlayer"`
}

func newConfigFromFile(configFilepath string) (Config, error) {
	path, err := configPath(configFilepath)
	if err != nil {
		return Config{}, fmt.Errorf("error calculating the path to your config file: %w", err)
	}

	file, err := utilities.OpenFile(path)
	if err != nil {
		return Config{}, fmt.Errorf("error opening %s: %w", path, err)
	}
	defer file.Close()

	var cfg Config

	if err := json.NewDecoder(file).Decode(&cfg); err != nil {
		return Config{}, fmt.Errorf("error decoding the JSON data: %w", err)
	}

	cfg.Path = configFilepath
	cfg.populated = true

	return cfg, nil
}

func FileExists(configFilepath string) (bool, error) {
	path, err := configPath(configFilepath)
	if err != nil {
		return false, fmt.Errorf("error calculating the path to your configuration file: %w", err)
	}

	exists, err := utilities.FileExists(path)
	if err != nil {
		return false, fmt.Errorf("error checking if the configuration file is present: %w", err)
	}

	return exists, nil
}

func EnsureParentDir(configFilepath string) error {
	path, err := configPath(configFilepath)
	if err != nil {
		return fmt.Errorf("error calculating the path to your config file: %w", err)
	}

	if err := utilities.EnsureDirectory(filepath.Dir(path)); err != nil {
		return fmt.Errorf("error ensuring that the configuration file's parent directory is present: %w", err)
	}

	return nil
}

func SaveInitialConfigToFile(configFilepath string) error {
	path, err := configPath(configFilepath)
	if err != nil {
		return fmt.Errorf("unable to calculate the path to your config file: %w", err)
	}

	file, err := utilities.CreateFile(path)
	if err != nil {
		return fmt.Errorf("unable to create the file at %s: %w", path, err)
	}
	defer file.Close()

	cfg := initialConfig()

	credentialsFilePath, err := defaultCredentialsFilepath()
	if err != nil {
		return fmt.Errorf("unable to calculate the path to the credentials file: %w", err)
	}

	cacheDirPath, err := defaultCacheDir()
	if err != nil {
		return fmt.Errorf("error retrieving the path to the default cache directory: %w", err)
	}

	socketPath, err := newSocketPath()
	if err != nil {
		return fmt.Errorf("unable to calculate the path to the socket file: %w", err)
	}

	cfg.CredentialsFile = credentialsFilePath
	cfg.CacheDirectory = cacheDirPath
	cfg.Server.SocketPath = socketPath

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "    ")

	if err := encoder.Encode(cfg); err != nil {
		return fmt.Errorf("error writing the JSON data to the configuration file: %w", err)
	}

	return nil
}

func saveConfig(configFilepath string, cfg Config) error {
	file, err := utilities.CreateFile(configFilepath)
	if err != nil {
		return fmt.Errorf("error opening the configuration file: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "    ")

	if err := encoder.Encode(cfg); err != nil {
		return fmt.Errorf("error writing the JSON data to the configuration file: %w", err)
	}

	return nil
}

func initialConfig() Config {
	return Config{
		CredentialsFile: "",
		CacheDirectory:  "",
		Aliases:         make(map[string]string),
		GTSClient: GTSClient{
			Timeout:      defaultHTTPTimeout,
			MediaTimeout: defaultHTTPMediaTimeout,
		},
		Server: Server{
			SocketPath:  "",
			IdleTimeout: defaultServerIdleTimeout,
		},
		LineWrapMaxWidth: defaultLineWrapMaxWidth,
		Integrations: Integrations{
			Browser:     "",
			Editor:      "",
			Pager:       "",
			ImageViewer: "",
			VideoPlayer: "",
			AudioPlayer: "",
		},
	}
}
