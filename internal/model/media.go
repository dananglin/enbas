package model

type MediaAttachment struct {
	Meta             MediaMeta `json:"meta"`
	Blurhash         string    `json:"blurhash"`
	Description      string    `json:"description"`
	ID               string    `json:"id"`
	PreviewRemoteURL string    `json:"preview_remote_url"`
	PreviewURL       string    `json:"preview_url"`
	RemoteURL        string    `json:"remote_url"`
	TextURL          string    `json:"text_url"`
	Type             string    `json:"type"`
	URL              string    `json:"url"`
}

type MediaMeta struct {
	Focus    MediaFocus      `json:"focus"`
	Original MediaDimensions `json:"original"`
	Small    MediaDimensions `json:"small"`
}

type MediaFocus struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

type MediaDimensions struct {
	Aspect    float64 `json:"aspect"`
	Bitrate   int     `json:"bitrate"`
	Duration  float64 `json:"duration"`
	FrameRate string  `json:"frame_rate"`
	Size      string  `json:"size"`
	Height    int     `json:"height"`
	Width     int     `json:"width"`
}
