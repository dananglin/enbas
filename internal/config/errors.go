package config

type aliasAlreadyPresentError struct {
	alias string
}

func (e aliasAlreadyPresentError) Error() string {
	return "the alias '" + e.alias + "' is already present in your configuration"
}

func NewAliasAlreadyPresentError(alias string) error {
	return aliasAlreadyPresentError{alias: alias}
}

type aliasNotPresentError struct {
	alias string
}

func (e aliasNotPresentError) Error() string {
	return "the alias '" + e.alias + "' is not present in your configuration"
}

func NewAliasNotPresentError(alias string) error {
	return aliasNotPresentError{alias: alias}
}
