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
