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
	pollPath string = "/api/v1/polls"
)

func (g *Client) GetPoll(pollID string) (model.Poll, error) {
	url := g.Authentication.Instance + pollPath + "/" + pollID

	var poll model.Poll

	if err := g.sendRequest(http.MethodGet, url, nil, &poll); err != nil {
		return model.Poll{}, fmt.Errorf(
			"received an error after sending the request to get the poll: %w",
			err,
		)
	}

	return poll, nil
}

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

	if err := g.sendRequest(http.MethodPost, url, requestBody, nil); err != nil {
		return fmt.Errorf("received an error after sending the request to vote in the poll: %w", err)
	}

	return nil
}
