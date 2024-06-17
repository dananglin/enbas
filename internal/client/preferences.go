// SPDX-FileCopyrightText: 2024 Dan Anglin <d.n.i.anglin@gmail.com>
//
// SPDX-License-Identifier: GPL-3.0-or-later

package client

import (
	"fmt"
	"net/http"

	"codeflow.dananglin.me.uk/apollo/enbas/internal/model"
)

func (g *Client) GetUserPreferences() (*model.Preferences, error) {
	url := g.Authentication.Instance + "/api/v1/preferences"

	var preferences model.Preferences

	if err := g.sendRequest(http.MethodGet, url, nil, &preferences); err != nil {
		return nil, fmt.Errorf("received an error after sending the request to get the user preferences: %w", err)
	}

	return &preferences, nil
}
