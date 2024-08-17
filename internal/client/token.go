package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"codeflow.dananglin.me.uk/apollo/enbas/internal"
)

type tokenRequest struct {
	RedirectURI  string `json:"redirect_uri"`
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
	tokenReq := tokenRequest{
		RedirectURI:  internal.RedirectURI,
		ClientID:     g.Authentication.ClientID,
		ClientSecret: g.Authentication.ClientSecret,
		GrantType:    "authorization_code",
		Code:         code,
	}

	data, err := json.Marshal(tokenReq)
	if err != nil {
		return fmt.Errorf("unable to marshal the request body: %w", err)
	}

	requestBody := bytes.NewBuffer(data)
	url := g.Authentication.Instance + "/oauth/token"

	var response tokenResponse

	params := requestParameters{
		httpMethod:  http.MethodPost,
		url:         url,
		requestBody: requestBody,
		contentType: applicationJSON,
		output:      &response,
	}

	if err := g.sendRequest(params); err != nil {
		return fmt.Errorf("received an error after sending the token request: %w", err)
	}

	if response.AccessToken == "" {
		return Error{"received an empty access token"}
	}

	g.Authentication.AccessToken = response.AccessToken

	return nil
}
