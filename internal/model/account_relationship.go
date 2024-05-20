package model

import (
	"fmt"

	"codeflow.dananglin.me.uk/apollo/enbas/internal/utilities"
)

type AccountRelationship struct {
	ID                  string `json:"id"`
	PrivateNote         string `json:"note"`
	BlockedBy           bool   `json:"blocked_by"`
	Blocking            bool   `json:"blocking"`
	DomainBlocking      bool   `json:"domain_blocking"`
	Endorsed            bool   `json:"endorsed"`
	FollowedBy          bool   `json:"followed_by"`
	Following           bool   `json:"following"`
	Muting              bool   `json:"muting"`
	MutingNotifications bool   `json:"muting_notifications"`
	Notifying           bool   `json:"notifying"`
	FollowRequested     bool   `json:"requested"`
	FollowRequestedBy   bool   `json:"requested_by"`
	ShowingReblogs      bool   `json:"showing_reblogs"`
}

func (a AccountRelationship) String() string {
	format := `
%s
  %s: %t
  %s: %t
  %s: %t
  %s: %t
  %s: %t
  %s: %t
  %s: %t
  %s: %t
  %s: %t
  %s: %t
  %s: %t`

	privateNoteFormat := `
%s
  %s`

	output := fmt.Sprintf(
		format,
		utilities.HeaderFormat("YOUR RELATIONSHIP WITH THIS ACCOUNT:"),
		utilities.FieldFormat("Following"), a.Following,
		utilities.FieldFormat("Is following you"), a.FollowedBy,
		utilities.FieldFormat("A follow request was sent and is pending"), a.FollowRequested,
		utilities.FieldFormat("Received a pending follow request"), a.FollowRequestedBy,
		utilities.FieldFormat("Endorsed"), a.Endorsed,
		utilities.FieldFormat("Showing Reposts (boosts)"), a.ShowingReblogs,
		utilities.FieldFormat("Muted"), a.Muting,
		utilities.FieldFormat("Notifications muted"), a.MutingNotifications,
		utilities.FieldFormat("Blocking"), a.Blocking,
		utilities.FieldFormat("Is blocking you"), a.BlockedBy,
		utilities.FieldFormat("Blocking account's domain"), a.DomainBlocking,
	)

	if a.PrivateNote != "" {
		output += fmt.Sprintf(
			privateNoteFormat,
			utilities.HeaderFormat("YOUR PRIVATE NOTE ABOUT THIS ACCOUNT:"),
			utilities.WrapLines(a.PrivateNote, "\n  ", 80),
		)
	}

	return output
}
