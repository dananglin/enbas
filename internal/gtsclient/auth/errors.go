package auth

type NoAuthenticationDetailsError struct {
	account string
}

func NewNoAuthenticationDetailsError(account string) NoAuthenticationDetailsError {
	return NoAuthenticationDetailsError{account: account}
}

func (e NoAuthenticationDetailsError) Error() string {
	return "no authentication details found for the current account '" +
		e.account +
		"'"
}
