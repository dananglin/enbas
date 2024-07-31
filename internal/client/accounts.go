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
	baseAccountsPath       = "/api/v1/accounts"
	baseFollowRequestsPath = "/api/v1/follow_requests"
)

func (g *Client) VerifyCredentials() (model.Account, error) {
	url := g.Authentication.Instance + baseAccountsPath + "/verify_credentials"

	var account model.Account

	if err := g.sendRequest(http.MethodGet, url, nil, &account); err != nil {
		return model.Account{}, fmt.Errorf("received an error after sending the request to verify the credentials: %w", err)
	}

	return account, nil
}

func (g *Client) GetAccount(accountURI string) (model.Account, error) {
	url := g.Authentication.Instance + baseAccountsPath + "/lookup?acct=" + accountURI

	var account model.Account

	if err := g.sendRequest(http.MethodGet, url, nil, &account); err != nil {
		return model.Account{}, fmt.Errorf("received an error after sending the request to get the account information: %w", err)
	}

	return account, nil
}

func (g *Client) GetAccountRelationship(accountID string) (*model.AccountRelationship, error) {
	url := g.Authentication.Instance + baseAccountsPath + "/relationships?id=" + accountID

	var relationships []model.AccountRelationship

	if err := g.sendRequest(http.MethodGet, url, nil, &relationships); err != nil {
		return nil, fmt.Errorf(
			"received an error after sending the request to get the account relationship: %w",
			err,
		)
	}

	if len(relationships) != 1 {
		return nil, fmt.Errorf(
			"unexpected number of account relationships returned: want 1, got %d",
			len(relationships),
		)
	}

	return &relationships[0], nil
}

type FollowAccountForm struct {
	AccountID   string `json:"id"`
	ShowReposts bool   `json:"reblogs"`
	Notify      bool   `json:"notify"`
}

func (g *Client) FollowAccount(form FollowAccountForm) error {
	data, err := json.Marshal(form)
	if err != nil {
		return fmt.Errorf("unable to marshal the form: %w", err)
	}

	requestBody := bytes.NewBuffer(data)
	url := g.Authentication.Instance + baseAccountsPath + "/" + form.AccountID + "/follow"

	if err := g.sendRequest(http.MethodPost, url, requestBody, nil); err != nil {
		return fmt.Errorf("received an error after sending the follow request: %w", err)
	}

	return nil
}

func (g *Client) UnfollowAccount(accountID string) error {
	url := g.Authentication.Instance + baseAccountsPath + "/" + accountID + "/unfollow"

	if err := g.sendRequest(http.MethodPost, url, nil, nil); err != nil {
		return fmt.Errorf("received an error after sending the request to unfollow the account: %w", err)
	}

	return nil
}

func (g *Client) GetFollowers(accountID string, limit int) (model.AccountList, error) {
	url := g.Authentication.Instance + fmt.Sprintf("%s/%s/followers?limit=%d", baseAccountsPath, accountID, limit)

	accounts := make([]model.Account, limit)

	if err := g.sendRequest(http.MethodGet, url, nil, &accounts); err != nil {
		return model.AccountList{}, fmt.Errorf("received an error after sending the request to get the list of followers: %w", err)
	}

	followers := model.AccountList{
		Type:     model.AccountListFollowers,
		Accounts: accounts,
	}

	return followers, nil
}

func (g *Client) GetFollowing(accountID string, limit int) (model.AccountList, error) {
	url := g.Authentication.Instance + fmt.Sprintf("%s/%s/following?limit=%d", baseAccountsPath, accountID, limit)

	accounts := make([]model.Account, limit)

	if err := g.sendRequest(http.MethodGet, url, nil, &accounts); err != nil {
		return model.AccountList{}, fmt.Errorf("received an error after sending the request to get the list of followed accounts: %w", err)
	}

	following := model.AccountList{
		Type:     model.AccountListFollowing,
		Accounts: accounts,
	}

	return following, nil
}

func (g *Client) BlockAccount(accountID string) error {
	url := g.Authentication.Instance + baseAccountsPath + "/" + accountID + "/block"

	if err := g.sendRequest(http.MethodPost, url, nil, nil); err != nil {
		return fmt.Errorf("received an error after sending the request to block the account: %w", err)
	}

	return nil
}

func (g *Client) UnblockAccount(accountID string) error {
	url := g.Authentication.Instance + baseAccountsPath + "/" + accountID + "/unblock"

	if err := g.sendRequest(http.MethodPost, url, nil, nil); err != nil {
		return fmt.Errorf("received an error after sending the request to unblock the account: %w", err)
	}

	return nil
}

func (g *Client) GetBlockedAccounts(limit int) (model.AccountList, error) {
	url := g.Authentication.Instance + fmt.Sprintf("/api/v1/blocks?limit=%d", limit)

	var accounts []model.Account

	if err := g.sendRequest(http.MethodGet, url, nil, &accounts); err != nil {
		return model.AccountList{}, fmt.Errorf("received an error after sending the request to get the list of blocked accounts: %w", err)
	}

	blocked := model.AccountList{
		Type:     model.AccountListBlockedAccount,
		Accounts: accounts,
	}

	return blocked, nil
}

func (g *Client) SetPrivateNote(accountID, note string) error {
	form := struct {
		Comment string `json:"comment"`
	}{
		Comment: note,
	}

	data, err := json.Marshal(form)
	if err != nil {
		return fmt.Errorf("unable to marshal the form: %w", err)
	}

	requestBody := bytes.NewBuffer(data)
	url := g.Authentication.Instance + baseAccountsPath + "/" + accountID + "/note"

	if err := g.sendRequest(http.MethodPost, url, requestBody, nil); err != nil {
		return fmt.Errorf("received an error after sending the request to set the private note: %w", err)
	}

	return nil
}

func (g *Client) GetFollowRequests(limit int) (model.AccountList, error) {
	url := g.Authentication.Instance + fmt.Sprintf("%s?limit=%d", baseFollowRequestsPath, limit)

	var accounts []model.Account

	if err := g.sendRequest(http.MethodGet, url, nil, &accounts); err != nil {
		return model.AccountList{}, fmt.Errorf("received an error after sending the request to get the list of follow requests: %w", err)
	}

	requests := model.AccountList{
		Type:     model.AccountListFollowRequests,
		Accounts: accounts,
	}

	return requests, nil
}

func (g *Client) AcceptFollowRequest(accountID string) error {
	url := g.Authentication.Instance + baseFollowRequestsPath + "/" + accountID + "/authorize"

	if err := g.sendRequest(http.MethodPost, url, nil, nil); err != nil {
		return fmt.Errorf("received an error after sending the request to accept the follow request: %w", err)
	}

	return nil
}

func (g *Client) RejectFollowRequest(accountID string) error {
	url := g.Authentication.Instance + baseFollowRequestsPath + "/" + accountID + "/reject"

	if err := g.sendRequest(http.MethodPost, url, nil, nil); err != nil {
		return fmt.Errorf("received an error after sending the request to reject the follow request: %w", err)
	}

	return nil
}

func (g *Client) GetMutedAccounts(limit int) (model.AccountList, error) {
	url := g.Authentication.Instance + fmt.Sprintf("/api/v1/mutes?limit=%d", limit)

	var accounts []model.Account

	if err := g.sendRequest(http.MethodGet, url, nil, &accounts); err != nil {
		return model.AccountList{}, fmt.Errorf("received an error after sending the request to get the list of muted accounts: %w", err)
	}

	muted := model.AccountList{
		Type:     model.AccountListMuted,
		Accounts: accounts,
	}

	return muted, nil
}

type MuteAccountForm struct {
	Notifications bool `json:"notifications"`
	Duration      int  `json:"duration"`
}

func (g *Client) MuteAccount(accountID string, form MuteAccountForm) error {
	data, err := json.Marshal(form)
	if err != nil {
		return fmt.Errorf("unable to marshal the form: %w", err)
	}

	requestBody := bytes.NewBuffer(data)
	url := g.Authentication.Instance + baseAccountsPath + "/" + accountID + "/mute"

	if err := g.sendRequest(http.MethodPost, url, requestBody, nil); err != nil {
		return fmt.Errorf("received an error after sending the request to mute the account: %w", err)
	}

	return nil
}

func (g *Client) UnmuteAccount(accountID string) error {
	url := g.Authentication.Instance + baseAccountsPath + "/" + accountID + "/unmute"

	if err := g.sendRequest(http.MethodPost, url, nil, nil); err != nil {
		return fmt.Errorf("received an error after sending the request to unmute the account: %w", err)
	}

	return nil
}

type GetAccountStatusesForm struct {
	AccountID      string
	Limit          int
	ExcludeReplies bool
	ExcludeReblogs bool
	Pinned         bool
	OnlyMedia      bool
	OnlyPublic     bool
}

func (g *Client) GetAccountStatuses(form GetAccountStatusesForm) (*model.StatusList, error) {
	path := baseAccountsPath + "/" + form.AccountID + "/statuses"
	query := fmt.Sprintf(
		"?limit=%d&exclude_replies=%t&exclude_reblogs=%t&pinned=%t&only_media=%t&only_public=%t",
		form.Limit,
		form.ExcludeReplies,
		form.ExcludeReblogs,
		form.Pinned,
		form.OnlyMedia,
		form.OnlyPublic,
	)
	url := g.Authentication.Instance + path + query

	var statuses []model.Status

	if err := g.sendRequest(http.MethodGet, url, nil, &statuses); err != nil {
		return nil, fmt.Errorf("received an error after sending the request to get the account's statuses: %w", err)
	}

	statusList := model.StatusList{
		Name:     "STATUSES:",
		Statuses: statuses,
	}

	return &statusList, nil
}
