package gtsclient

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

func (g *GTSClient) GetAllLists(_ NoRPCArgs, lists *[]model.List) error {
	params := requestParameters{
		httpMethod:  http.MethodGet,
		url:         g.authentication.Instance + baseListPath,
		requestBody: nil,
		contentType: "",
		output:      lists,
	}

	if err := g.sendRequest(params); err != nil {
		return fmt.Errorf(
			"received an error after sending the request to get the list of lists: %w",
			err,
		)
	}

	return nil
}

func (g *GTSClient) GetList(listID string, list *model.List) error {
	params := requestParameters{
		httpMethod:  http.MethodGet,
		url:         g.authentication.Instance + baseListPath + "/" + listID,
		requestBody: nil,
		contentType: "",
		output:      list,
	}

	if err := g.sendRequest(params); err != nil {
		return fmt.Errorf(
			"received an error after sending the request to get the list: %w",
			err,
		)
	}

	return nil
}

type CreateListArgs struct {
	Title         string
	RepliesPolicy model.ListRepliesPolicy
	Exclusive     bool
}

func (g *GTSClient) CreateList(args CreateListArgs, list *model.List) error {
	form := struct {
		Title         string                  `json:"title"`
		RepliesPolicy model.ListRepliesPolicy `json:"replies_policy"`
		Exclusive     bool                    `json:"exclusive"`
	}{
		Title:         args.Title,
		RepliesPolicy: args.RepliesPolicy,
		Exclusive:     args.Exclusive,
	}

	data, err := json.Marshal(form)
	if err != nil {
		return fmt.Errorf("unable to marshal the form: %w", err)
	}

	requestBody := bytes.NewBuffer(data)

	params := requestParameters{
		httpMethod:  http.MethodPost,
		url:         g.authentication.Instance + baseListPath,
		requestBody: requestBody,
		contentType: applicationJSON,
		output:      list,
	}

	if err := g.sendRequest(params); err != nil {
		return fmt.Errorf(
			"received an error after sending the request to create the list: %w",
			err,
		)
	}

	return nil
}

func (g *GTSClient) UpdateList(listToUpdate model.List, updatedList *model.List) error {
	form := struct {
		Title         string                  `json:"title"`
		RepliesPolicy model.ListRepliesPolicy `json:"replies_policy"`
	}{
		Title:         listToUpdate.Title,
		RepliesPolicy: listToUpdate.RepliesPolicy,
	}

	data, err := json.Marshal(form)
	if err != nil {
		return fmt.Errorf("unable to marshal the form: %w", err)
	}

	requestBody := bytes.NewBuffer(data)

	params := requestParameters{
		httpMethod:  http.MethodPut,
		url:         g.authentication.Instance + baseListPath + "/" + listToUpdate.ID,
		requestBody: requestBody,
		contentType: applicationJSON,
		output:      updatedList,
	}

	if err := g.sendRequest(params); err != nil {
		return fmt.Errorf(
			"received an error after sending the request to update the list: %w",
			err,
		)
	}

	return nil
}

func (g *GTSClient) DeleteList(listID string, _ *NoRPCResults) error {
	params := requestParameters{
		httpMethod:  http.MethodDelete,
		url:         g.authentication.Instance + baseListPath + "/" + listID,
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

type AddAccountsToListArgs struct {
	ListID     string
	AccountIDs []string
}

func (g *GTSClient) AddAccountsToList(args AddAccountsToListArgs, _ *NoRPCResults) error {
	form := struct {
		AccountIDs []string `json:"account_ids"`
	}{
		AccountIDs: args.AccountIDs,
	}

	data, err := json.Marshal(form)
	if err != nil {
		return fmt.Errorf("unable to marshal the form: %w", err)
	}

	requestBody := bytes.NewBuffer(data)

	params := requestParameters{
		httpMethod:  http.MethodPost,
		url:         g.authentication.Instance + baseListPath + "/" + args.ListID + "/accounts",
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

type RemoveAccountsFromListArgs struct {
	ListID     string
	AccountIDs []string
}

func (g *GTSClient) RemoveAccountsFromList(args RemoveAccountsFromListArgs, _ *NoRPCResults) error {
	form := struct {
		AccountIDs []string `json:"account_ids"`
	}{
		AccountIDs: args.AccountIDs,
	}

	data, err := json.Marshal(form)
	if err != nil {
		return fmt.Errorf("unable to marshal the form: %w", err)
	}

	requestBody := bytes.NewBuffer(data)

	params := requestParameters{
		httpMethod:  http.MethodDelete,
		url:         g.authentication.Instance + baseListPath + "/" + args.ListID + "/accounts",
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

type GetAccountsFromListArgs struct {
	ListID string
	Limit  int
}

func (g *GTSClient) GetAccountsFromList(args GetAccountsFromListArgs, list *model.AccountList) error {
	path := fmt.Sprintf("%s/%s/accounts?limit=%d", baseListPath, args.ListID, args.Limit)

	var accounts []model.Account

	params := requestParameters{
		httpMethod:  http.MethodGet,
		url:         g.authentication.Instance + path,
		requestBody: nil,
		contentType: "",
		output:      &accounts,
	}

	if err := g.sendRequest(params); err != nil {
		return fmt.Errorf(
			"received an error after sending the request to get the accounts from the list: %w",
			err,
		)
	}

	*list = model.AccountList{
		Label:           "Accounts",
		Accounts:        accounts,
		BlockedAccounts: false,
	}

	return nil
}
