package model

import (
	"time"
)

type Account struct {
	Acct              string               `json:"acct"`
	Avatar            string               `json:"avatar"`
	AvatarDescription string               `json:"avatar_description"`
	AvatarMediaID     string               `json:"avatar_media_id"`
	AvatarStatic      string               `json:"avatar_static"`
	CustomCSS         string               `json:"custom_css"`
	Header            string               `json:"header"`
	HeaderDescription string               `json:"header_description"`
	HeaderMediaID     string               `json:"header_media_id"`
	HeaderStatic      string               `json:"header_static"`
	ID                string               `json:"id"`
	LastStatusAt      string               `json:"last_status_at"`
	DisplayName       string               `json:"display_name"`
	Emojis            []Emoji              `json:"emojis"`
	EnableRSS         bool                 `json:"enable_rss"`
	Bot               bool                 `json:"bot"`
	Locked            bool                 `json:"locked"`
	Suspended         bool                 `json:"suspended"`
	Discoverable      bool                 `json:"discoverable"`
	HideCollections   bool                 `json:"hide_collections"`
	Fields            []Field              `json:"fields"`
	FollowersCount    int                  `json:"followers_count"`
	FollowingCount    int                  `json:"following_count"`
	CreatedAt         time.Time            `json:"created_at"`
	MuteExpiresAt     time.Time            `json:"mute_expires_at"`
	Note              string               `json:"note"`
	Role              AccountRole          `json:"role"`
	Roles             []AccountDisplayRole `json:"roles"`
	Source            Source               `json:"source"`
	StatusCount       int                  `json:"statuses_count"`
	Theme             string               `json:"theme"`
	URL               string               `json:"url"`
	Username          string               `json:"username"`
}

type AccountRole struct {
	Color       string `json:"color"`
	Highlighted bool   `json:"highlighted"`
	ID          string `json:"id"`
	Name        string `json:"name"`
	Permissions string `json:"permissions"`
}

type AccountDisplayRole struct {
	Color string `json:"color"`
	ID    string `json:"id"`
	Name  string `json:"name"`
}

type Source struct {
	Fields             []Field  `json:"fields"`
	FollowRequestCount int      `json:"follow_requests_count"`
	Language           string   `json:"language"`
	Note               string   `json:"note"`
	Privacy            string   `json:"privacy"`
	StatusContentType  string   `json:"status_content_type"`
	WebVisibility      string   `json:"web_visibility"`
	Sensitive          bool     `json:"sensitive"`
	AlsoKnownAsURIs    []string `json:"also_known_as_uris"`
}

type Field struct {
	Name       string    `json:"name"`
	Value      string    `json:"value"`
	VerifiedAt time.Time `json:"verified_at"`
}

type AccountRelationship struct {
	Print               bool   `json:"-"`
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

type AccountList struct {
	Label           string
	Accounts        []Account
	BlockedAccounts bool
}
