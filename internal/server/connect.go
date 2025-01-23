package server

import (
	"fmt"
	"net/rpc"
	"os"
	"os/exec"
	"time"

	"codeflow.dananglin.me.uk/apollo/enbas/internal/config"
	"codeflow.dananglin.me.uk/apollo/enbas/internal/utilities"
)

// Connect attempts to connect to the server and return
// the reference to the created RPC client. If the server is not
// running then an attempt is made to run and connect to a
// temporary server.
func Connect(cfg config.Server, cfgDir string) (*rpc.Client, error) {
	socketPath, err := utilities.AbsolutePath(cfg.SocketPath)
	if err != nil {
		return nil, fmt.Errorf(
			"unable to calculate the absolute path of the socket file: %w",
			err,
		)
	}

	exists, err := utilities.FileExists(socketPath)
	if err != nil {
		return nil, fmt.Errorf(
			"received an unexpected error after checking for the socket file: %w",
			err,
		)
	}

	// The function to attempt the connection to the server.
	connect := func() (*rpc.Client, error) {
		client, err := rpc.Dial("unix", socketPath)
		if err != nil {
			return nil, fmt.Errorf("error connecting to the RPC server: %w", err)
		}

		return client, nil
	}

	// Attempt the server connection if the socket file is present.
	if exists {
		return connect()
	}

	// The socket file is not present so we will
	// attempt to start a new server process.
	server := exec.Command(os.Args[0], "--config-dir", cfgDir, "server")

	if err := server.Start(); err != nil {
		return nil, fmt.Errorf("error starting the server: %w", err)
	}

	// Attempt to create a connection to the server.
	var client *rpc.Client

	for range 3 {
		time.Sleep(100 * time.Millisecond)

		client, err = connect()
		if err == nil {
			break
		}
	}

	if err != nil {
		return nil, fmt.Errorf("error connecting to the server: %w", err)
	}

	return client, nil
}
