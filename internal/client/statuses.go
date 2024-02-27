package client

import (
	"fmt"
	"net/http"

	"codeflow.dananglin.me.uk/apollo/enbas/internal/model"
)

func (g *Client) GetStatus(statusID string) (model.Status, error) {
	path := "/api/v1/statuses/" + statusID
	url := g.Authentication.Instance + path

	var status model.Status

	if err := g.sendRequest(http.MethodGet, url, nil, &status); err != nil {
		return model.Status{}, fmt.Errorf(
			"received an error after sending the request to get the status information; %w",
			err,
		)
	}

	return status, nil
}
