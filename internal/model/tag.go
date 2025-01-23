package model

type Tag struct {
	Following bool   `json:"following"`
	History   []any  `json:"history"`
	Name      string `json:"name"`
	URL       string `json:"url"`
}

type TagList struct {
	Name string
	Tags []Tag
}
