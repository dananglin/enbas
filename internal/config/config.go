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
	CredentialsFile  string       `json:"credentialsFile"`
	CacheDirectory   string       `json:"cacheDirectory"`
	LineWrapMaxWidth int          `json:"lineWrapMaxWidth"`
	GTSClient        GTSClient    `json:"gtsClient"`
	Server           Server       `json:"server"`
	Integrations     Integrations `json:"integrations"`
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

func NewConfigFromFile(configFilepath string) (Config, error) {
	path, err := configPath(configFilepath)
	if err != nil {
		return Config{}, fmt.Errorf("error calculating the path to your config file: %w", err)
	}

	file, err := utilities.OpenFile(path)
	if err != nil {
		return Config{}, fmt.Errorf("unable to open %s: %w", path, err)
	}
	defer file.Close()

	var config Config

	if err := json.NewDecoder(file).Decode(&config); err != nil {
		return Config{}, fmt.Errorf("unable to decode the JSON data: %w", err)
	}

	return config, nil
}

func FileExists(configFilepath string) (bool, error) {
	path, err := configPath(configFilepath)
	if err != nil {
		return false, fmt.Errorf("unable to calculate the path to your config file: %w", err)
	}

	return utilities.FileExists(path)
}

func EnsureParentDir(configFilepath string) error {
	path, err := configPath(configFilepath)
	if err != nil {
		return fmt.Errorf("error calculating the path to your config file: %w", err)
	}

	return utilities.EnsureDirectory(filepath.Dir(path))
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

	config := initialConfig()

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

	config.CredentialsFile = credentialsFilePath
	config.CacheDirectory = cacheDirPath
	config.Server.SocketPath = socketPath

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "    ")

	if err := encoder.Encode(config); err != nil {
		return fmt.Errorf("unable to save the JSON data to the config file: %w", err)
	}

	return nil
}

func initialConfig() Config {
	return Config{
		CredentialsFile: "",
		CacheDirectory:  "",
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
