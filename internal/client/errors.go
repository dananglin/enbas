package client

import "fmt"

type ResponseError struct {
	StatusCode       int
	Message          string
	MessageDecodeErr error
}

func (e ResponseError) Error() string {
	if e.MessageDecodeErr != nil {
		return fmt.Sprintf(
			"received HTTP code %d from the instance but was unable to decode the error message: %v",
			e.StatusCode,
			e.MessageDecodeErr,
		)
	}

	if e.Message == "" {
		return fmt.Sprintf(
			"received HTTP code %d from the instance but no error message was provided",
			e.StatusCode,
		)
	}

	return fmt.Sprintf(
		"message received from the instance: (%d) %q",
		e.StatusCode,
		e.Message,
	)
}

type BadStatusCodeError struct {
	statusCode int
}

func (e BadStatusCodeError) Error() string {
	return fmt.Sprintf(
		"did not receive an OK response from the GoToSocial server: got %d",
		e.statusCode,
	)
}

type Error struct {
	message string
}

func (e Error) Error() string {
	return e.message
}
