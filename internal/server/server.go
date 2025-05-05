package server

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/rpc"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"codeflow.dananglin.me.uk/apollo/enbas/internal/gtsclient"
	"codeflow.dananglin.me.uk/apollo/enbas/internal/printer"
	"codeflow.dananglin.me.uk/apollo/enbas/internal/utilities"
)

const minIdleTimeout = 60

func Run(
	printSettings printer.Settings,
	client *gtsclient.GTSClient,
	socketPath string,
	withoutIdleTimeout bool,
	idleTimeout int,
) error {
	if socketPath == "" {
		return SocketFileNotSpecifiedError{}
	}

	socketPath, err := utilities.AbsolutePath(socketPath)
	if err != nil {
		return fmt.Errorf(
			"error calculating the absolute path to the socket file: %w",
			err,
		)
	}

	// Ensure that the socket file's parent folder is present.
	if err := utilities.EnsureDirectory(filepath.Dir(socketPath)); err != nil {
		return fmt.Errorf(
			"error ensuring the presence of the socket's parent directory: %w",
			err,
		)
	}

	if err := removeUnusedSocketFile(socketPath); err != nil {
		return fmt.Errorf("error removing the unused socket file: %w", err)
	}

	// Create the RPC server and register the GTS client methods.
	server := rpc.NewServer()

	if err := server.Register(client); err != nil {
		return fmt.Errorf("error registering the GTSClient methods to the server: %w", err)
	}

	if withoutIdleTimeout {
		// Run the server without a timer.
		return runWithoutIdleTimeout(
			printSettings,
			server,
			socketPath,
		)
	}

	// Run the server with a timer.
	return runWithIdleTimeout(
		printSettings,
		server,
		socketPath,
		idleTimeout,
	)
}

// runWithIdleTimeout runs the RPC server. The server closes after a specified amount of idle time or when the
// shutdown signal is received.
func runWithIdleTimeout(
	printSettings printer.Settings,
	server *rpc.Server,
	socketPath string,
	idleTimeout int,
) error {
	listener, err := net.Listen("unix", socketPath)
	if err != nil {
		return fmt.Errorf("error serving socket connection: %w", err)
	}
	defer listener.Close()

	printer.PrintInfo("Running the server using socket path: " + socketPath + "\n")

	// Create a timer for the idle timeout.
	if idleTimeout < minIdleTimeout {
		idleTimeout = minIdleTimeout
	}

	timeout := time.Duration(idleTimeout) * time.Second
	timer := time.NewTimer(timeout)

	// Listen and serve connections from the client in a separate goroutine.
	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				if errors.Is(err, net.ErrClosed) {
					printer.PrintInfo("Network connection closed.\n")

					return
				}

				printer.PrintFailure(
					printSettings,
					"Error accepting the connection: "+err.Error()+".",
				)

				os.Exit(1)
			}

			timer.Reset(timeout)

			go server.ServeConn(conn)
		}
	}()

	// Create a context for receiving the shutdown signal.
	ctx, stop := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	defer stop()

	select {
	case <-timer.C:
		printer.PrintInfo("Server idle timeout.\n")

		return nil
	case <-ctx.Done():
		stop()
		printer.PrintInfo("\nShutdown signal received.\n")

		return nil
	}
}

// runWithoutIdleTimeout runs the RPC server. The server closes when the shutdown signal is received.
func runWithoutIdleTimeout(
	printSettings printer.Settings,
	server *rpc.Server,
	socketPath string,
) error {
	listener, err := net.Listen("unix", socketPath)
	if err != nil {
		return fmt.Errorf("error serving socket connection: %w", err)
	}
	defer listener.Close()

	printer.PrintInfo("Running the server using socket path: " + socketPath + "\n")

	// Listen and serve connections from the client in a separate goroutine.
	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				if errors.Is(err, net.ErrClosed) {
					printer.PrintInfo("Network connection closed.\n")

					return
				}

				printer.PrintFailure(
					printSettings,
					"Error accepting the connection: "+err.Error()+".",
				)

				os.Exit(1)
			}

			go server.ServeConn(conn)
		}
	}()

	// Create a context for receiving the shutdown signal.
	ctx, stop := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	defer stop()

	<-ctx.Done()
	stop()

	printer.PrintInfo("\nShutdown signal received.\n")

	return nil
}

// removeUnusedSocketFile removes the socket file if it already exists and
// is not being used by a running server.
func removeUnusedSocketFile(path string) error {
	// Check for the existence of the socket path.
	exists, err := utilities.FileExists(path)
	if err != nil {
		return fmt.Errorf("error checking for the socket file: %w", err)
	}

	if !exists {
		return nil
	}

	// Attempt a connection to the socket path to see if it is in use.
	_, err = rpc.Dial("unix", path)

	// If the connection is successful, then the socket file is currently in
	// use by another running server.
	if err == nil {
		return SocketFileInUseError{}
	}

	// If no connection can be made then it should be safe to remove the file.
	if err := os.Remove(path); err != nil {
		return fmt.Errorf("error removing the unused socket file: %w", err)
	}

	return nil
}
