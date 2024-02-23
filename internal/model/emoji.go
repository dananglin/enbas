package model

type Emoji struct {
	Category        string `json:"category"`
	Shortcode       string `json:"shortcode"`
	StaticURL       string `json:"static_url"`
	URL             string `json:"url"`
	VisibleInPicker bool   `json:"visible_in_picker"`
}
