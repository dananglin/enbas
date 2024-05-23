package executor

type FlagNotSetError struct {
	flagText string
}

func (e FlagNotSetError) Error() string {
	return "the flag '" + e.flagText + "' is not set"
}

type UnsupportedTypeError struct {
	resourceType string
}

func (e UnsupportedTypeError) Error() string {
	return "unsupported resource type '" + e.resourceType + "'"
}

type InvalidTimelineCategoryError struct {
	category string
}

func (e InvalidTimelineCategoryError) Error() string {
	return "'" + e.category + "' is not a valid timeline category (please choose home, public, tag or list)"
}

type NoAccountSpecifiedError struct{}

func (e NoAccountSpecifiedError) Error() string {
	return "no account specified in this request"
}

type UnsupportedAddOperationError struct {
	ResourceType      string
	AddToResourceType string
}

func (e UnsupportedAddOperationError) Error() string {
	return "adding '" + e.ResourceType + "' to '" + e.AddToResourceType + "' is not supported"
}

type UnsupportedRemoveOperationError struct {
	ResourceType           string
	RemoveFromResourceType string
}

func (e UnsupportedRemoveOperationError) Error() string {
	return "removing '" + e.ResourceType + "' from '" + e.RemoveFromResourceType + "' is not supported"
}

type EmptyContentError struct{}

func (e EmptyContentError) Error() string {
	return "content should not be empty"
}
