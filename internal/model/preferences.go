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

func (p Preferences) String() string {
	format := `
%s
  %s: %s
  %s: %s
  %s: %t`

	return fmt.Sprintf(
		format,
		utilities.HeaderFormat("YOUR PREFERENCES:"),
		utilities.FieldFormat("Default post language"), p.PostingDefaultLanguage,
		utilities.FieldFormat("Default post visibility"), p.PostingDefaultVisibility,
		utilities.FieldFormat("Mark posts as sensitive by default"), p.PostingDefaultSensitive,
	)
}