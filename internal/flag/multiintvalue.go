package flag

import (
	"fmt"
	"slices"
	"strconv"
	"strings"
)

type MultiIntValue struct {
	values []int
}

func NewMultiIntValue() MultiIntValue {
	return MultiIntValue{
		values: make([]int, 0, 3),
	}
}

func (v *MultiIntValue) String() string {
	var builder strings.Builder

	for idx, value := range slices.All(v.values) {
		if idx == len(v.values)-1 {
			builder.WriteString(strconv.Itoa(value))
		} else {
			builder.WriteString(strconv.Itoa(value))
			builder.WriteString(", ")
		}
	}

	return builder.String()
}

func (v *MultiIntValue) Set(text string) error {
	value, err := strconv.Atoi(text)
	if err != nil {
		return fmt.Errorf("error parsing the value to an integer: %w", err)
	}

	v.values = append(v.values, value)

	return nil
}

func (v *MultiIntValue) Values() []int {
	return v.values
}

func (v *MultiIntValue) Empty() bool {
	return len(v.values) == 0
}

func (v *MultiIntValue) Length() int {
	return len(v.values)
}

func (v *MultiIntValue) ExpectedLength(expectedLength int) bool {
	return v.Length() == expectedLength
}
