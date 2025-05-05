package server

type SocketFileInUseError struct{}

func (e SocketFileInUseError) Error() string {
	return "this socket file is used by another server"
}

type SocketFileNotSpecifiedError struct{}

func (e SocketFileNotSpecifiedError) Error() string {
	return "the path to the socket file is not specified"
}
