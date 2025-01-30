package model

import "time"

type InteractionRequest struct {
	Account    Account   `json:"account"`
	Reply      Status    `json:"reply"`
	Status     Status    `json:"status"`
	AcceptedAt time.Time `json:"accepted_at"`
	CreatedAt  time.Time `json:"created_at"`
	RejectedAt time.Time `json:"rejected_at"`
	ID         string    `json:"id"`
	Type       string    `json:"type"`
	URI        string    `json:"uri"`
}
