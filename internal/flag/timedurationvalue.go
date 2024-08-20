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
	Duration time.Duration
}

func NewTimeDurationValue() TimeDurationValue {
	return TimeDurationValue{
		Duration: time.Duration(0),
	}
}

func (v TimeDurationValue) String() string {
	return v.Duration.String()
}

func (v *TimeDurationValue) Set(value string) error {
	pattern := regexp.MustCompile(timeDurationRegexPattern)
	matches := pattern.FindAllString(value, -1)

	days, hours, minutes, seconds := 0, 0, 0, 0

	var err error

	for ind := range len(matches) {
		switch {
		case strings.Contains(matches[ind], "day"):
			days, err = parseInt(matches[ind])
			if err != nil {
				return fmt.Errorf("unable to parse the number of days from %s: %w", matches[ind], err)
			}
		case strings.Contains(matches[ind], "hour"):
			hours, err = parseInt(matches[ind])
			if err != nil {
				return fmt.Errorf("unable to parse the number of hours from %s: %w", matches[ind], err)
			}
		case strings.Contains(matches[ind], "minute"):
			minutes, err = parseInt(matches[ind])
			if err != nil {
				return fmt.Errorf("unable to parse the number of minutes from %s: %w", matches[ind], err)
			}
		case strings.Contains(matches[ind], "second"):
			seconds, err = parseInt(matches[ind])
			if err != nil {
				return fmt.Errorf("unable to parse the number of seconds from %s: %w", matches[ind], err)
			}
		}
	}

	durationValue := (days * 86400) + (hours * 3600) + (minutes * 60) + seconds

	v.Duration = time.Duration(durationValue) * time.Second

	return nil
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
