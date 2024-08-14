package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	pollPath string = "/api/v1/polls"
)

// func (g *Client) GetPoll(pollID string) (model.Poll, error) {
// 	url := g.Authentication.Instance + pollPath + "/" + pollID
//
// 	var poll model.Poll
//
// 	params := requestParameters{
// 		httpMethod:  http.MethodGet,
// 		url:         url,
// 		requestBody: nil,
// 		contentType: "",
// 		output:      &poll,
// 	}
//
// 	if err := g.sendRequest(params); err != nil {
// 		return model.Poll{}, fmt.Errorf(
// 			"received an error after sending the request to get the poll: %w",
// 			err,
// 		)
// 	}
//
// 	return poll, nil
// }

func (g *Client) VoteInPoll(pollID string, choices []int) error {
	form := struct {
		Choices []int `json:"choices"`
	}{
		Choices: choices,
	}

	data, err := json.Marshal(form)
	if err != nil {
		return fmt.Errorf("unable to encode the JSON form: %w", err)
	}

	requestBody := bytes.NewBuffer(data)
	url := g.Authentication.Instance + pollPath + "/" + pollID + "/votes"

	params := requestParameters{
		httpMethod:  http.MethodPost,
		url:         url,
		requestBody: requestBody,
		contentType: applicationJSON,
		output:      nil,
	}

	if err := g.sendRequest(params); err != nil {
		return fmt.Errorf("received an error after sending the request to vote in the poll: %w", err)
	}

	return nil
}
