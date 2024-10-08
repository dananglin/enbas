package client

import (
	"fmt"
	"net/http"

	"codeflow.dananglin.me.uk/apollo/enbas/internal/model"
)

func (g *Client) GetInstance() (model.InstanceV2, error) {
	path := "/api/v2/instance"
	url := g.Authentication.Instance + path

	var instance model.InstanceV2

	params := requestParameters{
		httpMethod:  http.MethodGet,
		url:         url,
		requestBody: nil,
		contentType: "",
		output:      &instance,
	}

	if err := g.sendRequest(params); err != nil {
		return model.InstanceV2{}, fmt.Errorf("received an error after sending the request to get the instance details: %w", err)
	}

	return instance, nil
}
