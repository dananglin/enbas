// SPDX-FileCopyrightText: 2024 Dan Anglin <d.n.i.anglin@gmail.com>
//
// SPDX-License-Identifier: GPL-3.0-or-later

package model

type Application struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	ID           string `json:"id"`
	Name         string `json:"name"`
	RedirectURI  string `json:"redirect_uri"`
	VapidKey     string `json:"vapid_key"`
	Website      string `json:"website"`
}
