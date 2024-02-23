package model

import "time"

type Account struct {
	Acct           string      `json:"acct"`
	Avatar         string      `json:"avatar"`
	AvatarStatic   string      `json:"avatar_static"`
	Bot            bool        `json:"bot"`
	CreatedAt      time.Time   `json:"created_at"`
	CustomCSS      string      `json:"custom_css"`
	Discoverable   bool        `json:"discoverable"`
	DisplayName    string      `json:"display_name"`
	Emojis         []Emoji     `json:"emojis"`
	EnableRSS      bool        `json:"enable_rss"`
	Fields         []Field     `json:"fields"`
	FollowersCount int         `json:"followers_count"`
	FollowingCount int         `json:"following_count"`
	Header         string      `json:"header"`
	HeaderStatic   string      `json:"header_static"`
	ID             string      `json:"id"`
	LastStatusAt   string      `json:"last_status_at"`
	Locked         bool        `json:"locked"`
	MuteExpiresAt  time.Time   `json:"mute_expires_at"`
	Note           string      `json:"note"`
	Role           AccountRole `json:"role"`
	Source         Source      `json:"source"`
	StatusCount    int         `json:"statuses_count"`
	Suspended      bool        `json:"suspended"`
	URL            string      `json:"url"`
	Username       string      `json:"username"`
}

type AccountRole struct {
	Name string `json:"name"`
}

type Source struct {
	Fields             []Field `json:"fields"`
	FollowRequestCount int     `json:"follow_requests_count"`
	Language           string  `json:"language"`
	Note               string  `json:"note"`
	Privacy            string  `json:"string"`
	Sensitive          bool    `json:"sensitive"`
	StatusContentType  string  `json:"status_content_type"`
}

type Field struct {
	Name       string `json:"name"`
	Value      string `json:"value"`
	VerifiedAt string `json:"verified_at"`
}
