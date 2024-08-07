package model

const (
	TimelineCategoryHome   = "home"
	TimelineCategoryPublic = "public"
	TimelineCategoryTag    = "tag"
	TimelineCategoryList   = "list"
)

type InvalidTimelineCategoryError struct {
	Value string
}

func (e InvalidTimelineCategoryError) Error() string {
	return "'" +
		e.Value +
		"' is not a valid timeline category (valid values are " +
		TimelineCategoryHome + ", " +
		TimelineCategoryPublic + ", " +
		TimelineCategoryTag + ", " +
		TimelineCategoryList + ")"
}
