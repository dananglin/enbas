package gtsclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"codeflow.dananglin.me.uk/apollo/enbas/internal/model"
)

const baseStatusesPath string = "/api/v1/statuses"

func (g *GTSClient) GetStatus(statusID string, status *model.Status) error {
	params := requestParameters{
		httpMethod:  http.MethodGet,
		url:         g.auth.GetInstanceURL() + baseStatusesPath + "/" + statusID,
		requestBody: nil,
		contentType: "",
		output:      status,
	}

	if err := g.sendRequest(params); err != nil {
		return fmt.Errorf(
			"received an error after sending the request to get the status information: %w",
			err,
		)
	}

	return nil
}

type CreateStatusForm struct {
	Content       string                `json:"status"`
	InReplyTo     string                `json:"in_reply_to_id"`
	Language      string                `json:"language"`
	SpoilerText   string                `json:"spoiler_text"`
	Boostable     bool                  `json:"boostable"`
	LocalOnly     bool                  `json:"local_only"`
	Likeable      bool                  `json:"likeable"`
	Replyable     bool                  `json:"replyable"`
	Sensitive     bool                  `json:"sensitive"`
	Poll          *CreateStatusPollForm `json:"poll,omitempty"`
	ContentType   string                `json:"content_type"`
	Visibility    string                `json:"visibility"`
	AttachmentIDs []string              `json:"media_ids,omitempty"`
}

type CreateStatusPollForm struct {
	Options    []string `json:"options"`
	ExpiresIn  int      `json:"expires_in"`
	Multiple   bool     `json:"multiple"`
	HideTotals bool     `json:"hide_totals"`
}

func (g *GTSClient) CreateStatus(form CreateStatusForm, status *model.Status) error {
	data, err := json.Marshal(form)
	if err != nil {
		return fmt.Errorf("unable to create the JSON form: %w", err)
	}

	requestBody := bytes.NewBuffer(data)

	params := requestParameters{
		httpMethod:  http.MethodPost,
		url:         g.auth.GetInstanceURL() + baseStatusesPath,
		requestBody: requestBody,
		contentType: applicationJSON,
		output:      status,
	}

	if err := g.sendRequest(params); err != nil {
		return fmt.Errorf(
			"received an error after sending the request to create the status: %w",
			err,
		)
	}

	return nil
}

func (g *GTSClient) GetBookmarks(limit int, bookmarks *model.StatusList) error {
	path := fmt.Sprintf("/api/v1/bookmarks?limit=%d", limit)

	var statuses []model.Status

	params := requestParameters{
		httpMethod:  http.MethodGet,
		url:         g.auth.GetInstanceURL() + path,
		requestBody: nil,
		contentType: "",
		output:      &statuses,
	}

	if err := g.sendRequest(params); err != nil {
		return fmt.Errorf(
			"received an error after sending the request to get the bookmarks: %w",
			err,
		)
	}

	*bookmarks = model.StatusList{
		Name:     "Your Bookmarks",
		Statuses: statuses,
	}

	return nil
}

func (g *GTSClient) AddStatusToBookmarks(statusID string, _ *NoRPCResults) error {
	params := requestParameters{
		httpMethod:  http.MethodPost,
		url:         g.auth.GetInstanceURL() + fmt.Sprintf("/api/v1/statuses/%s/bookmark", statusID),
		requestBody: nil,
		contentType: "",
		output:      nil,
	}

	if err := g.sendRequest(params); err != nil {
		return fmt.Errorf(
			"received an error after sending the request to add the status to the list of bookmarks: %w",
			err,
		)
	}

	return nil
}

func (g *GTSClient) RemoveStatusFromBookmarks(statusID string, _ *NoRPCResults) error {
	params := requestParameters{
		httpMethod:  http.MethodPost,
		url:         g.auth.GetInstanceURL() + fmt.Sprintf("/api/v1/statuses/%s/unbookmark", statusID),
		requestBody: nil,
		contentType: "",
		output:      nil,
	}

	if err := g.sendRequest(params); err != nil {
		return fmt.Errorf(
			"received an error after sending the request to remove the status from the list of bookmarks: %w",
			err,
		)
	}

	return nil
}

func (g *GTSClient) LikeStatus(statusID string, _ *NoRPCResults) error {
	params := requestParameters{
		httpMethod:  http.MethodPost,
		url:         g.auth.GetInstanceURL() + baseStatusesPath + "/" + statusID + "/favourite",
		requestBody: nil,
		contentType: "",
		output:      nil,
	}

	if err := g.sendRequest(params); err != nil {
		return fmt.Errorf(
			"received an error after sending the request to like the status: %w",
			err,
		)
	}

	return nil
}

func (g *GTSClient) UnlikeStatus(statusID string, _ *NoRPCResults) error {
	params := requestParameters{
		httpMethod:  http.MethodPost,
		url:         g.auth.GetInstanceURL() + baseStatusesPath + "/" + statusID + "/unfavourite",
		requestBody: nil,
		contentType: "",
		output:      nil,
	}

	if err := g.sendRequest(params); err != nil {
		return fmt.Errorf(
			"received an error after sending the request to unlike the status: %w",
			err,
		)
	}

	return nil
}

func (g *GTSClient) GetFavourites(limit int, favourites *model.StatusList) error {
	var statuses []model.Status

	params := requestParameters{
		httpMethod:  http.MethodGet,
		url:         g.auth.GetInstanceURL() + fmt.Sprintf("/api/v1/favourites?limit=%d", limit),
		requestBody: nil,
		contentType: "",
		output:      &statuses,
	}

	if err := g.sendRequest(params); err != nil {
		return fmt.Errorf(
			"received an error after sending the request to get the list of statuses: %w",
			err,
		)
	}

	*favourites = model.StatusList{
		Name:     "Your favourite statuses",
		Statuses: statuses,
	}

	return nil
}

func (g *GTSClient) ReblogStatus(statusID string, _ *NoRPCResults) error {
	params := requestParameters{
		httpMethod:  http.MethodPost,
		url:         g.auth.GetInstanceURL() + baseStatusesPath + "/" + statusID + "/reblog",
		requestBody: nil,
		contentType: "",
		output:      nil,
	}

	if err := g.sendRequest(params); err != nil {
		return fmt.Errorf(
			"received an error after sending the request to reblog the status: %w",
			err,
		)
	}

	return nil
}

func (g *GTSClient) UnreblogStatus(statusID string, _ *NoRPCResults) error {
	params := requestParameters{
		httpMethod:  http.MethodPost,
		url:         g.auth.GetInstanceURL() + baseStatusesPath + "/" + statusID + "/unreblog",
		requestBody: nil,
		contentType: "",
		output:      nil,
	}

	if err := g.sendRequest(params); err != nil {
		return fmt.Errorf(
			"received an error after sending the request to un-reblog the status: %w",
			err,
		)
	}

	return nil
}

func (g *GTSClient) MuteStatus(statusID string, _ *NoRPCResults) error {
	params := requestParameters{
		httpMethod:  http.MethodPost,
		url:         g.auth.GetInstanceURL() + baseStatusesPath + "/" + statusID + "/mute",
		requestBody: nil,
		contentType: "",
		output:      nil,
	}

	if err := g.sendRequest(params); err != nil {
		return fmt.Errorf(
			"received an error after sending the request to mute the status: %w",
			err,
		)
	}

	return nil
}

func (g *GTSClient) UnmuteStatus(statusID string, _ *NoRPCResults) error {
	params := requestParameters{
		httpMethod:  http.MethodPost,
		url:         g.auth.GetInstanceURL() + baseStatusesPath + "/" + statusID + "/unmute",
		requestBody: nil,
		contentType: "",
		output:      nil,
	}

	if err := g.sendRequest(params); err != nil {
		return fmt.Errorf(
			"received an error after sending the request to unmute the status: %w",
			err,
		)
	}

	return nil
}

func (g *GTSClient) DeleteStatus(statusID string, text *string) error {
	var status model.Status

	params := requestParameters{
		httpMethod:  http.MethodDelete,
		url:         g.auth.GetInstanceURL() + baseStatusesPath + "/" + statusID,
		requestBody: nil,
		contentType: "",
		output:      &status,
	}

	if err := g.sendRequest(params); err != nil {
		return fmt.Errorf(
			"received an error after sending the request to delete the status: %w",
			err,
		)
	}

	*text = status.Text

	return nil
}

func (g *GTSClient) GetAccountsWhoLikedStatus(statusID string, list *model.AccountList) error {
	var accounts []model.Account

	params := requestParameters{
		httpMethod:  http.MethodGet,
		url:         g.auth.GetInstanceURL() + baseStatusesPath + "/" + statusID + "/favourited_by",
		requestBody: nil,
		contentType: "",
		output:      &accounts,
	}

	if err := g.sendRequest(params); err != nil {
		return fmt.Errorf(
			"recevied an error after sending the request to retrieve the accounts that liked the status: %w",
			err,
		)
	}

	*list = model.AccountList{
		Label:           "LIKED BY",
		Accounts:        accounts,
		BlockedAccounts: false,
	}

	return nil
}

func (g *GTSClient) GetAccountsWhoRebloggedStatus(statusID string, list *model.AccountList) error {
	var accounts []model.Account

	params := requestParameters{
		httpMethod:  http.MethodGet,
		url:         g.auth.GetInstanceURL() + baseStatusesPath + "/" + statusID + "/reblogged_by",
		requestBody: nil,
		contentType: "",
		output:      &accounts,
	}

	if err := g.sendRequest(params); err != nil {
		return fmt.Errorf(
			"received an error after sending the request to retrieve the accounts that reblogged the status: %w",
			err,
		)
	}

	*list = model.AccountList{
		Label:           "BOOSTED BY",
		Accounts:        accounts,
		BlockedAccounts: false,
	}

	return nil
}

func (g *GTSClient) GetThread(statusID string, thread *model.Thread) error {
	obj := struct {
		Ancestors   []model.Status `json:"ancestors"`
		Descendants []model.Status `json:"descendants"`
	}{
		Ancestors:   nil,
		Descendants: nil,
	}

	params := requestParameters{
		httpMethod:  http.MethodGet,
		url:         g.auth.GetInstanceURL() + baseStatusesPath + "/" + statusID + "/context",
		requestBody: nil,
		contentType: "",
		output:      &obj,
	}

	if err := g.sendRequest(params); err != nil {
		return fmt.Errorf(
			"received an error after sending the request to retrieve the status thread: %w",
			err,
		)
	}

	*thread = model.Thread{
		Ancestors: model.StatusList{
			Name:     "Ancestors",
			Statuses: obj.Ancestors,
		},
		Descendants: model.StatusList{
			Name:     "Descendants",
			Statuses: obj.Descendants,
		},
	}

	return nil
}
