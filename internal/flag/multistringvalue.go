package flag

import "strings"

type MultiStringValue struct {
	values []string
}

func NewMultiStringValue() MultiStringValue {
	return MultiStringValue{
		values: make([]string, 0, 3),
	}
}

func (v *MultiStringValue) String() string {
	return strings.Join(v.values, ", ")
}

func (v *MultiStringValue) Set(value string) error {
	if value != "" {
		v.values = append(v.values, value)
	}

	return nil
}

func (v *MultiStringValue) Values() []string {
	return v.values
}

func (v *MultiStringValue) Empty() bool {
	return len(v.values) == 0
}

func (v *MultiStringValue) Length() int {
	return len(v.values)
}

func (v MultiStringValue) ExpectedLength(expectedLength int) bool {
	return v.Length() == expectedLength
}
