package model

import "time"

type Application struct {
	ClientID     string    `json:"client_id"`
	ClientSecret string    `json:"client_secret"`
	CreatedAt    time.Time `json:"created_at"`
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	RedirectURI  string    `json:"redirect_uri"`
	RedirectURIs []string  `json:"redirect_uris"`
	Scopes       []string  `json:"scopes"`
	VapidKey     string    `json:"vapid_key"`
	Website      string    `json:"website"`
}
