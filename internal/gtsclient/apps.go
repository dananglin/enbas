package gtsclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"codeflow.dananglin.me.uk/apollo/enbas/internal/info"
	"codeflow.dananglin.me.uk/apollo/enbas/internal/model"
)

const baseAppsPath = "/api/v1/apps"

type RegisterAppArgs struct {
	RedirectURI string
	Scopes      []string
}

type RegisteredApp struct {
	ClientID     string
	ClientSecret string
}

func (g *GTSClient) RegisterApp(args RegisterAppArgs, registeredApp *RegisteredApp) error {
	form := struct {
		ClientName   string `json:"client_name"`
		RedirectUris string `json:"redirect_uris"`
		Scopes       string `json:"scopes"`
		Website      string `json:"website"`
	}{
		ClientName:   info.ApplicationName,
		RedirectUris: args.RedirectURI,
		Scopes:       strings.Join(args.Scopes, " "),
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
		url:         g.auth.GetInstanceURL() + baseAppsPath,
		requestBody: requestBody,
		contentType: applicationJSON,
		output:      &app,
	}

	if err := g.sendRequest(requestParams); err != nil {
		return fmt.Errorf("received an error after sending the registration request: %w", err)
	}

	*registeredApp = RegisteredApp{
		ClientID:     app.ClientID,
		ClientSecret: app.ClientSecret,
	}

	return nil
}
