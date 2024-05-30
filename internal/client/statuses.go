package client

import (
	"bytes"
	"encoding/json"
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

type CreateStatusForm struct {
	Content     string                  `json:"status"`
	Language    string                  `json:"language"`
	SpoilerText string                  `json:"spoiler_text"`
	Boostable   bool                    `json:"boostable"`
	Federated   bool                    `json:"federated"`
	Likeable    bool                    `json:"likeable"`
	Replyable   bool                    `json:"replyable"`
	Sensitive   bool                    `json:"sensitive"`
	ContentType model.StatusContentType `json:"content_type"`
	Visibility  model.StatusVisibility  `json:"visibility"`
}

func (g *Client) CreateStatus(form CreateStatusForm) (model.Status, error) {
	data, err := json.Marshal(form)
	if err != nil {
		return model.Status{}, fmt.Errorf("unable to create the JSON form; %w", err)
	}

	requestBody := bytes.NewBuffer(data)
	url := g.Authentication.Instance + "/api/v1/statuses"

	var status model.Status

	if err := g.sendRequest(http.MethodPost, url, requestBody, &status); err != nil {
		return model.Status{}, fmt.Errorf(
			"received an error after sending the request to create the status; %w",
			err,
		)
	}

	return status, nil
}
