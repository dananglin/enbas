package model

import (
	"encoding/json"
	"fmt"

	"codeflow.dananglin.me.uk/apollo/enbas/internal/utilities"
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

func ParseListRepliesPolicy(value string) ListRepliesPolicy {
	mapped := map[string]ListRepliesPolicy{
		listRepliesPolicyFollowedValue: ListRepliesPolicyFollowed,
		listRepliesPolicyListValue:     ListRepliesPolicyList,
		listRepliesPolicyNoneValue:     ListRepliesPolicyNone,
	}

	output, ok := mapped[value]
	if !ok {
		return ListRepliesPolicyUnknown
	}

	return output
}

func (l ListRepliesPolicy) MarshalJSON() ([]byte, error) {
	value := l.String()
	if value == unknownValue {
		return nil, fmt.Errorf("%q is not a valid list replies policy")
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

	*l = ParseListRepliesPolicy(value)

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
