package config

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"codeflow.dananglin.me.uk/apollo/enbas/internal/info"
	"codeflow.dananglin.me.uk/apollo/enbas/internal/utilities"
)

const (
	configFileName string = "config.json"

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

func NewConfigFromFile(configDir string) (*Config, error) {
	path, err := configPath(configDir)
	if err != nil {
		return nil, fmt.Errorf("unable to calculate the path to your config file: %w", err)
	}

	file, err := utilities.OpenFile(path)
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

func SaveInitialConfigToFile(configDir string) error {
	path, err := configPath(configDir)
	if err != nil {
		return fmt.Errorf("unable to calculate the path to your config file: %w", err)
	}

	file, err := utilities.CreateFile(path)
	if err != nil {
		return fmt.Errorf("unable to create the file at %s: %w", path, err)
	}
	defer file.Close()

	config := initialConfig()

	credentialsFilePath, err := defaultCredentialsConfigFile(configDir)
	if err != nil {
		return fmt.Errorf("unable to calculate the path to the credentials file: %w", err)
	}

	socketPath, err := createSocketPath()
	if err != nil {
		return fmt.Errorf("unable to calculate the path to the socket file: %w", err)
	}

	config.CredentialsFile = credentialsFilePath
	config.Server.SocketPath = socketPath

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

func createSocketPath() (string, error) {
	runtimeDir := os.Getenv("XDG_RUNTIME_DIR")
	if runtimeDir == "" {
		return "", nil
	}

	randBytes := make([]byte, 4)

	if _, err := rand.Read(randBytes); err != nil {
		return "", fmt.Errorf("unable to create random bytes: %w", err)
	}

	path, err := utilities.AbsolutePath(filepath.Join(
		runtimeDir,
		info.ApplicationName,
		"server."+hex.EncodeToString(randBytes)+".socket",
	))
	if err != nil {
		return "", fmt.Errorf("unable to calculate the absolute path to the socket file: %w", err)
	}

	return path, nil
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
		},
	}
}
