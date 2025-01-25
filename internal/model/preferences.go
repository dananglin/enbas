package model

type Preferences struct {
	Print                    bool   `json:"-"`
	PostingDefaultVisibility string `json:"posting:default:visibility"`
	PostingDefaultSensitive  bool   `json:"posting:default:sensitive"`
	PostingDefaultLanguage   string `json:"posting:default:language"`
	ReadingExpandMedia       string `json:"reading:expand:media"`
	ReadingExpandSpoilers    bool   `json:"reading:expand:spoilers"`
	ReadingAutoplayGifs      bool   `json:"reading:autoplay:gifs"`
}
