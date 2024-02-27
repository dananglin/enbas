package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"codeflow.dananglin.me.uk/apollo/enbas/internal/model"
)

const (
	listPath string = "/api/v1/lists"
)

func (g *Client) GetAllLists() ([]model.List, error) {
	url := g.Authentication.Instance + listPath

	var lists []model.List

	if err := g.sendRequest(http.MethodGet, url, nil, &lists); err != nil {
		return nil, fmt.Errorf(
			"received an error after sending the request to get the list of lists; %w",
			err,
		)
	}

	return lists, nil
}

func (g *Client) GetList(listID string) (model.List, error) {
	url := g.Authentication.Instance + listPath + "/" + listID

	var list model.List

	if err := g.sendRequest(http.MethodGet, url, nil, &list); err != nil {
		return model.List{}, fmt.Errorf(
			"received an error after sending the request to get the list; %w",
			err,
		)
	}

	return list, nil
}

func (g *Client) CreateList(title string, repliesPolicy model.ListRepliesPolicy) (model.List, error) {
	params := struct {
		Title         string                  `json:"title"`
		RepliesPolicy model.ListRepliesPolicy `json:"replies_policy"`
	}{
		Title:         title,
		RepliesPolicy: repliesPolicy,
	}

	data, err := json.Marshal(params)
	if err != nil {
		return model.List{}, fmt.Errorf("unable to marshal the request body; %w", err)
	}

	requestBody := bytes.NewBuffer(data)
	url := g.Authentication.Instance + "/api/v1/lists"

	var list model.List

	if err := g.sendRequest(http.MethodPost, url, requestBody, &list); err != nil {
		return model.List{}, fmt.Errorf(
			"received an error after sending the request to create the list; %w",
			err,
		)
	}

	return list, nil
}

func (g *Client) UpdateList(listToUpdate model.List) (model.List, error) {
	params := struct {
		Title         string                  `json:"title"`
		RepliesPolicy model.ListRepliesPolicy `json:"replies_policy"`
	}{
		Title:         listToUpdate.Title,
		RepliesPolicy: listToUpdate.RepliesPolicy,
	}

	data, err := json.Marshal(params)
	if err != nil {
		return model.List{}, fmt.Errorf("unable to marshal the request body; %w", err)
	}

	requestBody := bytes.NewBuffer(data)
	url := g.Authentication.Instance + listPath + "/" + listToUpdate.ID

	var updatedList model.List

	if err := g.sendRequest(http.MethodPut, url, requestBody, &updatedList); err != nil {
		return model.List{}, fmt.Errorf(
			"received an error after sending the request to update the list; %w",
			err,
		)
	}

	return updatedList, nil
}

func (g *Client) DeleteList(listID string) error {
	url := g.Authentication.Instance + "/api/v1/lists/" + listID

	return g.sendRequest(http.MethodDelete, url, nil, nil)
}
