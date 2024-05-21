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

func (a Account) String() string {
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
			utilities.FieldFormat(field.Name),
			utilities.StripHTMLTags(field.Value),
		)
	}

	return fmt.Sprintf(
		format,
		utilities.DisplayNameFormat(a.DisplayName),
		a.Username,
		utilities.HeaderFormat("ACCOUNT ID:"),
		a.ID,
		utilities.HeaderFormat("JOINED ON:"),
		utilities.FormatDate(a.CreatedAt),
		utilities.HeaderFormat("STATS:"),
		utilities.FieldFormat("Followers:"), a.FollowersCount,
		utilities.FieldFormat("Following:"), a.FollowingCount,
		utilities.FieldFormat("Statuses:"), a.StatusCount,
		utilities.HeaderFormat("BIOGRAPHY:"),
		utilities.WrapLines(utilities.StripHTMLTags(a.Note), "\n  ", 80),
		utilities.HeaderFormat("METADATA:"),
		metadata,
		utilities.HeaderFormat("ACCOUNT URL:"),
		a.URL,
	)
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

func (a AccountList) String() string {
	output := "\n"

	switch a.Type {
	case AccountListFollowers:
		output += utilities.HeaderFormat("FOLLOWED BY:")
	case AccountListFollowing:
		output += utilities.HeaderFormat("FOLLOWING:")
	case AccountListBlockedAccount:
		output += utilities.HeaderFormat("BLOCKED ACCOUNTS:")
	default:
		output += utilities.HeaderFormat("ACCOUNTS:")
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
				utilities.DisplayNameFormat(a.Accounts[i].DisplayName),
				a.Accounts[i].Acct,
			)
		}
	}

	return output
}
