package gtsclient

import (
	"fmt"
	"net/http"

	"codeflow.dananglin.me.uk/apollo/enbas/internal/model"
)

const baseSearchPath string = "/api/v2/search"

type searchResults struct {
	Accounts []model.Account `json:"accounts"`
	Tags     []model.Tag     `json:"hashtags"`
	Statuses []model.Status  `json:"statuses"`
}

type SearchTagsArgs struct {
	Limit int
	Query string
}

func (g *GTSClient) SearchTags(args SearchTagsArgs, list *model.TagList) error {
	query := fmt.Sprintf(
		"?type=hashtags&limit=%d&q=%s",
		args.Limit,
		args.Query,
	)

	var results searchResults

	params := requestParameters{
		httpMethod:  http.MethodGet,
		url:         g.authentication.Instance + baseSearchPath + query,
		requestBody: nil,
		contentType: "",
		output:      &results,
	}

	if err := g.sendRequest(params); err != nil {
		return fmt.Errorf(
			"received an error after sending the search request for tags: %w",
			err,
		)
	}

	*list = model.TagList{
		Name: "Search results",
		Tags: results.Tags,
	}

	return nil
}

type SearchAccountsArgs struct {
	Limit     int
	Query     string
	Resolve   bool
	Following bool
}

func (g *GTSClient) SearchAccounts(args SearchAccountsArgs, list *model.AccountList) error {
	query := fmt.Sprintf(
		"?type=accounts&limit=%d&q=%s&resolve=%t&following=%t",
		args.Limit,
		args.Query,
		args.Resolve,
		args.Following,
	)

	var results searchResults

	params := requestParameters{
		httpMethod:  http.MethodGet,
		url:         g.authentication.Instance + baseSearchPath + query,
		requestBody: nil,
		contentType: "",
		output:      &results,
	}

	if err := g.sendRequest(params); err != nil {
		return fmt.Errorf(
			"received an error after sending the search request for accounts: %w",
			err,
		)
	}

	*list = model.AccountList{
		Label:           "Search results",
		Accounts:        results.Accounts,
		BlockedAccounts: false,
	}

	return nil
}

type SearchStatusesArgs struct {
	Limit     int
	Query     string
	AccountID string
	Resolve   bool
}

func (g *GTSClient) SearchStatuses(args SearchStatusesArgs, list *model.StatusList) error {
	query := fmt.Sprintf(
		"?type=statuses&limit=%d&q=%s&resolve=%t",
		args.Limit,
		args.Query,
		args.Resolve,
	)

	if args.AccountID != "" {
		query = query + "&account_id=" + args.AccountID
	}

	var results searchResults

	params := requestParameters{
		httpMethod:  http.MethodGet,
		url:         g.authentication.Instance + baseSearchPath + query,
		requestBody: nil,
		contentType: "",
		output:      &results,
	}

	if err := g.sendRequest(params); err != nil {
		return fmt.Errorf(
			"received an error after sending the search request for statuses: %w",
			err,
		)
	}

	*list = model.StatusList{
		Name:     "Search results",
		Statuses: results.Statuses,
	}

	return nil
}
