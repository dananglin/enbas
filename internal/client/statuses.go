// SPDX-FileCopyrightText: 2024 Dan Anglin <d.n.i.anglin@gmail.com>
//
// SPDX-License-Identifier: GPL-3.0-or-later

package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"codeflow.dananglin.me.uk/apollo/enbas/internal/model"
)

const (
	baseStatusesPath string = "/api/v1/statuses"
)

func (g *Client) GetStatus(statusID string) (model.Status, error) {
	path := baseStatusesPath + "/" + statusID
	url := g.Authentication.Instance + path

	var status model.Status

	if err := g.sendRequest(http.MethodGet, url, nil, &status); err != nil {
		return model.Status{}, fmt.Errorf(
			"received an error after sending the request to get the status information: %w",
			err,
		)
	}

	return status, nil
}

type CreateStatusForm struct {
	Content     string                  `json:"status"`
	InReplyTo   string                  `json:"in_reply_to_id"`
	Language    string                  `json:"language"`
	SpoilerText string                  `json:"spoiler_text"`
	Boostable   bool                    `json:"boostable"`
	Federated   bool                    `json:"federated"`
	Likeable    bool                    `json:"likeable"`
	Replyable   bool                    `json:"replyable"`
	Sensitive   bool                    `json:"sensitive"`
	Poll        *CreateStatusPollForm   `json:"poll,omitempty"`
	ContentType model.StatusContentType `json:"content_type"`
	Visibility  model.StatusVisibility  `json:"visibility"`
}

type CreateStatusPollForm struct {
	Options    []string `json:"options"`
	ExpiresIn  int      `json:"expires_in"`
	Multiple   bool     `json:"multiple"`
	HideTotals bool     `json:"hide_totals"`
}

func (g *Client) CreateStatus(form CreateStatusForm) (model.Status, error) {
	data, err := json.Marshal(form)
	if err != nil {
		return model.Status{}, fmt.Errorf("unable to create the JSON form: %w", err)
	}

	requestBody := bytes.NewBuffer(data)
	url := g.Authentication.Instance + baseStatusesPath

	var status model.Status

	if err := g.sendRequest(http.MethodPost, url, requestBody, &status); err != nil {
		return model.Status{}, fmt.Errorf(
			"received an error after sending the request to create the status: %w",
			err,
		)
	}

	return status, nil
}

func (g *Client) GetBookmarks(limit int) (model.StatusList, error) {
	path := fmt.Sprintf("/api/v1/bookmarks?limit=%d", limit)
	url := g.Authentication.Instance + path

	bookmarks := model.StatusList{
		Name:     "Your Bookmarks",
		Statuses: nil,
	}

	if err := g.sendRequest(http.MethodGet, url, nil, &bookmarks.Statuses); err != nil {
		return bookmarks, fmt.Errorf(
			"received an error after sending the request to get the bookmarks: %w",
			err,
		)
	}

	return bookmarks, nil
}

func (g *Client) AddStatusToBookmarks(statusID string) error {
	path := fmt.Sprintf("/api/v1/statuses/%s/bookmark", statusID)
	url := g.Authentication.Instance + path

	if err := g.sendRequest(http.MethodPost, url, nil, nil); err != nil {
		return fmt.Errorf(
			"received an error after sending the request to add the status to the list of bookmarks: %w",
			err,
		)
	}

	return nil
}

func (g *Client) RemoveStatusFromBookmarks(statusID string) error {
	path := fmt.Sprintf("/api/v1/statuses/%s/unbookmark", statusID)
	url := g.Authentication.Instance + path

	if err := g.sendRequest(http.MethodPost, url, nil, nil); err != nil {
		return fmt.Errorf(
			"received an error after sending the request to remove the status from the list of bookmarks: %w",
			err,
		)
	}

	return nil
}

func (g *Client) LikeStatus(statusID string) error {
	url := g.Authentication.Instance + baseStatusesPath + "/" + statusID + "/favourite"

	if err := g.sendRequest(http.MethodPost, url, nil, nil); err != nil {
		return fmt.Errorf(
			"received an error after sending the request to like the status: %w",
			err,
		)
	}

	return nil
}

func (g *Client) UnlikeStatus(statusID string) error {
	url := g.Authentication.Instance + baseStatusesPath + "/" + statusID + "/unfavourite"

	if err := g.sendRequest(http.MethodPost, url, nil, nil); err != nil {
		return fmt.Errorf(
			"received an error after sending the request to unlike the status: %w",
			err,
		)
	}

	return nil
}

func (g *Client) GetLikedStatuses(limit int, resourceName string) (model.StatusList, error) {
	url := g.Authentication.Instance + fmt.Sprintf("/api/v1/favourites?limit=%d", limit)

	liked := model.StatusList{
		Name:     "Your " + resourceName + " statuses",
		Statuses: nil,
	}

	if err := g.sendRequest(http.MethodGet, url, nil, &liked.Statuses); err != nil {
		return model.StatusList{}, fmt.Errorf(
			"received an error after sending the request to get the list of statuses: %w",
			err,
		)
	}

	return liked, nil
}

func (g *Client) ReblogStatus(statusID string) error {
	url := g.Authentication.Instance + baseStatusesPath + "/" + statusID + "/reblog"

	if err := g.sendRequest(http.MethodPost, url, nil, nil); err != nil {
		return fmt.Errorf(
			"received an error after sending the request to reblog the status; %w",
			err,
		)
	}

	return nil
}

func (g *Client) UnreblogStatus(statusID string) error {
	url := g.Authentication.Instance + baseStatusesPath + "/" + statusID + "/unreblog"

	if err := g.sendRequest(http.MethodPost, url, nil, nil); err != nil {
		return fmt.Errorf(
			"received an error after sending the request to un-reblog the status; %w",
			err,
		)
	}

	return nil
}
