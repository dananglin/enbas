package flag

import (
	"fmt"
	"strconv"
)

type BoolPtrValue struct {
	Value *bool
}

func NewBoolPtrValue() BoolPtrValue {
	return BoolPtrValue{
		Value: nil,
	}
}

func (b BoolPtrValue) String() string {
	if b.Value == nil {
		return "NOT SET"
	}

	return strconv.FormatBool(*b.Value)
}

func (b *BoolPtrValue) Set(value string) error {
	boolVar, err := strconv.ParseBool(value)
	if err != nil {
		return fmt.Errorf("unable to parse %q as a boolean value: %w", value, err)
	}

	b.Value = new(bool)
	*b.Value = boolVar

	return nil
}

func (b *BoolPtrValue) IsBoolFlag() bool { return true }
