package main

type unknownCommandError struct {
	subcommand string
}

func (e unknownCommandError) Error() string {
	return "unknown command '" + e.subcommand + "'"
}
