package flag

import "strings"

type StringSliceValue []string

func NewStringSliceValue() StringSliceValue {
	arr := make([]string, 0, 3)

	return StringSliceValue(arr)
}

func (v StringSliceValue) String() string {
	return strings.Join(v, ", ")
}

func (v *StringSliceValue) Set(value string) error {
	if len(value) > 0 {
		*v = append(*v, value)
	}

	return nil
}

func (v StringSliceValue) Empty() bool {
	return len(v) == 0
}

func (v StringSliceValue) Length() int {
	return len(v)
}

func (v StringSliceValue) ExpectedLength(expectedLength int) bool {
	return v.Length() == expectedLength
}
