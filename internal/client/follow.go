package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"codeflow.dananglin.me.uk/apollo/enbas/internal/model"
)

func (g *Client) FollowAccount(accountID string, reblogs, notify bool) error {
	form := struct {
		ID      string `json:"id"`
		Reblogs bool   `json:"reblogs"`
		Notify  bool   `json:"notify"`
	}{
		ID:      accountID,
		Reblogs: reblogs,
		Notify:  notify,
	}

	data, err := json.Marshal(form)
	if err != nil {
		return fmt.Errorf("unable to marshal the form; %w", err)
	}

	requestBody := bytes.NewBuffer(data)
	url := g.Authentication.Instance + fmt.Sprintf("/api/v1/accounts/%s/follow", accountID)

	if err := g.sendRequest(http.MethodPost, url, requestBody, nil); err != nil {
		return fmt.Errorf("received an error after sending the follow request; %w", err)
	}

	return nil
}

func (g *Client) UnfollowAccount(accountID string) error {
	url := g.Authentication.Instance + fmt.Sprintf("/api/v1/accounts/%s/unfollow", accountID)

	if err := g.sendRequest(http.MethodPost, url, nil, nil); err != nil {
		return fmt.Errorf("received an error after sending the request to unfollow the account; %w", err)
	}

	return nil
}

func (g *Client) GetFollowers(accountID string, limit int) (model.Followers, error) {
	url := g.Authentication.Instance + fmt.Sprintf("/api/v1/accounts/%s/followers?limit=%d", accountID, limit)

	var followers model.Followers

	if err := g.sendRequest(http.MethodGet, url, nil, &followers); err != nil {
		return nil, fmt.Errorf("received an error after sending the request to get the list of followers; %w", err)
	}

	return followers, nil
}

func (g *Client) GetFollowing(accountID string, limit int) (model.Following, error) {
	url := g.Authentication.Instance + fmt.Sprintf("/api/v1/accounts/%s/following?limit=%d", accountID, limit)

	var following model.Following

	if err := g.sendRequest(http.MethodGet, url, nil, &following); err != nil {
		return nil, fmt.Errorf("received an error after sending the request to get the list of followed accounts; %w", err)
	}

	return following, nil
}
