package executor

import "fmt"

type Executor interface {
	Name() string
	Parse(args []string) error
	Execute() error
}

func Execute(executor Executor, args []string) error {
	if err := executor.Parse(args); err != nil {
		return fmt.Errorf("flag parsing error: %w", err)
	}

	if err := executor.Execute(); err != nil {
		return fmt.Errorf("execution error: %w", err)
	}

	return nil
}
