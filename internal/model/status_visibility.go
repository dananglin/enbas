// SPDX-FileCopyrightText: 2024 Dan Anglin <d.n.i.anglin@gmail.com>
//
// SPDX-License-Identifier: GPL-3.0-or-later

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

func ParseStatusVisibility(value string) (StatusVisibility, error) {
	mapped := map[string]StatusVisibility{
		statusVisibilityPublicValue:      StatusVisibilityPublic,
		statusVisibilityPrivateValue:     StatusVisibilityPrivate,
		statusVisibilityUnlistedValue:    StatusVisibilityUnlisted,
		statusVisibilityMutualsOnlyValue: StatusVisibilityMutualsOnly,
		statusVisibilityDirectValue:      StatusVisibilityDirect,
	}

	output, ok := mapped[value]
	if !ok {
		return StatusVisibilityUnknown, InvalidStatusVisibilityError{Value: value}
	}

	return output, nil
}

func (s StatusVisibility) MarshalJSON() ([]byte, error) {
	value := s.String()
	if value == unknownValue {
		return nil, InvalidStatusVisibilityError{Value: value}
	}

	data, err := json.Marshal(value)
	if err != nil {
		return nil, fmt.Errorf("unable to encode %s to JSON: %w", value, err)
	}

	return data, nil
}

func (s *StatusVisibility) UnmarshalJSON(data []byte) error {
	var (
		value string
		err   error
	)

	if err = json.Unmarshal(data, &value); err != nil {
		return fmt.Errorf("unable to decode the JSON data: %w", err)
	}

	// Ignore the error if the visibility from another service is
	// not known by enbas. It will be replaced with 'unknown'.
	*s, _ = ParseStatusVisibility(value)

	return nil
}

type InvalidStatusVisibilityError struct {
	Value string
}

func (e InvalidStatusVisibilityError) Error() string {
	return "'" + e.Value + "' is not a valid status visibility value: valid values are " +
		statusVisibilityPublicValue + ", " +
		statusVisibilityUnlistedValue + ", " +
		statusVisibilityPrivateValue + ", " +
		statusVisibilityMutualsOnlyValue + ", " +
		statusVisibilityDirectValue
}
