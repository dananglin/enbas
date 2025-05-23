package flag

import (
	"fmt"
	"maps"
	"slices"
	"strings"
)

type EnumValue struct {
	value     string
	validOpts map[string]struct{}
}

func NewEnumValue(
	options []string,
	defaultVal string,
) EnumValue {
	if len(options) == 0 {
		return EnumValue{
			value:     "",
			validOpts: nil,
		}
	}

	validOpts := make(map[string]struct{})

	for _, opt := range slices.All(options) {
		validOpts[opt] = struct{}{}
	}

	if defaultVal != "" {
		if _, ok := validOpts[defaultVal]; !ok {
			defaultVal = options[0]
		}
	}

	return EnumValue{
		value:     defaultVal,
		validOpts: validOpts,
	}
}

func (v *EnumValue) String() string {
	return v.value
}

func (v *EnumValue) Set(value string) error {
	if _, ok := v.validOpts[value]; !ok {
		return fmt.Errorf(
			"valid values are: %s",
			listOpts(v.validOpts),
		)
	}

	v.value = value

	return nil
}

func (v *EnumValue) Value() string {
	return v.value
}

type MultiEnumValue struct {
	values    []string
	validOpts map[string]struct{}
}

func NewMultiEnumValue(options []string) MultiEnumValue {
	if len(options) == 0 {
		return MultiEnumValue{
			values:    make([]string, 0),
			validOpts: nil,
		}
	}

	validOpts := make(map[string]struct{})

	for _, opt := range slices.All(options) {
		validOpts[opt] = struct{}{}
	}

	return MultiEnumValue{
		values:    make([]string, 0),
		validOpts: validOpts,
	}
}

func (v *MultiEnumValue) String() string {
	return strings.Join(v.values, ", ")
}

func (v *MultiEnumValue) Set(value string) error {
	if _, ok := v.validOpts[value]; !ok {
		return fmt.Errorf(
			"valid values are: %s",
			listOpts(v.validOpts),
		)
	}

	v.values = append(v.values, value)

	return nil
}

func (v *MultiEnumValue) Values() []string {
	return v.values
}

func (v *MultiEnumValue) Empty() bool {
	return len(v.values) == 0
}

func listOpts(validOpts map[string]struct{}) string {
	opts := []string{}

	for key := range maps.Keys(validOpts) {
		opts = append(opts, key)
	}

	slices.Sort(opts)

	output := ""
	numOpts := len(opts)

	for idx := range numOpts {
		if idx == numOpts-1 {
			output += opts[idx]

			continue
		}

		if idx == numOpts-2 {
			output += opts[idx] + ", and "

			continue
		}

		output += opts[idx] + ", "
	}

	return output
}
