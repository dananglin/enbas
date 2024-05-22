package main

type flagNotSetError struct {
	flagText string
}

func (e flagNotSetError) Error() string {
	return "the flag '" + e.flagText + "' is not set"
}

type unsupportedResourceTypeError struct {
	resourceType string
}

func (e unsupportedResourceTypeError) Error() string {
	return "unsupported resource type '" + e.resourceType + "'"
}

type invalidTimelineCategoryError struct {
	category string
}

func (e invalidTimelineCategoryError) Error() string {
	return "'" + e.category + "' is not a valid timeline category (please choose home, public, tag or list)"
}

type unknownSubcommandError struct {
	subcommand string
}

func (e unknownSubcommandError) Error() string {
	return "unknown subcommand '" + e.subcommand + "'"
}

type noAccountSpecifiedError struct{}

func (e noAccountSpecifiedError) Error() string {
	return "no account specified in this request"
}

type unsupportedAddOperationError struct {
	ResourceType      string
	AddToResourceType string
}

func (e unsupportedAddOperationError) Error() string {
	return "adding '" + e.ResourceType + "' to '" + e.AddToResourceType + "' is not supported"
}

type unsupportedRemoveOperationError struct {
	ResourceType           string
	RemoveFromResourceType string
}

func (e unsupportedRemoveOperationError) Error() string {
	return "removing '" + e.ResourceType + "' from '" + e.RemoveFromResourceType + "' is not supported"
}
