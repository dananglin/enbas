package gtsclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"codeflow.dananglin.me.uk/apollo/enbas/internal/info"
	"codeflow.dananglin.me.uk/apollo/enbas/internal/model"
)

func (g *GTSClient) Register(_ NoRPCArgs, _ *NoRPCResults) error {
	form := struct {
		ClientName   string `json:"client_name"`
		RedirectUris string `json:"redirect_uris"`
		Scopes       string `json:"scopes"`
		Website      string `json:"website"`
	}{
		ClientName:   info.ApplicationName,
		RedirectUris: redirectURI,
		Scopes:       "read write",
		Website:      info.ApplicationWebsite,
	}

	data, err := json.Marshal(form)
	if err != nil {
		return fmt.Errorf("error marshalling the request body: %w", err)
	}

	requestBody := bytes.NewBuffer(data)

	var app model.Application

	requestParams := requestParameters{
		httpMethod:  http.MethodPost,
		url:         g.Authentication.Instance + "/api/v1/apps",
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

func (g *GTSClient) AuthCodeURL(_ NoRPCArgs, authCodeURL *string) error {
	escapedRedirectURI := url.QueryEscape(redirectURI)

	*authCodeURL = fmt.Sprintf(
		authCodeURLFormat,
		g.Authentication.Instance,
		g.Authentication.ClientID,
		escapedRedirectURI,
	)

	return nil
}
