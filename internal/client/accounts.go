package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"codeflow.dananglin.me.uk/apollo/enbas/internal/model"
)

func (g *Client) VerifyCredentials() (model.Account, error) {
	path := "/api/v1/accounts/verify_credentials"
	url := g.Authentication.Instance + path

	var account model.Account

	if err := g.sendRequest(http.MethodGet, url, nil, &account); err != nil {
		return model.Account{}, fmt.Errorf("received an error after sending the request to verify the credentials; %w", err)
	}

	return account, nil
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

func (g *Client) GetAccountRelationship(accountID string) (model.AccountRelationship, error) {
	path := "/api/v1/accounts/relationships?id=" + accountID
	url := g.Authentication.Instance + path

	var relationships []model.AccountRelationship

	if err := g.sendRequest(http.MethodGet, url, nil, &relationships); err != nil {
		return model.AccountRelationship{}, fmt.Errorf("received an error after sending the request to get the account relationship; %w", err)
	}

	if len(relationships) != 1 {
		return model.AccountRelationship{}, fmt.Errorf("unexpected number of account relationships returned; want 1, got %d", len(relationships))
	}

	return relationships[0], nil
}

type FollowAccountForm struct {
	AccountID   string `json:"id"`
	ShowReposts bool   `json:"reblogs"`
	Notify      bool   `json:"notify"`
}

func (g *Client) FollowAccount(form FollowAccountForm) error {
	data, err := json.Marshal(form)
	if err != nil {
		return fmt.Errorf("unable to marshal the form; %w", err)
	}

	requestBody := bytes.NewBuffer(data)
	url := g.Authentication.Instance + fmt.Sprintf("/api/v1/accounts/%s/follow", form.AccountID)

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

func (g *Client) GetFollowers(accountID string, limit int) (model.AccountList, error) {
	url := g.Authentication.Instance + fmt.Sprintf("/api/v1/accounts/%s/followers?limit=%d", accountID, limit)

	accounts := make([]model.Account, limit)

	if err := g.sendRequest(http.MethodGet, url, nil, &accounts); err != nil {
		return model.AccountList{}, fmt.Errorf("received an error after sending the request to get the list of followers; %w", err)
	}

	followers := model.AccountList{
		Type:     model.AccountListFollowers,
		Accounts: accounts,
	}

	return followers, nil
}

func (g *Client) GetFollowing(accountID string, limit int) (model.AccountList, error) {
	url := g.Authentication.Instance + fmt.Sprintf("/api/v1/accounts/%s/following?limit=%d", accountID, limit)

	accounts := make([]model.Account, limit)

	if err := g.sendRequest(http.MethodGet, url, nil, &accounts); err != nil {
		return model.AccountList{}, fmt.Errorf("received an error after sending the request to get the list of followed accounts; %w", err)
	}

	following := model.AccountList{
		Type:     model.AccountListFollowing,
		Accounts: accounts,
	}

	return following, nil
}

func (g *Client) BlockAccount(accountID string) error {
	url := g.Authentication.Instance + fmt.Sprintf("/api/v1/accounts/%s/block", accountID)

	if err := g.sendRequest(http.MethodPost, url, nil, nil); err != nil {
		return fmt.Errorf("received an error after sending the request to block the account; %w", err)
	}

	return nil
}

func (g *Client) UnblockAccount(accountID string) error {
	url := g.Authentication.Instance + fmt.Sprintf("/api/v1/accounts/%s/unblock", accountID)

	if err := g.sendRequest(http.MethodPost, url, nil, nil); err != nil {
		return fmt.Errorf("received an error after sending the request to unblock the account; %w", err)
	}

	return nil
}

func (g *Client) GetBlockedAccounts(limit int) (model.AccountList, error) {
	url := g.Authentication.Instance + fmt.Sprintf("/api/v1/blocks?limit=%d", limit)

	accounts := make([]model.Account, limit)

	if err := g.sendRequest(http.MethodGet, url, nil, &accounts); err != nil {
		return model.AccountList{}, fmt.Errorf("received an error after sending the request to get the list of blocked accounts; %w", err)
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
		return fmt.Errorf("unable to marshal the form; %w", err)
	}

	requestBody := bytes.NewBuffer(data)
	url := g.Authentication.Instance + fmt.Sprintf("/api/v1/accounts/%s/note", accountID)

	if err := g.sendRequest(http.MethodPost, url, requestBody, nil); err != nil {
		return fmt.Errorf("received an error after sending the request to set the private note; %w", err)
	}

	return nil
}
