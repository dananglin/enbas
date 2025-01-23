package gtsclient

import (
	"fmt"
	"net/http"

	"codeflow.dananglin.me.uk/apollo/enbas/internal/model"
)

const preferencesPath = "/api/v1/preferences"

func (g *GTSClient) GetUserPreferences(_ NoRPCArgs, preferences *model.Preferences) error {
	params := requestParameters{
		httpMethod:  http.MethodGet,
		url:         g.authentication.Instance + preferencesPath,
		requestBody: nil,
		contentType: "",
		output:      preferences,
	}

	if err := g.sendRequest(params); err != nil {
		return fmt.Errorf(
			"received an error after sending the request to get the user preferences: %w",
			err,
		)
	}

	return nil
}
