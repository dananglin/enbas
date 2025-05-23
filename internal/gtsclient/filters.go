package gtsclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"codeflow.dananglin.me.uk/apollo/enbas/internal/model"
)

const (
	baseFiltersV2Path      string = "/api/v2/filters"
	baseFilterKeywordsPath string = baseFiltersV2Path + "/keywords"
	baseFilterStatusesPath string = baseFiltersV2Path + "/statuses"
)

func (g *GTSClient) GetAllFilters(_ NoRPCArgs, filters *[]model.FilterV2) error {
	params := requestParameters{
		httpMethod:  http.MethodGet,
		url:         g.authentication.Instance + baseFiltersV2Path,
		requestBody: nil,
		contentType: "",
		output:      filters,
	}

	if err := g.sendRequest(params); err != nil {
		return fmt.Errorf(
			"received an error after sending the request to get the list of filters: %w",
			err,
		)
	}

	return nil
}

func (g *GTSClient) GetFilter(filterID string, filter *model.FilterV2) error {
	params := requestParameters{
		httpMethod:  http.MethodGet,
		url:         g.authentication.Instance + baseFiltersV2Path + "/" + filterID,
		requestBody: nil,
		contentType: "",
		output:      filter,
	}

	if err := g.sendRequest(params); err != nil {
		return fmt.Errorf(
			"received an error after sending the request to get the filter: %w",
			err,
		)
	}

	return nil
}

type CreateFilterArgs struct {
	Title        string
	FilterAction string
	Context      []string
	ExpiresIn    time.Duration
}

func (g *GTSClient) CreateFilter(args CreateFilterArgs, filter *model.FilterV2) error {
	form := struct {
		Title        string   `json:"title"`
		FilterAction string   `json:"filter_action"`
		Context      []string `json:"context"`
		ExpiresIn    int      `json:"expires_in"`
	}{
		Title:        args.Title,
		FilterAction: args.FilterAction,
		Context:      args.Context,
		ExpiresIn:    int(args.ExpiresIn.Seconds()),
	}

	data, err := json.Marshal(form)
	if err != nil {
		return fmt.Errorf("error marshalling the form: %w", err)
	}

	requestBody := bytes.NewBuffer(data)

	params := requestParameters{
		httpMethod:  http.MethodPost,
		url:         g.authentication.Instance + baseFiltersV2Path,
		requestBody: requestBody,
		contentType: applicationJSON,
		output:      filter,
	}

	if err := g.sendRequest(params); err != nil {
		return fmt.Errorf(
			"received an error after sending the request to create the filter: %w",
			err,
		)
	}

	return nil
}

type EditFilterArgs struct {
	FilterID     string
	Title        string
	FilterAction string
	Context      []string
	ExpiresIn    time.Duration
}

func (g *GTSClient) EditFilter(args EditFilterArgs, _ *NoRPCResults) error {
	form := struct {
		Title        string   `json:"title"`
		FilterAction string   `json:"filter_action"`
		Context      []string `json:"context"`
		ExpiresIn    int      `json:"expires_in"`
	}{
		Title:        args.Title,
		FilterAction: args.FilterAction,
		Context:      args.Context,
		ExpiresIn:    int(args.ExpiresIn.Seconds()),
	}

	data, err := json.Marshal(form)
	if err != nil {
		return fmt.Errorf("error marshalling the form: %w", err)
	}

	requestBody := bytes.NewBuffer(data)

	params := requestParameters{
		httpMethod:  http.MethodPut,
		url:         g.authentication.Instance + baseFiltersV2Path + "/" + args.FilterID,
		requestBody: requestBody,
		contentType: applicationJSON,
		output:      nil,
	}

	if err := g.sendRequest(params); err != nil {
		return fmt.Errorf(
			"received an error after sending the request to update the filter: %w",
			err,
		)
	}

	return nil
}

func (g *GTSClient) DeleteFilter(filterID string, _ *NoRPCResults) error {
	params := requestParameters{
		httpMethod:  http.MethodDelete,
		url:         g.authentication.Instance + baseFiltersV2Path + "/" + filterID,
		requestBody: nil,
		contentType: "",
		output:      nil,
	}

	if err := g.sendRequest(params); err != nil {
		return fmt.Errorf(
			"received an error after sending the request to delete the filter: %w",
			err,
		)
	}

	return nil
}

type AddFilterKeywordToFilterArgs struct {
	FilterID  string
	Keyword   string
	WholeWord bool
}

func (g *GTSClient) AddFilterKeywordToFilter(args AddFilterKeywordToFilterArgs, filterKeyword *model.FilterKeyword) error {
	form := struct {
		Keyword   string `json:"keyword"`
		WholeWord bool   `json:"whole_word"`
	}{
		Keyword:   args.Keyword,
		WholeWord: args.WholeWord,
	}

	data, err := json.Marshal(form)
	if err != nil {
		return fmt.Errorf("error marshalling the form: %w", err)
	}

	requestBody := bytes.NewBuffer(data)

	params := requestParameters{
		httpMethod:  http.MethodPost,
		url:         g.authentication.Instance + baseFiltersV2Path + "/" + args.FilterID + "/keywords",
		requestBody: requestBody,
		contentType: applicationJSON,
		output:      filterKeyword,
	}

	if err := g.sendRequest(params); err != nil {
		return fmt.Errorf(
			"received an error after sending the request to add the filter keyword to the filter: %w",
			err,
		)
	}

	return nil
}

func (g *GTSClient) DeleteFilterKeyword(filterKeywordID string, _ *NoRPCResults) error {
	params := requestParameters{
		httpMethod:  http.MethodDelete,
		url:         g.authentication.Instance + baseFilterKeywordsPath + "/" + filterKeywordID,
		requestBody: nil,
		contentType: "",
		output:      nil,
	}

	if err := g.sendRequest(params); err != nil {
		return fmt.Errorf(
			"received an error after sending the request to delete the filter keyword: %w",
			err,
		)
	}

	return nil
}

func (g *GTSClient) GetFilterKeyword(filterKeywordID string, filterKeyword *model.FilterKeyword) error {
	params := requestParameters{
		httpMethod:  http.MethodGet,
		url:         g.authentication.Instance + baseFilterKeywordsPath + "/" + filterKeywordID,
		requestBody: nil,
		contentType: "",
		output:      filterKeyword,
	}

	if err := g.sendRequest(params); err != nil {
		return fmt.Errorf(
			"received an error after sending the request to get the filter keyword: %w",
			err,
		)
	}

	return nil
}

type UpdateFilterKeywordArgs struct {
	FilterKeywordID string
	Keyword         string
	WholeWord       bool
}

func (g *GTSClient) UpdateFilterKeyword(args UpdateFilterKeywordArgs, filterKeyword *model.FilterKeyword) error {
	form := struct {
		Keyword   string `json:"keyword"`
		WholeWord bool   `json:"whole_word"`
	}{
		Keyword:   args.Keyword,
		WholeWord: args.WholeWord,
	}

	data, err := json.Marshal(form)
	if err != nil {
		return fmt.Errorf("error marshalling the form: %w", err)
	}

	requestBody := bytes.NewBuffer(data)

	params := requestParameters{
		httpMethod:  http.MethodPut,
		url:         g.authentication.Instance + baseFilterKeywordsPath + "/" + args.FilterKeywordID,
		requestBody: requestBody,
		contentType: applicationJSON,
		output:      filterKeyword,
	}

	if err := g.sendRequest(params); err != nil {
		return fmt.Errorf(
			"received an error after sending the request to update the filter keyword: %w",
			err,
		)
	}

	return nil
}

type AddFilterStatusToFilterArgs struct {
	FilterID string
	StatusID string
}

func (g *GTSClient) AddFilterStatusToFilter(args AddFilterStatusToFilterArgs, filterStatus *model.FilterStatus) error {
	form := struct {
		StatusID string `json:"status_id"`
	}{
		StatusID: args.StatusID,
	}

	data, err := json.Marshal(form)
	if err != nil {
		return fmt.Errorf("error marshalling the form: %w", err)
	}

	requestBody := bytes.NewBuffer(data)

	params := requestParameters{
		httpMethod:  http.MethodPost,
		url:         g.authentication.Instance + baseFiltersV2Path + "/" + args.FilterID + "/statuses",
		requestBody: requestBody,
		contentType: applicationJSON,
		output:      filterStatus,
	}

	if err := g.sendRequest(params); err != nil {
		return fmt.Errorf(
			"received an error after sending the request to add the filter status to the filter: %w",
			err,
		)
	}

	return nil
}

func (g *GTSClient) DeleteFilterStatus(filterStatusID string, _ *NoRPCResults) error {
	params := requestParameters{
		httpMethod:  http.MethodDelete,
		url:         g.authentication.Instance + baseFilterStatusesPath + "/" + filterStatusID,
		requestBody: nil,
		contentType: "",
		output:      nil,
	}

	if err := g.sendRequest(params); err != nil {
		return fmt.Errorf(
			"received an error after sending the request to delete the filter status: %w",
			err,
		)
	}

	return nil
}

func (g *GTSClient) GetFilterStatus(filterStatusID string, filterStatus *model.FilterStatus) error {
	params := requestParameters{
		httpMethod:  http.MethodGet,
		url:         g.authentication.Instance + baseFilterStatusesPath + "/" + filterStatusID,
		requestBody: nil,
		contentType: "",
		output:      filterStatus,
	}

	if err := g.sendRequest(params); err != nil {
		return fmt.Errorf(
			"received an error after sending the request to get the filter status: %w",
			err,
		)
	}

	return nil
}
