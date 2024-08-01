package model

import (
	"encoding/json"
	"fmt"
)

type ListRepliesPolicy int

const (
	ListRepliesPolicyFollowed ListRepliesPolicy = iota
	ListRepliesPolicyList
	ListRepliesPolicyNone
	ListRepliesPolicyUnknown
)

const (
	listRepliesPolicyFollowedValue = "followed"
	listRepliesPolicyListValue     = "list"
	listRepliesPolicyNoneValue     = "none"
)

func (l ListRepliesPolicy) String() string {
	mapped := map[ListRepliesPolicy]string{
		ListRepliesPolicyFollowed: listRepliesPolicyFollowedValue,
		ListRepliesPolicyList:     listRepliesPolicyListValue,
		ListRepliesPolicyNone:     listRepliesPolicyNoneValue,
	}

	output, ok := mapped[l]
	if !ok {
		return unknownValue
	}

	return output
}

func ParseListRepliesPolicy(value string) (ListRepliesPolicy, error) {
	mapped := map[string]ListRepliesPolicy{
		listRepliesPolicyFollowedValue: ListRepliesPolicyFollowed,
		listRepliesPolicyListValue:     ListRepliesPolicyList,
		listRepliesPolicyNoneValue:     ListRepliesPolicyNone,
	}

	output, ok := mapped[value]
	if !ok {
		return ListRepliesPolicyUnknown, InvalidListRepliesPolicyError{value}
	}

	return output, nil
}

func (l ListRepliesPolicy) MarshalJSON() ([]byte, error) {
	value := l.String()
	if value == unknownValue {
		return nil, InvalidListRepliesPolicyError{value}
	}

	data, err := json.Marshal(value)
	if err != nil {
		return nil, fmt.Errorf("unable to encode %s to JSON: %w", value, err)
	}

	return data, nil
}

func (l *ListRepliesPolicy) UnmarshalJSON(data []byte) error {
	var (
		value string
		err   error
	)

	if err = json.Unmarshal(data, &value); err != nil {
		return fmt.Errorf("unable to decode the JSON data: %w", err)
	}

	*l, err = ParseListRepliesPolicy(value)
	if err != nil {
		return err
	}

	return nil
}

type InvalidListRepliesPolicyError struct {
	Value string
}

func (e InvalidListRepliesPolicyError) Error() string {
	return "'" +
		e.Value +
		"' is not a valid list replies policy: valid values are " +
		listRepliesPolicyFollowedValue + ", " +
		listRepliesPolicyListValue + ", " +
		listRepliesPolicyNoneValue
}

type List struct {
	ID            string            `json:"id"`
	RepliesPolicy ListRepliesPolicy `json:"replies_policy"`
	Title         string            `json:"title"`
	Accounts      map[string]string
}
