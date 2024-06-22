// SPDX-FileCopyrightText: 2024 Dan Anglin <d.n.i.anglin@gmail.com>
//
// SPDX-License-Identifier: GPL-3.0-or-later

package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"

	"codeflow.dananglin.me.uk/apollo/enbas/internal"
	"codeflow.dananglin.me.uk/apollo/enbas/internal/config"
)

type Client struct {
	Authentication config.Credentials
	HTTPClient     http.Client
	UserAgent      string
	Timeout        time.Duration
}

func NewClientFromConfig(configDir string) (*Client, error) {
	config, err := config.NewCredentialsConfigFromFile(configDir)
	if err != nil {
		return nil, fmt.Errorf("unable to get the authentication configuration: %w", err)
	}

	currentAuthentication := config.Credentials[config.CurrentAccount]

	return NewClient(currentAuthentication), nil
}

func NewClient(authentication config.Credentials) *Client {
	httpClient := http.Client{}

	gtsClient := Client{
		Authentication: authentication,
		HTTPClient:     httpClient,
		UserAgent:      internal.UserAgent,
		Timeout:        5 * time.Second,
	}

	return &gtsClient
}

func (g *Client) AuthCodeURL() string {
	format := "%s/oauth/authorize?client_id=%s&redirect_uri=%s&response_type=code"
	escapedRedirectURI := url.QueryEscape(internal.RedirectUri)

	return fmt.Sprintf(
		format,
		g.Authentication.Instance,
		g.Authentication.ClientID,
		escapedRedirectURI,
	)
}

func (g *Client) DownloadMedia(url, path string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("unable to create the HTTP request: %w", err)
	}

	request.Header.Set("User-Agent", g.UserAgent)

	response, err := g.HTTPClient.Do(request)
	if err != nil {
		return fmt.Errorf("received an error after attempting the download: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf(
			"did not receive an OK response from the GoToSocial server: got %d",
			response.StatusCode,
		)
	}

	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("unable to create %s: %w", path, err)
	}
	defer file.Close()

	if _, err = io.Copy(file, response.Body); err != nil {
		return fmt.Errorf("unable to save the download to %s: %w", path, err)
	}

	return nil
}

func (g *Client) sendRequest(method string, url string, requestBody io.Reader, object any) error {
	ctx, cancel := context.WithTimeout(context.Background(), g.Timeout)
	defer cancel()

	request, err := http.NewRequestWithContext(ctx, method, url, requestBody)
	if err != nil {
		return fmt.Errorf("unable to create the HTTP request: %w", err)
	}

	request.Header.Set("Content-Type", "application/json; charset=utf-8")
	request.Header.Set("Accept", "application/json; charset=utf-8")
	request.Header.Set("User-Agent", g.UserAgent)

	if len(g.Authentication.AccessToken) > 0 {
		request.Header.Set("Authorization", "Bearer "+g.Authentication.AccessToken)
	}

	response, err := g.HTTPClient.Do(request)
	if err != nil {
		return fmt.Errorf("received an error after sending the request: %w", err)
	}

	defer response.Body.Close()

	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusBadRequest {
		return fmt.Errorf(
			"did not receive an OK response from the GoToSocial server: got %d",
			response.StatusCode,
		)
	}

	if object == nil {
		return nil
	}

	if err := json.NewDecoder(response.Body).Decode(object); err != nil {
		return fmt.Errorf(
			"unable to decode the response from the GoToSocial server: %w",
			err,
		)
	}

	return nil
}
