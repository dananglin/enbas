package model

import "time"

const (
	FilterActionHide string = "hide"
	FilterActionWarn string = "warn"
)

type FilterV2 struct {
	Context   []string        `json:"context"`
	ExpiresAt time.Time       `json:"expires_at"`
	Action    string          `json:"filter_action"`
	ID        string          `json:"id"`
	Keywords  []FilterKeyword `json:"keywords"`
	Statuses  []FilterStatus  `json:"statuses"`
	Title     string          `json:"title"`
}

type FilterKeyword struct {
	ID        string `json:"id"`
	Keyword   string `json:"keyword"`
	WholeWord bool   `json:"whole_word"`
}

type FilterStatus struct {
	ID       string `json:"id"`
	StatusID string `json:"phrase"`
}

type FilterResult struct {
	Filter         FilterV2 `json:"filter"`
	KeywordMatches []string `json:"keyword_matches"`
	StatusMatches  []string `json:"status_matches"`
}
