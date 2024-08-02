package client

import (
	"fmt"
	"net/http"

	"codeflow.dananglin.me.uk/apollo/enbas/internal/model"
)

func (g *Client) GetUserPreferences() (*model.Preferences, error) {
	url := g.Authentication.Instance + "/api/v1/preferences"

	var preferences model.Preferences

	params := requestParameters{
		httpMethod:  http.MethodGet,
		url:         url,
		requestBody: nil,
		contentType: "",
		output:      &preferences,
	}

	if err := g.sendRequest(params); err != nil {
		return nil, fmt.Errorf(
			"received an error after sending the request to get the user preferences: %w",
			err,
		)
	}

	return &preferences, nil
}
