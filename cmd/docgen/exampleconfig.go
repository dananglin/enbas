package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"codeflow.dananglin.me.uk/apollo/enbas/internal/config"
)

func generateExampleConfig(
	applicationName string,
	dir string,
) error {
	path := filepath.Join(dir, "config.json")

	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("error creating %q: %w", path, err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "    ")

	if err := encoder.Encode(exampleConfig(applicationName)); err != nil {
		return fmt.Errorf("error saving the example configuration to %q: %w", path, err)
	}

	return nil
}

func exampleConfig(applicationName string) config.Config {
	return config.Config{
		CredentialsFile: "/home/user/.local/config/" + applicationName + "/credentials/credentials.json",
		CacheDirectory:  "/home/user/.local/cache/" + applicationName,
		Aliases: map[string]string{
			"aliases":      "show aliases",
			"boost":        "reblog status --status-id",
			"fave":         "favourite status --status-id",
			"my-followers": "show followers from account --my-account",
			"toot":         "create status --content-type plain --visibility public --content",
		},
		LineWrapMaxWidth: 80,
		GTSClient: config.GTSClient{
			Timeout:      30,
			MediaTimeout: 60,
		},
		Server: config.Server{
			SocketPath:  "/var/run/user/1000/" + applicationName + "/server.psqm2yeo.socket",
			IdleTimeout: 300,
		},
		Integrations: config.Integrations{
			Browser:     "firefox --new-window",
			Editor:      "vim",
			Pager:       "less -FIRX",
			ImageViewer: "feh --scale-down",
			VideoPlayer: "mpv --loop-file=inf",
			AudioPlayer: "mpv --force-window",
		},
	}
}
