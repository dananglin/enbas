package gtsclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type accessTokenRequest struct {
	RedirectURI  string `json:"redirect_uri"`
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	GrantType    string `json:"grant_type"`
	Code         string `json:"code"`
}

type accessTokenResponse struct {
	AccessToken string `json:"access_token"`
	CreatedAt   int    `json:"created_at"`
	Scope       string `json:"scope"`
	TokenType   string `json:"token_type"`
}

type GetAccessTokenArgs struct {
	ClientID     string
	ClientSecret string
	Code         string
	RedirectURI  string
}

func (g *GTSClient) GetAccessToken(args GetAccessTokenArgs, token *string) error {
	tokenReq := accessTokenRequest{
		RedirectURI:  args.RedirectURI,
		ClientID:     args.ClientID,
		ClientSecret: args.ClientSecret,
		GrantType:    "authorization_code",
		Code:         args.Code,
	}

	data, err := json.Marshal(tokenReq)
	if err != nil {
		return fmt.Errorf(
			"error marshalling the request body: %w",
			err,
		)
	}

	requestBody := bytes.NewBuffer(data)
	url := g.auth.GetInstanceURL() + "/oauth/token"

	var response accessTokenResponse

	params := requestParameters{
		httpMethod:  http.MethodPost,
		url:         url,
		requestBody: requestBody,
		contentType: applicationJSON,
		output:      &response,
	}

	if err := g.sendRequest(params); err != nil {
		return fmt.Errorf(
			"received an error after sending the request to get the access token: %w",
			err,
		)
	}

	if response.AccessToken == "" {
		return EmptyAccessTokenError{}
	}

	*token = response.AccessToken

	return nil
}
