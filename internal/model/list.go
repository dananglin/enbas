package model

import (
	"encoding/json"
	"errors"
	"fmt"

	"codeflow.dananglin.me.uk/apollo/enbas/internal/utilities"
)

type ListRepliesPolicy int

const (
	ListRepliesPolicyFollowed ListRepliesPolicy = iota
	ListRepliesPolicyList
	ListRepliesPolicyNone
)

func ParseListRepliesPolicy(policy string) (ListRepliesPolicy, error) {
	switch policy {
	case "followed":
		return ListRepliesPolicyFollowed, nil
	case "list":
		return ListRepliesPolicyList, nil
	case "none":
		return ListRepliesPolicyNone, nil
	}

	return ListRepliesPolicy(-1), errors.New("invalid list replies policy")
}

func (l ListRepliesPolicy) String() string {
	switch l {
	case ListRepliesPolicyFollowed:
		return "followed"
	case ListRepliesPolicyList:
		return "list"
	case ListRepliesPolicyNone:
		return "none"
	}

	return ""
}

func (l ListRepliesPolicy) MarshalJSON() ([]byte, error) {
	str := l.String()
	if str == "" {
		return nil, errors.New("invalid list replies policy")
	}

	return json.Marshal(str)
}

func (l *ListRepliesPolicy) UnmarshalJSON(data []byte) error {
	var (
		value string
		err   error
	)

	if err = json.Unmarshal(data, &value); err != nil {
		return fmt.Errorf("unable to unmarshal the data; %w", err)
	}

	*l, err = ParseListRepliesPolicy(value)
	if err != nil {
		return fmt.Errorf("unable to parse %s as a list replies policy; %w", value, err)
	}

	return nil
}

type List struct {
	ID            string            `json:"id"`
	RepliesPolicy ListRepliesPolicy `json:"replies_policy"`
	Title         string            `json:"title"`
}

func (l List) String() string {
	format := `%s %s
%s %s
%s %s`

	return fmt.Sprintf(
		format,
		utilities.FieldFormat("List ID:"), l.ID,
		utilities.FieldFormat("Title:"), l.Title,
		utilities.FieldFormat("Replies Policy:"), l.RepliesPolicy,
	)
}
