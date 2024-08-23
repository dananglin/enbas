package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"codeflow.dananglin.me.uk/apollo/enbas/internal/info"
	"codeflow.dananglin.me.uk/apollo/enbas/internal/model"
)

type registerRequest struct {
	ClientName   string `json:"client_name"`
	RedirectUris string `json:"redirect_uris"`
	Scopes       string `json:"scopes"`
	Website      string `json:"website"`
}

func (g *Client) Register() error {
	registerParams := registerRequest{
		ClientName:   info.ApplicationName,
		RedirectUris: redirectURI,
		Scopes:       "read write",
		Website:      info.ApplicationWebsite,
	}

	data, err := json.Marshal(registerParams)
	if err != nil {
		return fmt.Errorf("unable to marshal the request body: %w", err)
	}

	requestBody := bytes.NewBuffer(data)
	url := g.Authentication.Instance + "/api/v1/apps"

	var app model.Application

	requestParams := requestParameters{
		httpMethod:  http.MethodPost,
		url:         url,
		requestBody: requestBody,
		contentType: applicationJSON,
		output:      &app,
	}

	if err := g.sendRequest(requestParams); err != nil {
		return fmt.Errorf("received an error after sending the registration request: %w", err)
	}

	g.Authentication.ClientID = app.ClientID
	g.Authentication.ClientSecret = app.ClientSecret

	return nil
}
