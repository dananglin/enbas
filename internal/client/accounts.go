package client

import (
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
