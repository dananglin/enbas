package flag

import (
	"fmt"
	"strconv"
	"strings"
)

type IntSliceValue []int

func NewIntSliceValue() IntSliceValue {
	arr := make([]int, 0, 3)

	return IntSliceValue(arr)
}

func (v IntSliceValue) String() string {
	var builder strings.Builder

	for ind, value := range v {
		if ind == len(v)-1 {
			builder.WriteString(strconv.Itoa(value))
		} else {
			builder.WriteString(strconv.Itoa(value))
			builder.WriteString(", ")
		}
	}

	return builder.String()
}

func (v *IntSliceValue) Set(text string) error {
	value, err := strconv.Atoi(text)
	if err != nil {
		return fmt.Errorf("unable to parse the value to an integer: %w", err)
	}

	*v = append(*v, value)

	return nil
}

func (v IntSliceValue) Empty() bool {
	return len(v) == 0
}

func (v IntSliceValue) ExpectedLength(expectedLength int) bool {
	return len(v) == expectedLength
}
