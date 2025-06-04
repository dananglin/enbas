package gtsclient

import (
	"fmt"
	"net/http"

	"codeflow.dananglin.me.uk/apollo/enbas/internal/model"
)

const baseTokenPath string = "/api/v1/tokens"

func (g *GTSClient) GetTokens(limit int, tokenList *model.TokenList) error {
	var tokens []model.Token

	params := requestParameters{
		httpMethod:  http.MethodGet,
		url:         g.auth.GetInstanceURL() + fmt.Sprintf("%s?limit=%d", baseTokenPath, limit),
		requestBody: nil,
		contentType: "",
		output:      &tokens,
	}

	if err := g.sendRequest(params); err != nil {
		return fmt.Errorf(
			"received an error after sending the request to get the tokens: %w",
			err,
		)
	}

	*tokenList = model.TokenList{
		Label:  "Your tokens",
		Tokens: tokens,
	}

	return nil
}

func (g *GTSClient) GetToken(tokenID string, token *model.Token) error {
	params := requestParameters{
		httpMethod:  http.MethodGet,
		url:         g.auth.GetInstanceURL() + baseTokenPath + "/" + tokenID,
		requestBody: nil,
		contentType: "",
		output:      token,
	}

	if err := g.sendRequest(params); err != nil {
		return fmt.Errorf(
			"received an error after sending the request to get the token: %w",
			err,
		)
	}

	return nil
}

func (g *GTSClient) InvalidateToken(tokenID string, _ *NoRPCResults) error {
	params := requestParameters{
		httpMethod:  http.MethodPost,
		url:         g.auth.GetInstanceURL() + baseTokenPath + "/" + tokenID + "/invalidate",
		requestBody: nil,
		contentType: "",
		output:      nil,
	}

	if err := g.sendRequest(params); err != nil {
		return fmt.Errorf(
			"received an error after sending the request to invalidate the token: %w",
			err,
		)
	}

	return nil
}
