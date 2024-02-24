package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"codeflow.dananglin.me.uk/apollo/enbas/internal"
	"codeflow.dananglin.me.uk/apollo/enbas/internal/config"
	"codeflow.dananglin.me.uk/apollo/enbas/internal/model"
)

type Client struct {
	Authentication config.Authentication
	HTTPClient     http.Client
	UserAgent      string
	Timeout        time.Duration
}

func NewClientFromConfig() (*Client, error) {
	config, err := config.NewAuthenticationConfigFromFile()
	if err != nil {
		return nil, fmt.Errorf("unable to get the authentication configuration; %w", err)
	}

	currentAuthentication := config.Authentications[config.CurrentAccount]

	return NewClient(currentAuthentication), nil
}

func NewClient(authentication config.Authentication) *Client {
	httpClient := http.Client{}

	gtsClient := Client{
		Authentication: authentication,
		HTTPClient:     httpClient,
		UserAgent:      internal.UserAgent,
		Timeout:        5 * time.Second,
	}

	return &gtsClient
}

func (g *Client) VerifyCredentials() (model.Account, error) {
	path := "/api/v1/accounts/verify_credentials"
	url := g.Authentication.Instance + path

	var account model.Account

	if err := g.sendRequest(http.MethodGet, url, nil, &account); err != nil {
		return model.Account{}, fmt.Errorf("received an error after sending the request to verify the credentials; %w", err)
	}

	return account, nil
}

func (g *Client) GetInstance() (model.InstanceV2, error) {
	path := "/api/v2/instance"
	url := g.Authentication.Instance + path

	var instance model.InstanceV2

	if err := g.sendRequest(http.MethodGet, url, nil, &instance); err != nil {
		return model.InstanceV2{}, fmt.Errorf("received an error after sending the request to get the instance details; %w", err)
	}

	return instance, nil
}

func (g *Client) GetAccount(accountURI string) (model.Account, error) {
	path := "/api/v1/accounts/lookup?acct=" + accountURI
	url := g.Authentication.Instance + path

	var account model.Account

	if err := g.sendRequest(http.MethodGet, url, nil, &account); err != nil {
		return model.Account{}, fmt.Errorf("received an error after sending the request to get the account information; %w", err)
	}

	return account, nil
}

func (g *Client) GetStatus(statusID string) (model.Status, error) {
	path := "/api/v1/statuses/" + statusID
	url := g.Authentication.Instance + path

	var status model.Status

	if err := g.sendRequest(http.MethodGet, url, nil, &status); err != nil {
		return model.Status{}, fmt.Errorf("received an error after sending the request to get the status information; %w", err)
	}

	return status, nil
}

func (g *Client) GetHomeTimeline(limit int) (model.Timeline, error) {
	path := fmt.Sprintf("/api/v1/timelines/home?limit=%d", limit)

	timeline := model.Timeline{
		Name:     "HOME TIMELINE",
		Statuses: nil,
	}

	return g.getTimeline(path, timeline)
}

func (g *Client) GetPublicTimeline(limit int) (model.Timeline, error) {
	path := fmt.Sprintf("/api/v1/timelines/public?limit=%d", limit)

	timeline := model.Timeline{
		Name:     "PUBLIC TIMELINE",
		Statuses: nil,
	}

	return g.getTimeline(path, timeline)
}

func (g *Client) GetListTimeline(listID string, limit int) (model.Timeline, error) {
	path := fmt.Sprintf("/api/v1/timelines/list/%s?limit=%d", listID, limit)

	timeline := model.Timeline{
		Name:     "LIST: " + listID,
		Statuses: nil,
	}

	return g.getTimeline(path, timeline)
}

func (g *Client) GetTagTimeline(tag string, limit int) (model.Timeline, error) {
	path := fmt.Sprintf("/api/v1/timelines/tag/%s?limit=%d", tag, limit)

	timeline := model.Timeline{
		Name:     "TAG: " + tag,
		Statuses: nil,
	}

	return g.getTimeline(path, timeline)
}

func (g *Client) getTimeline(path string, timeline model.Timeline) (model.Timeline, error) {
	url := g.Authentication.Instance + path

	var statuses []model.Status

	if err := g.sendRequest(http.MethodGet, url, nil, &statuses); err != nil {
		return timeline, fmt.Errorf("received an error after sending the request to get the timeline; %w", err)
	}

	timeline.Statuses = statuses

	return timeline, nil
}

func (g *Client) sendRequest(method string, url string, requestBody io.Reader, object any) error {
	ctx, cancel := context.WithTimeout(context.Background(), g.Timeout)
	defer cancel()

	request, err := http.NewRequestWithContext(ctx, method, url, requestBody)
	if err != nil {
		return fmt.Errorf("unable to create the HTTP request, %w", err)
	}

	request.Header.Set("Content-Type", "application/json; charset=utf-8")
	request.Header.Set("Accept", "application/json; charset=utf-8")
	request.Header.Set("User-Agent", g.UserAgent)

	if len(g.Authentication.AccessToken) > 0 {
		request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", g.Authentication.AccessToken))
	}

	response, err := g.HTTPClient.Do(request)
	if err != nil {
		return fmt.Errorf("received an error after sending the request; %w", err)
	}

	defer response.Body.Close()

	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusBadRequest {
		return fmt.Errorf(
			"did not receive an OK response from the GoToSocial server; got %d",
			response.StatusCode,
		)
	}

	if err := json.NewDecoder(response.Body).Decode(object); err != nil {
		return fmt.Errorf(
			"unable to decode the response from the GoToSocial server; %w",
			err,
		)
	}

	return nil
}
