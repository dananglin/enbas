package printer

import "fmt"

type NumFlagDescriptionsError struct {
	NumFlags            int
	NumFlagDescriptions int
}

func (e NumFlagDescriptionsError) Error() string {
	return fmt.Sprintf(
		"the number of flag descriptions does not match the number of flags: number of flags: %d, number of flag descriptions: %d",
		e.NumFlags,
		e.NumFlagDescriptions,
	)
}
