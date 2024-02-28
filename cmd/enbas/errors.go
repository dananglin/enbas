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
