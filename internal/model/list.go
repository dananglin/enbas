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

	return ListRepliesPolicy(-1), fmt.Errorf("%q is not a valid list replies policy", policy)
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
	value := l.String()
	if value == "" {
		return nil, errors.New("invalid list replies policy")
	}

	return json.Marshal(value)
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
	Accounts      map[string]string
}

func (l List) String() string {
	format := `
%s
  %s

%s
  %s

%s
  %s

%s`

	output := fmt.Sprintf(
		format,
		utilities.HeaderFormat("LIST TITLE:"), l.Title,
		utilities.HeaderFormat("LIST ID:"), l.ID,
		utilities.HeaderFormat("REPLIES POLICY:"), l.RepliesPolicy,
		utilities.HeaderFormat("ADDED ACCOUNTS:"),
	)

	if len(l.Accounts) > 0 {
		for id, name := range l.Accounts {
			output += fmt.Sprintf(
				"\n  • %s (%s)",
				utilities.DisplayNameFormat(name),
				id,
			)
		}
	} else {
		output += "\n  None"
	}

	return output
}

type Lists []List

func (l Lists) String() string {
	output := ""

	for i := range l {
		output += fmt.Sprintf(
			"\n%s (%s)",
			l[i].Title,
			l[i].ID,
		)
	}

	return output
}
