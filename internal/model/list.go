package model

type List struct {
	ID            string            `json:"id"`
	RepliesPolicy string            `json:"replies_policy"`
	Title         string            `json:"title"`
	Exclusive     bool              `json:"exclusive"`
	Accounts      map[string]string `json:"-"`
}
