// SPDX-FileCopyrightText: 2024 Dan Anglin <d.n.i.anglin@gmail.com>
//
// SPDX-License-Identifier: GPL-3.0-or-later

package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"codeflow.dananglin.me.uk/apollo/enbas/internal"
)

var errEmptyAccessToken = errors.New("received an empty access token")

type tokenRequest struct {
	RedirectUri  string `json:"redirect_uri"`
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	GrantType    string `json:"grant_type"`
	Code         string `json:"code"`
}

type tokenResponse struct {
	AccessToken string `json:"access_token"`
	CreatedAt   int    `json:"created_at"`
	Scope       string `json:"scope"`
	TokenType   string `json:"token_type"`
}

func (g *Client) UpdateToken(code string) error {
	params := tokenRequest{
		RedirectUri:  internal.RedirectUri,
		ClientID:     g.Authentication.ClientID,
		ClientSecret: g.Authentication.ClientSecret,
		GrantType:    "authorization_code",
		Code:         code,
	}

	data, err := json.Marshal(params)
	if err != nil {
		return fmt.Errorf("unable to marshal the request body: %w", err)
	}

	requestBody := bytes.NewBuffer(data)
	url := g.Authentication.Instance + "/oauth/token"

	var response tokenResponse

	if err := g.sendRequest(http.MethodPost, url, requestBody, &response); err != nil {
		return fmt.Errorf("received an error after sending the token request: %w", err)
	}

	if response.AccessToken == "" {
		return errEmptyAccessToken
	}

	g.Authentication.AccessToken = response.AccessToken

	return nil
}
