package model

import (
	"time"
)

type Status struct {
	Account            Account           `json:"account"`
	Application        Application       `json:"application"`
	Bookmarked         bool              `json:"bookmarked"`
	Card               Card              `json:"card"`
	Content            string            `json:"content"`
	CreatedAt          time.Time         `json:"created_at"`
	Emojis             []Emoji           `json:"emojis"`
	Favourited         bool              `json:"favourited"`
	FavouritesCount    int               `json:"favourites_count"`
	ID                 string            `json:"id"`
	InReplyToAccountID string            `json:"in_reply_to_account_id"`
	InReplyToID        string            `json:"in_reply_to_id"`
	InteractionPolicy  InteractionPolicy `json:"interaction_policy"`
	Language           string            `json:"language"`
	LocalOnly          bool              `json:"local_only"`
	MediaAttachments   []MediaAttachment `json:"media_attachments"`
	Mentions           []Mention         `json:"mentions"`
	Muted              bool              `json:"muted"`
	Pinned             bool              `json:"pinned"`
	Poll               *Poll             `json:"poll"`
	Reblog             *StatusReblogged  `json:"reblog"`
	Reblogged          bool              `json:"reblogged"`
	ReblogsCount       int               `json:"reblogs_count"`
	RepliesCount       int               `json:"replies_count"`
	Sensitive          bool              `json:"sensitive"`
	SpoilerText        string            `json:"spoiler_text"`
	Tags               []Tag             `json:"tags"`
	Text               string            `json:"text"`
	URI                string            `json:"uri"`
	URL                string            `json:"url"`
	Visibility         StatusVisibility  `json:"visibility"`
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
	Account            Account           `json:"account"`
	Application        Application       `json:"application"`
	Bookmarked         bool              `json:"bookmarked"`
	Card               Card              `json:"card"`
	Content            string            `json:"content"`
	CreatedAt          time.Time         `json:"created_at"`
	Emojis             []Emoji           `json:"emojis"`
	Favourited         bool              `json:"favourited"`
	FavouritesCount    int               `json:"favourites_count"`
	ID                 string            `json:"id"`
	InReplyToAccountID string            `json:"in_reply_to_account_id"`
	InReplyToID        string            `json:"in_reply_to_id"`
	Language           string            `json:"language"`
	MediaAttachments   []MediaAttachment `json:"media_attachments"`
	Mentions           []Mention         `json:"mentions"`
	Muted              bool              `json:"muted"`
	Pinned             bool              `json:"pinned"`
	Poll               *Poll             `json:"poll"`
	Reblogged          bool              `json:"reblogged"`
	ReblogsCount       int               `json:"reblogs_count"`
	RepliesCount       int               `json:"replies_count"`
	Sensitive          bool              `json:"sensitive"`
	SpoilerText        string            `json:"spoiler_text"`
	Tags               []Tag             `json:"tags"`
	Text               string            `json:"text"`
	URI                string            `json:"uri"`
	URL                string            `json:"url"`
	Visibility         StatusVisibility  `json:"visibility"`
}

type StatusList struct {
	Name     string
	Statuses []Status
}
