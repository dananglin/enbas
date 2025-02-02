package executor

import (
	"fmt"

	"codeflow.dananglin.me.uk/apollo/enbas/internal/gtsclient"
	"codeflow.dananglin.me.uk/apollo/enbas/internal/server"
)

func (s *ServerExecutor) Execute() error {
	gtsClient, err := gtsclient.NewGTSClient(s.config)
	if err != nil {
		return fmt.Errorf("unable to create the GoToSocial client: %w", err)
	}

	if err := server.Run(
		s.printSettings,
		gtsClient,
		s.config.Server.SocketPath,
		s.noIdleTimeout,
		s.config.Server.IdleTimeout,
	); err != nil {
		return fmt.Errorf("error running the daemon process: %w", err)
	}

	return nil
}
