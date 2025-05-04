package model

import "time"

type Notification struct {
	Account   *Account  `json:"account"`
	CreatedAt time.Time `json:"created_at"`
	ID        string    `json:"id"`
	Status    *Status   `json:"status"`
	Type      string    `json:"type"`
}
