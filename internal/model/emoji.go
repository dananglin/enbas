// SPDX-FileCopyrightText: 2024 Dan Anglin <d.n.i.anglin@gmail.com>
//
// SPDX-License-Identifier: GPL-3.0-or-later

package model

type Emoji struct {
	Category        string `json:"category"`
	Shortcode       string `json:"shortcode"`
	StaticURL       string `json:"static_url"`
	URL             string `json:"url"`
	VisibleInPicker bool   `json:"visible_in_picker"`
}
