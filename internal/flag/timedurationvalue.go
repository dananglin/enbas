package flag

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const timeDurationRegexPattern string = `[0-9]{1,4}\s+(days?|hours?|minutes?|seconds?)`

type TimeDurationValue struct {
	duration time.Duration
	isSet    bool
}

func NewTimeDurationValue(defaultDuration time.Duration) TimeDurationValue {
	return TimeDurationValue{
		duration: defaultDuration,
		isSet:    false,
	}
}

func (v *TimeDurationValue) String() string {
	return v.duration.String()
}

func (v *TimeDurationValue) Set(value string) error {
	pattern := regexp.MustCompile(timeDurationRegexPattern)
	matches := pattern.FindAllString(value, -1)

	days, hours, minutes, seconds := 0, 0, 0, 0

	var err error

	for idx := range matches {
		switch {
		case strings.Contains(matches[idx], "day"):
			days, err = parseInt(matches[idx])
			if err != nil {
				return fmt.Errorf(
					"unable to parse the number of days from %s: %w",
					matches[idx],
					err,
				)
			}
		case strings.Contains(matches[idx], "hour"):
			hours, err = parseInt(matches[idx])
			if err != nil {
				return fmt.Errorf(
					"unable to parse the number of hours from %s: %w",
					matches[idx],
					err,
				)
			}
		case strings.Contains(matches[idx], "minute"):
			minutes, err = parseInt(matches[idx])
			if err != nil {
				return fmt.Errorf(
					"unable to parse the number of minutes from %s: %w",
					matches[idx],
					err,
				)
			}
		case strings.Contains(matches[idx], "second"):
			seconds, err = parseInt(matches[idx])
			if err != nil {
				return fmt.Errorf(
					"unable to parse the number of seconds from %s: %w",
					matches[idx],
					err,
				)
			}
		}
	}

	durationValue := (days * 86400) + (hours * 3600) + (minutes * 60) + seconds

	v.duration = time.Duration(durationValue) * time.Second
	v.isSet = true

	return nil
}

func (v *TimeDurationValue) IsSet() bool {
	return v.isSet
}

func (v *TimeDurationValue) Value() time.Duration {
	return v.duration
}

func parseInt(text string) (int, error) {
	split := strings.SplitN(text, " ", 2)
	if len(split) != 2 {
		return 0, fmt.Errorf("unexpected number of split for %s: want 2, got %d", text, len(split))
	}

	output, err := strconv.Atoi(split[0])
	if err != nil {
		return 0, fmt.Errorf("unable to convert %s to an integer: %w", text, err)
	}

	return output, nil
}
