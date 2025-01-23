package gtsclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

const pollPath string = "/api/v1/polls"

type VoteInPollArgs struct {
	PollID  string
	Choices []int
}

func (g *GTSClient) VoteInPoll(args VoteInPollArgs, _ *NoRPCResults) error {
	form := struct {
		Choices []int `json:"choices"`
	}{
		Choices: args.Choices,
	}

	data, err := json.Marshal(form)
	if err != nil {
		return fmt.Errorf("unable to encode the JSON form: %w", err)
	}

	requestBody := bytes.NewBuffer(data)

	params := requestParameters{
		httpMethod:  http.MethodPost,
		url:         g.Authentication.Instance + pollPath + "/" + args.PollID + "/votes",
		requestBody: requestBody,
		contentType: applicationJSON,
		output:      nil,
	}

	if err := g.sendRequest(params); err != nil {
		return fmt.Errorf("received an error after sending the request to vote in the poll: %w", err)
	}

	return nil
}
