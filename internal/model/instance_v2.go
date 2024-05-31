package model

import (
	"fmt"

	"codeflow.dananglin.me.uk/apollo/enbas/internal/utilities"
)

type InstanceV2 struct {
	AccountDomain   string                  `json:"account_domain"`
	Configuration   InstanceConfiguration   `json:"configuration"`
	Contact         InstanceV2Contact       `json:"contact"`
	Description     string                  `json:"description"`
	DescriptionText string                  `json:"description_text"`
	Domain          string                  `json:"domain"`
	Languages       []string                `json:"languages"`
	Registrations   InstanceV2Registrations `json:"registrations"`
	Rules           []InstanceRule          `json:"rules"`
	SourceURL       string                  `json:"source_url"`
	Terms           string                  `json:"terms"`
	TermsText       string                  `json:"terms_text"`
	Thumbnail       InstanceV2Thumbnail     `json:"thumbnail"`
	Title           string                  `json:"title"`
	Usage           InstanceV2Usage         `json:"usage"`
	Version         string                  `json:"version"`
}

type InstanceConfiguration struct {
	Accounts         InstanceConfigurationAccounts         `json:"accounts"`
	Emojis           InstanceConfigurationEmojis           `json:"emojis"`
	MediaAttachments InstanceConfigurationMediaAttachments `json:"media_attachments"`
	Polls            InstanceConfigurationPolls            `json:"polls"`
	Statuses         InstanceConfigurationStatuses         `json:"statuses"`
	Translation      InstanceV2ConfigurationTranslation    `json:"translation"`
	URLs             InstanceV2URLs                        `json:"urls"`
}

type InstanceConfigurationAccounts struct {
	AllowCustomCSS   bool `json:"allow_custom_css"`
	MaxFeaturedTags  int  `json:"max_featured_tags"`
	MaxProfileFields int  `json:"max_profile_fields"`
}

type InstanceConfigurationEmojis struct {
	EmojiSizeLimit int `json:"emoji_size_limit"`
}

type InstanceConfigurationMediaAttachments struct {
	ImageMatrixLimit    int      `json:"image_matrix_limit"`
	ImageSizeLimit      int      `json:"image_size_limit"`
	SupportedMimeTypes  []string `json:"supported_mime_types"`
	VideoFrameRateLimit int      `json:"video_frame_rate_limit"`
	VideoMatrixLimit    int      `json:"video_matrix_limit"`
	VideoSizeLimit      int      `json:"video_size_limit"`
}

type InstanceConfigurationPolls struct {
	MaxCharactersPerOption int `json:"max_characters_per_option"`
	MaxExpiration          int `json:"max_expiration"`
	MaxOptions             int `json:"max_options"`
	MinExpiration          int `json:"min_expiration"`
}

type InstanceConfigurationStatuses struct {
	CharactersReservedPerURL int      `json:"characters_reserved_per_url"`
	MaxCharacters            int      `json:"max_characters"`
	MaxMediaAttachments      int      `json:"max_media_attachments"`
	SupportedMimeTypes       []string `json:"supported_mime_types"`
}

type InstanceV2ConfigurationTranslation struct {
	Enabled bool `json:"enabled"`
}

type InstanceV2URLs struct {
	Streaming string `json:"streaming"`
}

type InstanceV2Contact struct {
	Account Account `json:"account"`
	Email   string  `json:"email"`
}

type InstanceV2Registrations struct {
	ApprovalRequired bool   `json:"approval_required"`
	Enabled          bool   `json:"enabled"`
	Message          string `json:"message"`
}

type InstanceRule struct {
	ID   string `json:"id"`
	Text string `json:"text"`
}

type InstanceV2Thumbnail struct {
	BlurHash             string                      `json:"blurhash"`
	ThumbnailDescription string                      `json:"thumbnail_description"`
	ThumbnailType        string                      `json:"thumbnail_type"`
	URL                  string                      `json:"url"`
	Versions              InstanceV2ThumbnailVersions `json:"versions"`
}

type InstanceV2ThumbnailVersions struct {
	Size1URL string `json:"@1x"`
	Size2URL string `json:"@2x"`
}

type InstanceV2Usage struct {
	Users InstanceV2Users `json:"users"`
}

type InstanceV2Users struct {
	ActiveMonth int `json:"active_month"`
}

func (i InstanceV2) Display(noColor bool) string {
	format := `
%s
  %s

%s
  %s

%s
  %s

%s
  %s

%s
  Running GoToSocial %s

%s
  %s %s
  %s %s
  %s %s
`

	return fmt.Sprintf(
		format,
		utilities.HeaderFormat(noColor, "INSTANCE TITLE:"),
		i.Title,
		utilities.HeaderFormat(noColor, "INSTANCE DESCRIPTION:"),
		utilities.WrapLines(i.DescriptionText, "\n  ", 80),
		utilities.HeaderFormat(noColor, "DOMAIN:"),
		i.Domain,
		utilities.HeaderFormat(noColor, "TERMS AND CONDITIONS:"),
		utilities.WrapLines(i.TermsText, "\n  ", 80),
		utilities.HeaderFormat(noColor, "VERSION:"),
		i.Version,
		utilities.HeaderFormat(noColor, "CONTACT:"),
		utilities.FieldFormat(noColor, "Name:"),
		utilities.DisplayNameFormat(noColor, i.Contact.Account.DisplayName),
		utilities.FieldFormat(noColor, "Username:"),
		i.Contact.Account.Username,
		utilities.FieldFormat(noColor, "Email:"),
		i.Contact.Email,
	)
}
