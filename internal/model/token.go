package model

import "time"

type Token struct {
	Application Application `json:"application"`
	CreatedAt   time.Time   `json:"created_at"`
	LastUsed    time.Time   `json:"last_used"`
	ID          string      `json:"id"`
	Scope       string      `json:"scope"`
}

type TokenList struct {
	Label  string
	Tokens []Token
}
