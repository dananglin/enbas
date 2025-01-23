package gtsclient

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type requestParameters struct {
	httpMethod  string
	url         string
	contentType string
	requestBody io.Reader
	output      any
}

func (g *GTSClient) sendRequest(params requestParameters) error {
	ctx, cancel := context.WithTimeout(context.Background(), g.Timeout)
	defer cancel()

	request, err := http.NewRequestWithContext(ctx, params.httpMethod, params.url, params.requestBody)
	if err != nil {
		return fmt.Errorf("unable to create the HTTP request: %w", err)
	}

	if params.contentType != "" {
		request.Header.Set("Content-Type", params.contentType)
	}

	request.Header.Set("Accept", applicationJSON)
	request.Header.Set("User-Agent", userAgent)

	if len(g.Authentication.AccessToken) > 0 {
		request.Header.Set("Authorization", "Bearer "+g.Authentication.AccessToken)
	}

	response, err := g.HTTPClient.Do(request)
	if err != nil {
		return fmt.Errorf("received an error after sending the request: %w", err)
	}

	defer response.Body.Close()

	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusBadRequest {
		message := struct {
			Error string `json:"error"`
		}{
			Error: "",
		}

		if err := json.NewDecoder(response.Body).Decode(&message); err != nil {
			return ResponseError{
				StatusCode:       response.StatusCode,
				Message:          "",
				MessageDecodeErr: err,
			}
		}

		return ResponseError{
			StatusCode:       response.StatusCode,
			Message:          message.Error,
			MessageDecodeErr: nil,
		}
	}

	if params.output == nil {
		return nil
	}

	if err := json.NewDecoder(response.Body).Decode(params.output); err != nil {
		return fmt.Errorf(
			"unable to decode the response from the GoToSocial server: %w",
			err,
		)
	}

	return nil
}
