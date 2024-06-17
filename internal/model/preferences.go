// SPDX-FileCopyrightText: 2024 Dan Anglin <d.n.i.anglin@gmail.com>
//
// SPDX-License-Identifier: GPL-3.0-or-later

package model

type Preferences struct {
	PostingDefaultVisibility string `json:"posting:default:visibility"`
	PostingDefaultSensitive  bool   `json:"posting:default:sensitive"`
	PostingDefaultLanguage   string `json:"posting:default:language"`
	ReadingExpandMedia       string `json:"reading:expand:media"`
	ReadingExpandSpoilers    bool   `json:"reading:expand:spoilers"`
	ReadingAutoplayGifs      bool   `json:"reading:autoplay:gifs"`
}
