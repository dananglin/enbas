package model

import (
	"fmt"
	"time"

	"codeflow.dananglin.me.uk/apollo/enbas/internal/utilities"
)

type Account struct {
	Acct            string      `json:"acct"`
	Avatar          string      `json:"avatar"`
	AvatarStatic    string      `json:"avatar_static"`
	CustomCSS       string      `json:"custom_css"`
	Header          string      `json:"header"`
	HeaderStatic    string      `json:"header_static"`
	ID              string      `json:"id"`
	LastStatusAt    string      `json:"last_status_at"`
	DisplayName     string      `json:"display_name"`
	Emojis          []Emoji     `json:"emojis"`
	EnableRSS       bool        `json:"enable_rss"`
	Bot             bool        `json:"bot"`
	Locked          bool        `json:"locked"`
	Suspended       bool        `json:"suspended"`
	Discoverable    bool        `json:"discoverable"`
	HideCollections bool        `json:"hide_collections"`
	Fields          []Field     `json:"fields"`
	FollowersCount  int         `json:"followers_count"`
	FollowingCount  int         `json:"following_count"`
	CreatedAt       time.Time   `json:"created_at"`
	MuteExpiresAt   time.Time   `json:"mute_expires_at"`
	Note            string      `json:"note"`
	Role            AccountRole `json:"role"`
	Source          Source      `json:"source"`
	StatusCount     int         `json:"statuses_count"`
	Theme           string      `json:"theme"`
	URL             string      `json:"url"`
	Username        string      `json:"username"`
}

type AccountRole struct {
	Name string `json:"name"`
}

type Source struct {
	Fields             []Field  `json:"fields"`
	FollowRequestCount int      `json:"follow_requests_count"`
	Language           string   `json:"language"`
	Note               string   `json:"note"`
	Privacy            string   `json:"string"`
	Sensitive          bool     `json:"sensitive"`
	StatusContentType  string   `json:"status_content_type"`
	AlsoKnownAsURIs    []string `json:"also_known_as_uris"`
}

type Field struct {
	Name       string `json:"name"`
	Value      string `json:"value"`
	VerifiedAt string `json:"verified_at"`
}

func (a Account) Display(noColor bool) string {
	format := `
%s (@%s)

%s
  %s

%s
  %s

%s
  %s %d
  %s %d
  %s %d

%s
  %s

%s %s

%s
  %s`

	metadata := ""

	for _, field := range a.Fields {
		metadata += fmt.Sprintf(
			"\n  %s: %s",
			utilities.FieldFormat(noColor, field.Name),
			utilities.StripHTMLTags(field.Value),
		)
	}

	return fmt.Sprintf(
		format,
		utilities.DisplayNameFormat(noColor, a.DisplayName),
		a.Username,
		utilities.HeaderFormat(noColor, "ACCOUNT ID:"),
		a.ID,
		utilities.HeaderFormat(noColor, "JOINED ON:"),
		utilities.FormatDate(a.CreatedAt),
		utilities.HeaderFormat(noColor, "STATS:"),
		utilities.FieldFormat(noColor, "Followers:"), a.FollowersCount,
		utilities.FieldFormat(noColor, "Following:"), a.FollowingCount,
		utilities.FieldFormat(noColor, "Statuses:"), a.StatusCount,
		utilities.HeaderFormat(noColor, "BIOGRAPHY:"),
		utilities.WrapLines(utilities.StripHTMLTags(a.Note), "\n  ", 80),
		utilities.HeaderFormat(noColor, "METADATA:"),
		metadata,
		utilities.HeaderFormat(noColor, "ACCOUNT URL:"),
		a.URL,
	)
}

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

func (a AccountRelationship) Display(noColor bool) string {
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
		utilities.HeaderFormat(noColor, "YOUR RELATIONSHIP WITH THIS ACCOUNT:"),
		utilities.FieldFormat(noColor, "Following"), a.Following,
		utilities.FieldFormat(noColor, "Is following you"), a.FollowedBy,
		utilities.FieldFormat(noColor, "A follow request was sent and is pending"), a.FollowRequested,
		utilities.FieldFormat(noColor, "Received a pending follow request"), a.FollowRequestedBy,
		utilities.FieldFormat(noColor, "Endorsed"), a.Endorsed,
		utilities.FieldFormat(noColor, "Showing Reposts (boosts)"), a.ShowingReblogs,
		utilities.FieldFormat(noColor, "Muted"), a.Muting,
		utilities.FieldFormat(noColor, "Notifications muted"), a.MutingNotifications,
		utilities.FieldFormat(noColor, "Blocking"), a.Blocking,
		utilities.FieldFormat(noColor, "Is blocking you"), a.BlockedBy,
		utilities.FieldFormat(noColor, "Blocking account's domain"), a.DomainBlocking,
	)

	if a.PrivateNote != "" {
		output += "\n"
		output += fmt.Sprintf(
			privateNoteFormat,
			utilities.HeaderFormat(noColor, "YOUR PRIVATE NOTE ABOUT THIS ACCOUNT:"),
			utilities.WrapLines(a.PrivateNote, "\n  ", 80),
		)
	}

	return output
}

type AccountListType int

const (
	AccountListFollowers AccountListType = iota
	AccountListFollowing
	AccountListBlockedAccount
)

type AccountList struct {
	Type     AccountListType
	Accounts []Account
}

func (a AccountList) Display(noColor bool) string {
	output := "\n"

	switch a.Type {
	case AccountListFollowers:
		output += utilities.HeaderFormat(noColor, "FOLLOWED BY:")
	case AccountListFollowing:
		output += utilities.HeaderFormat(noColor, "FOLLOWING:")
	case AccountListBlockedAccount:
		output += utilities.HeaderFormat(noColor, "BLOCKED ACCOUNTS:")
	default:
		output += utilities.HeaderFormat(noColor, "ACCOUNTS:")
	}

	if a.Type == AccountListBlockedAccount {
		for i := range a.Accounts {
			output += fmt.Sprintf(
				"\n  • %s (%s)",
				a.Accounts[i].Acct,
				a.Accounts[i].ID,
			)
		}
	} else {
		for i := range a.Accounts {
			output += fmt.Sprintf(
				"\n  • %s (%s)",
				utilities.DisplayNameFormat(noColor, a.Accounts[i].DisplayName),
				a.Accounts[i].Acct,
			)
		}
	}

	return output
}
