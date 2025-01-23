package gtsclient

import (
	"fmt"
	"net/http"

	"codeflow.dananglin.me.uk/apollo/enbas/internal/model"
)

const instancePath = "/api/v2/instance"

func (g *GTSClient) GetInstance(_ NoRPCArgs, instance *model.InstanceV2) error {
	params := requestParameters{
		httpMethod:  http.MethodGet,
		url:         g.Authentication.Instance + instancePath,
		requestBody: nil,
		contentType: "",
		output:      instance,
	}

	if err := g.sendRequest(params); err != nil {
		return fmt.Errorf("received an error after sending the request to get the instance details: %w", err)
	}

	return nil
}
