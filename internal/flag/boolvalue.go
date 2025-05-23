package flag

import (
	"fmt"
	"strconv"
)

type BoolValue struct {
	value bool
	isSet bool
}

func NewBoolValue(defaultVal bool) BoolValue {
	return BoolValue{
		value: defaultVal,
		isSet: false,
	}
}

func (v *BoolValue) Set(value string) error {
	boolVal, err := strconv.ParseBool(value)
	if err != nil {
		return fmt.Errorf("error parsing %q as a boolean value: %w", value, err)
	}

	v.value = boolVal
	v.isSet = true

	return nil
}

func (v *BoolValue) IsSet() bool {
	return v.isSet
}

func (v *BoolValue) IsBoolFlag() bool { return true }

func (v *BoolValue) Value() bool {
	return v.value
}

func (v *BoolValue) String() string {
	return strconv.FormatBool(v.value)
}
