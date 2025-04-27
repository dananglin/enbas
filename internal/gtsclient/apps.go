package gtsclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"codeflow.dananglin.me.uk/apollo/enbas/internal/info"
	"codeflow.dananglin.me.uk/apollo/enbas/internal/model"
)

const baseAppsPath = "/api/v1/apps"

func (g *GTSClient) RegisterApp(scopes []string, _ *NoRPCResults) error {
	form := struct {
		ClientName   string `json:"client_name"`
		RedirectUris string `json:"redirect_uris"`
		Scopes       string `json:"scopes"`
		Website      string `json:"website"`
	}{
		ClientName:   info.ApplicationName,
		RedirectUris: redirectURI,
		Scopes:       strings.Join(scopes, " "),
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
		url:         g.authentication.Instance + baseAppsPath,
		requestBody: requestBody,
		contentType: applicationJSON,
		output:      &app,
	}

	if err := g.sendRequest(requestParams); err != nil {
		return fmt.Errorf("received an error after sending the registration request: %w", err)
	}

	g.authentication.ClientID = app.ClientID
	g.authentication.ClientSecret = app.ClientSecret

	return nil
}

func (g *GTSClient) AuthCodeURL(scopes []string, authCodeURL *string) error {
	escapedRedirectURI := url.QueryEscape(redirectURI)

	*authCodeURL = fmt.Sprintf(
		authCodeURLFormat,
		g.authentication.Instance,
		g.authentication.ClientID,
		escapedRedirectURI,
		strings.Join(scopes, "+"),
	)

	return nil
}
