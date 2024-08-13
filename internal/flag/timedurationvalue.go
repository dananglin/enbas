package flag

import (
	"fmt"
	"time"
)

type TimeDurationValue struct {
	Duration time.Duration
}

func NewTimeDurationValue() TimeDurationValue {
	return TimeDurationValue{
		Duration: 0 * time.Second,
	}
}

func (v TimeDurationValue) String() string {
	return v.Duration.String()
}

func (v *TimeDurationValue) Set(text string) error {
	duration, err := time.ParseDuration(text)
	if err != nil {
		return fmt.Errorf("unable to parse the value as time duration: %w", err)
	}

	v.Duration = duration

	return nil
}
