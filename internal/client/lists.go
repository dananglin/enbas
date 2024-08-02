package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"codeflow.dananglin.me.uk/apollo/enbas/internal/model"
)

const (
	baseListPath string = "/api/v1/lists"
)

func (g *Client) GetAllLists() ([]model.List, error) {
	url := g.Authentication.Instance + baseListPath

	var lists []model.List

	params := requestParameters{
		httpMethod:  http.MethodGet,
		url:         url,
		requestBody: nil,
		contentType: "",
		output:      &lists,
	}

	if err := g.sendRequest(params); err != nil {
		return nil, fmt.Errorf(
			"received an error after sending the request to get the list of lists: %w",
			err,
		)
	}

	return lists, nil
}

func (g *Client) GetList(listID string) (model.List, error) {
	url := g.Authentication.Instance + baseListPath + "/" + listID

	var list model.List

	params := requestParameters{
		httpMethod:  http.MethodGet,
		url:         url,
		requestBody: nil,
		contentType: "",
		output:      &list,
	}

	if err := g.sendRequest(params); err != nil {
		return model.List{}, fmt.Errorf(
			"received an error after sending the request to get the list: %w",
			err,
		)
	}

	return list, nil
}

type CreateListForm struct {
	Title         string                  `json:"title"`
	RepliesPolicy model.ListRepliesPolicy `json:"replies_policy"`
}

func (g *Client) CreateList(form CreateListForm) (model.List, error) {
	data, err := json.Marshal(form)
	if err != nil {
		return model.List{}, fmt.Errorf("unable to marshal the form: %w", err)
	}

	requestBody := bytes.NewBuffer(data)
	url := g.Authentication.Instance + baseListPath

	var list model.List

	params := requestParameters{
		httpMethod:  http.MethodPost,
		url:         url,
		requestBody: requestBody,
		contentType: applicationJSON,
		output:      &list,
	}

	if err := g.sendRequest(params); err != nil {
		return model.List{}, fmt.Errorf(
			"received an error after sending the request to create the list: %w",
			err,
		)
	}

	return list, nil
}

func (g *Client) UpdateList(listToUpdate model.List) (model.List, error) {
	form := struct {
		Title         string                  `json:"title"`
		RepliesPolicy model.ListRepliesPolicy `json:"replies_policy"`
	}{
		Title:         listToUpdate.Title,
		RepliesPolicy: listToUpdate.RepliesPolicy,
	}

	data, err := json.Marshal(form)
	if err != nil {
		return model.List{}, fmt.Errorf("unable to marshal the form: %w", err)
	}

	requestBody := bytes.NewBuffer(data)
	url := g.Authentication.Instance + baseListPath + "/" + listToUpdate.ID

	var updatedList model.List

	params := requestParameters{
		httpMethod:  http.MethodPut,
		url:         url,
		requestBody: requestBody,
		contentType: applicationJSON,
		output:      &updatedList,
	}

	if err := g.sendRequest(params); err != nil {
		return model.List{}, fmt.Errorf(
			"received an error after sending the request to update the list: %w",
			err,
		)
	}

	return updatedList, nil
}

func (g *Client) DeleteList(listID string) error {
	url := g.Authentication.Instance + baseListPath + "/" + listID

	params := requestParameters{
		httpMethod:  http.MethodDelete,
		url:         url,
		requestBody: nil,
		contentType: "",
		output:      nil,
	}

	if err := g.sendRequest(params); err != nil {
		return fmt.Errorf(
			"received an error after sending the request to delete the list: %w",
			err,
		)
	}

	return nil
}

func (g *Client) AddAccountsToList(listID string, accountIDs []string) error {
	form := struct {
		AccountIDs []string `json:"account_ids"`
	}{
		AccountIDs: accountIDs,
	}

	data, err := json.Marshal(form)
	if err != nil {
		return fmt.Errorf("unable to marshal the form: %w", err)
	}

	requestBody := bytes.NewBuffer(data)
	url := g.Authentication.Instance + baseListPath + "/" + listID + "/accounts"

	params := requestParameters{
		httpMethod:  http.MethodPost,
		url:         url,
		requestBody: requestBody,
		contentType: applicationJSON,
		output:      nil,
	}

	if err := g.sendRequest(params); err != nil {
		return fmt.Errorf(
			"received an error after sending the request to add the accounts to the list: %w",
			err,
		)
	}

	return nil
}

func (g *Client) RemoveAccountsFromList(listID string, accountIDs []string) error {
	form := struct {
		AccountIDs []string `json:"account_ids"`
	}{
		AccountIDs: accountIDs,
	}

	data, err := json.Marshal(form)
	if err != nil {
		return fmt.Errorf("unable to marshal the form: %w", err)
	}

	requestBody := bytes.NewBuffer(data)
	url := g.Authentication.Instance + baseListPath + "/" + listID + "/accounts"

	params := requestParameters{
		httpMethod:  http.MethodDelete,
		url:         url,
		requestBody: requestBody,
		contentType: applicationJSON,
		output:      nil,
	}

	if err := g.sendRequest(params); err != nil {
		return fmt.Errorf(
			"received an error after sending the request to remove the accounts from the list: %w",
			err,
		)
	}

	return nil
}

func (g *Client) GetAccountsFromList(listID string, limit int) ([]model.Account, error) {
	path := fmt.Sprintf("%s/%s/accounts?limit=%d", baseListPath, listID, limit)
	url := g.Authentication.Instance + path

	var accounts []model.Account

	params := requestParameters{
		httpMethod:  http.MethodGet,
		url:         url,
		requestBody: nil,
		contentType: "",
		output:      &accounts,
	}

	if err := g.sendRequest(params); err != nil {
		return nil, fmt.Errorf(
			"received an error after sending the request to get the accounts from the list: %w",
			err,
		)
	}

	return accounts, nil
}
