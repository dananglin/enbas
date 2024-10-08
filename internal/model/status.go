package model

import (
	"time"
)

type Status struct {
	Account            Account          `json:"account"`
	Application        Application      `json:"application"`
	Bookmarked         bool             `json:"bookmarked"`
	Card               Card             `json:"card"`
	Content            string           `json:"content"`
	CreatedAt          time.Time        `json:"created_at"`
	Emojis             []Emoji          `json:"emojis"`
	Favourited         bool             `json:"favourited"`
	FavouritesCount    int              `json:"favourites_count"`
	ID                 string           `json:"id"`
	InReplyToAccountID string           `json:"in_reply_to_account_id"`
	InReplyToID        string           `json:"in_reply_to_id"`
	Language           string           `json:"language"`
	MediaAttachments   []Attachment     `json:"media_attachments"`
	Mentions           []Mention        `json:"mentions"`
	Muted              bool             `json:"muted"`
	Pinned             bool             `json:"pinned"`
	Poll               *Poll            `json:"poll"`
	Reblog             *StatusReblogged `json:"reblog"`
	Reblogged          bool             `json:"reblogged"`
	ReblogsCount       int              `json:"reblogs_count"`
	RepliesCount       int              `json:"replies_count"`
	Sensitive          bool             `json:"sensitive"`
	SpoilerText        string           `json:"spoiler_text"`
	Tags               []Tag            `json:"tags"`
	Text               string           `json:"text"`
	URI                string           `json:"uri"`
	URL                string           `json:"url"`
	Visibility         StatusVisibility `json:"visibility"`
}

type Card struct {
	AuthorName   string `json:"author_name"`
	AuthorURL    string `json:"author_url"`
	Blurhash     string `json:"blurhash"`
	Description  string `json:"description"`
	EmbedURL     string `json:"embed_url"`
	HTML         string `json:"html"`
	Image        string `json:"image"`
	ProviderName string `json:"provider_name"`
	ProviderURL  string `json:"provider_url"`
	Title        string `json:"title"`
	Type         string `json:"type"`
	URL          string `json:"url"`
	Height       int    `json:"height"`
	Width        int    `json:"width"`
}

type Mention struct {
	Acct     string `json:"acct"`
	ID       string `json:"id"`
	URL      string `json:"url"`
	Username string `json:"username"`
}

type StatusReblogged struct {
	Account            Account          `json:"account"`
	Application        Application      `json:"application"`
	Bookmarked         bool             `json:"bookmarked"`
	Card               Card             `json:"card"`
	Content            string           `json:"content"`
	CreatedAt          time.Time        `json:"created_at"`
	Emojis             []Emoji          `json:"emojis"`
	Favourited         bool             `json:"favourited"`
	FavouritesCount    int              `json:"favourites_count"`
	ID                 string           `json:"id"`
	InReplyToAccountID string           `json:"in_reply_to_account_id"`
	InReplyToID        string           `json:"in_reply_to_id"`
	Language           string           `json:"language"`
	MediaAttachments   []Attachment     `json:"media_attachments"`
	Mentions           []Mention        `json:"mentions"`
	Muted              bool             `json:"muted"`
	Pinned             bool             `json:"pinned"`
	Poll               *Poll            `json:"poll"`
	Reblogged          bool             `json:"reblogged"`
	RebloggsCount      int              `json:"reblogs_count"`
	RepliesCount       int              `json:"replies_count"`
	Sensitive          bool             `json:"sensitive"`
	SpoilerText        string           `json:"spoiler_text"`
	Tags               []Tag            `json:"tags"`
	Text               string           `json:"text"`
	URI                string           `json:"uri"`
	URL                string           `json:"url"`
	Visibility         StatusVisibility `json:"visibility"`
}

type Tag struct {
	History []any  `json:"history"`
	Name    string `json:"name"`
	URL     string `json:"url"`
}

type Attachment struct {
	Meta             MediaMeta `json:"meta"`
	Blurhash         string    `json:"blurhash"`
	Description      string    `json:"description"`
	ID               string    `json:"id"`
	PreviewRemoteURL string    `json:"preview_remote_url"`
	PreviewURL       string    `json:"preview_url"`
	RemoteURL        string    `json:"remote_url"`
	TextURL          string    `json:"text_url"`
	Type             string    `json:"type"`
	URL              string    `json:"url"`
}

type MediaMeta struct {
	Focus    MediaFocus      `json:"focus"`
	Original MediaDimensions `json:"original"`
	Small    MediaDimensions `json:"small"`
}

type MediaFocus struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

type MediaDimensions struct {
	Aspect    float64 `json:"aspect"`
	Bitrate   int     `json:"bitrate"`
	Duration  float64 `json:"duration"`
	FrameRate string  `json:"frame_rate"`
	Size      string  `json:"size"`
	Height    int     `json:"height"`
	Width     int     `json:"width"`
}

type StatusList struct {
	Name     string
	Statuses []Status
}
