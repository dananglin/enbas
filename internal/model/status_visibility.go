package model

import (
	"encoding/json"
	"fmt"
)

type StatusVisibility int

const (
	StatusVisibilityPublic StatusVisibility = iota
	StatusVisibilityPrivate
	StatusVisibilityUnlisted
	StatusVisibilityMutualsOnly
	StatusVisibilityDirect
	StatusVisibilityUnknown
)

const (
	statusVisibilityPublicValue      = "public"
	statusVisibilityPrivateValue     = "private"
	statusVisibilityUnlistedValue    = "unlisted"
	statusVisibilityMutualsOnlyValue = "mutuals_only"
	statusVisibilityDirectValue      = "direct"
)

func (s StatusVisibility) String() string {
	mapped := map[StatusVisibility]string{
		StatusVisibilityPublic:      statusVisibilityPublicValue,
		StatusVisibilityPrivate:     statusVisibilityPrivateValue,
		StatusVisibilityUnlisted:    statusVisibilityUnlistedValue,
		StatusVisibilityMutualsOnly: statusVisibilityMutualsOnlyValue,
		StatusVisibilityDirect:      statusVisibilityDirectValue,
	}

	output, ok := mapped[s]
	if !ok {
		return unknownValue
	}

	return output
}

func ParseStatusVisibility(value string) StatusVisibility {
	mapped := map[string]StatusVisibility{
		statusVisibilityPublicValue:      StatusVisibilityPublic,
		statusVisibilityPrivateValue:     StatusVisibilityPrivate,
		statusVisibilityUnlistedValue:    StatusVisibilityUnlisted,
		statusVisibilityMutualsOnlyValue: StatusVisibilityMutualsOnly,
		statusVisibilityDirectValue:      StatusVisibilityDirect,
	}

	output, ok := mapped[value]
	if !ok {
		return StatusVisibilityUnknown
	}

	return output
}

func (s StatusVisibility) MarshalJSON() ([]byte, error) {
	value := s.String()
	if value == unknownValue {
		return nil, fmt.Errorf("%q is not a valid status visibility", value)
	}

	return json.Marshal(value)
}

func (s *StatusVisibility) UnmarshalJSON(data []byte) error {
	var (
		value string
		err   error
	)

	if err = json.Unmarshal(data, &value); err != nil {
		return fmt.Errorf("unable to unmarshal the data; %w", err)
	}

	*s = ParseStatusVisibility(value)

	return nil
}
