// SPDX-FileCopyrightText: 2024 Dan Anglin <d.n.i.anglin@gmail.com>
//
// SPDX-License-Identifier: GPL-3.0-or-later

package model

import (
	"fmt"

	"codeflow.dananglin.me.uk/apollo/enbas/internal/utilities"
)

type Preferences struct {
	PostingDefaultVisibility string `json:"posting:default:visibility"`
	PostingDefaultSensitive  bool   `json:"posting:default:sensitive"`
	PostingDefaultLanguage   string `json:"posting:default:language"`
	ReadingExpandMedia       string `json:"reading:expand:media"`
	ReadingExpandSpoilers    bool   `json:"reading:expand:spoilers"`
	ReadingAutoplayGifs      bool   `json:"reading:autoplay:gifs"`
}

func (p Preferences) Display(noColor bool) string {
	format := `
%s
  %s: %s
  %s: %s
  %s: %t`

	return fmt.Sprintf(
		format,
		utilities.HeaderFormat(noColor, "YOUR PREFERENCES:"),
		utilities.FieldFormat(noColor, "Default post language"), p.PostingDefaultLanguage,
		utilities.FieldFormat(noColor, "Default post visibility"), p.PostingDefaultVisibility,
		utilities.FieldFormat(noColor, "Mark posts as sensitive by default"), p.PostingDefaultSensitive,
	)
}
