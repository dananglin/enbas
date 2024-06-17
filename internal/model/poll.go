// SPDX-FileCopyrightText: 2024 Dan Anglin <d.n.i.anglin@gmail.com>
//
// SPDX-License-Identifier: GPL-3.0-or-later

package model

import (
	"time"
)

type Poll struct {
	Emojis      []Emoji      `json:"emojis"`
	Expired     bool         `json:"expired"`
	Voted       bool         `json:"voted"`
	Multiple    bool         `json:"multiple"`
	ExpiredAt   time.Time    `json:"expires_at"`
	ID          string       `json:"id"`
	OwnVotes    []int        `json:"own_votes"`
	VotersCount int          `json:"voters_count"`
	VotesCount  int          `json:"votes_count"`
	Options     []PollOption `json:"options"`
}

type PollOption struct {
	Title      string `json:"title"`
	VotesCount int    `json:"votes_count"`
}
