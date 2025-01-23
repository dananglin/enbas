package gtsclient

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

func (g *GTSClient) VerifyCredentials(_ NoRPCArgs, account *model.Account) error {
	params := requestParameters{
		httpMethod:  http.MethodGet,
		url:         g.authentication.Instance + baseAccountsPath + "/verify_credentials",
		requestBody: nil,
		contentType: "",
		output:      account,
	}

	if err := g.sendRequest(params); err != nil {
		return fmt.Errorf(
			"received an error after sending the request to verify the credentials: %w",
			err,
		)
	}

	return nil
}

func (g *GTSClient) GetAccount(accountURI string, account *model.Account) error {
	params := requestParameters{
		httpMethod:  http.MethodGet,
		url:         g.authentication.Instance + baseAccountsPath + "/lookup?acct=" + accountURI,
		requestBody: nil,
		contentType: "",
		output:      account,
	}

	if err := g.sendRequest(params); err != nil {
		return fmt.Errorf(
			"received an error after sending the request to get the account information: %w",
			err,
		)
	}

	return nil
}

func (g *GTSClient) GetAccountRelationship(accountID string, relationship *model.AccountRelationship) error {
	var relationships []model.AccountRelationship

	params := requestParameters{
		httpMethod:  http.MethodGet,
		url:         g.authentication.Instance + baseAccountsPath + "/relationships?id=" + accountID,
		requestBody: nil,
		contentType: "",
		output:      &relationships,
	}

	if err := g.sendRequest(params); err != nil {
		return fmt.Errorf(
			"received an error after sending the request to get the account relationship: %w",
			err,
		)
	}

	if len(relationships) != 1 {
		return fmt.Errorf(
			"unexpected number of account relationships returned: want 1, got %d",
			len(relationships),
		)
	}

	*relationship = relationships[0]

	return nil
}

type FollowAccountArgs struct {
	AccountID   string
	ShowReposts bool
	Notify      bool
}

func (g *GTSClient) FollowAccount(args FollowAccountArgs, _ *NoRPCResults) error {
	form := struct {
		AccountID   string `json:"id"`
		ShowReposts bool   `json:"reblogs"`
		Notify      bool   `json:"notify"`
	}{
		AccountID:   args.AccountID,
		ShowReposts: args.ShowReposts,
		Notify:      args.Notify,
	}

	data, err := json.Marshal(form)
	if err != nil {
		return fmt.Errorf("unable to marshal the form: %w", err)
	}

	requestBody := bytes.NewBuffer(data)

	params := requestParameters{
		httpMethod:  http.MethodPost,
		url:         g.authentication.Instance + baseAccountsPath + "/" + form.AccountID + "/follow",
		requestBody: requestBody,
		contentType: applicationJSON,
		output:      nil,
	}

	if err := g.sendRequest(params); err != nil {
		return fmt.Errorf("received an error after sending the follow request: %w", err)
	}

	return nil
}

func (g *GTSClient) UnfollowAccount(accountID string, _ *NoRPCResults) error {
	params := requestParameters{
		httpMethod:  http.MethodPost,
		url:         g.authentication.Instance + baseAccountsPath + "/" + accountID + "/unfollow",
		requestBody: nil,
		contentType: "",
		output:      nil,
	}

	if err := g.sendRequest(params); err != nil {
		return fmt.Errorf("received an error after sending the request to unfollow the account: %w", err)
	}

	return nil
}

type GetFollowersArgs struct {
	AccountID string
	Limit     int
}

func (g *GTSClient) GetFollowers(args GetFollowersArgs, followers *model.AccountList) error {
	var accounts []model.Account

	params := requestParameters{
		httpMethod:  http.MethodGet,
		url:         g.authentication.Instance + fmt.Sprintf("%s/%s/followers?limit=%d", baseAccountsPath, args.AccountID, args.Limit),
		requestBody: nil,
		contentType: "",
		output:      &accounts,
	}

	if err := g.sendRequest(params); err != nil {
		return fmt.Errorf(
			"received an error after sending the request to get the list of followers: %w",
			err,
		)
	}

	*followers = model.AccountList{
		Type:     model.AccountListFollowers,
		Accounts: accounts,
	}

	return nil
}

type GetFollowingsArgs struct {
	AccountID string
	Limit     int
}

func (g *GTSClient) GetFollowing(args GetFollowingsArgs, following *model.AccountList) error {
	var accounts []model.Account

	params := requestParameters{
		httpMethod:  http.MethodGet,
		url:         g.authentication.Instance + fmt.Sprintf("%s/%s/following?limit=%d", baseAccountsPath, args.AccountID, args.Limit),
		requestBody: nil,
		contentType: "",
		output:      &accounts,
	}

	if err := g.sendRequest(params); err != nil {
		return fmt.Errorf(
			"received an error after sending the request to get the list of followed accounts: %w",
			err,
		)
	}

	*following = model.AccountList{
		Type:     model.AccountListFollowing,
		Accounts: accounts,
	}

	return nil
}

func (g *GTSClient) BlockAccount(accountID string, _ *NoRPCResults) error {
	params := requestParameters{
		httpMethod:  http.MethodPost,
		url:         g.authentication.Instance + baseAccountsPath + "/" + accountID + "/block",
		requestBody: nil,
		contentType: "",
		output:      nil,
	}

	if err := g.sendRequest(params); err != nil {
		return fmt.Errorf("received an error after sending the request to block the account: %w", err)
	}

	return nil
}

func (g *GTSClient) UnblockAccount(accountID string, _ *NoRPCResults) error {
	params := requestParameters{
		httpMethod:  http.MethodPost,
		url:         g.authentication.Instance + baseAccountsPath + "/" + accountID + "/unblock",
		requestBody: nil,
		contentType: "",
		output:      nil,
	}

	if err := g.sendRequest(params); err != nil {
		return fmt.Errorf("received an error after sending the request to unblock the account: %w", err)
	}

	return nil
}

func (g *GTSClient) GetBlockedAccounts(limit int, blocked *model.AccountList) error {
	var accounts []model.Account

	params := requestParameters{
		httpMethod:  http.MethodGet,
		url:         g.authentication.Instance + fmt.Sprintf("/api/v1/blocks?limit=%d", limit),
		requestBody: nil,
		contentType: "",
		output:      &accounts,
	}

	if err := g.sendRequest(params); err != nil {
		return fmt.Errorf(
			"received an error after sending the request to get the list of blocked accounts: %w",
			err,
		)
	}

	*blocked = model.AccountList{
		Type:     model.AccountListBlockedAccount,
		Accounts: accounts,
	}

	return nil
}

type SetPrivateNoteArgs struct {
	AccountID string
	Note      string
}

func (g *GTSClient) SetPrivateNote(args SetPrivateNoteArgs, _ *NoRPCResults) error {
	form := struct {
		Comment string `json:"comment"`
	}{
		Comment: args.Note,
	}

	data, err := json.Marshal(form)
	if err != nil {
		return fmt.Errorf("unable to marshal the form: %w", err)
	}

	requestBody := bytes.NewBuffer(data)

	params := requestParameters{
		httpMethod:  http.MethodPost,
		url:         g.authentication.Instance + baseAccountsPath + "/" + args.AccountID + "/note",
		requestBody: requestBody,
		contentType: applicationJSON,
		output:      nil,
	}

	if err := g.sendRequest(params); err != nil {
		return fmt.Errorf("received an error after sending the request to set the private note: %w", err)
	}

	return nil
}

func (g *GTSClient) GetFollowRequests(limit int, requests *model.AccountList) error {
	var accounts []model.Account

	params := requestParameters{
		httpMethod:  http.MethodGet,
		url:         g.authentication.Instance + fmt.Sprintf("%s?limit=%d", baseFollowRequestsPath, limit),
		requestBody: nil,
		contentType: "",
		output:      &accounts,
	}

	if err := g.sendRequest(params); err != nil {
		return fmt.Errorf(
			"received an error after sending the request to get the list of follow requests: %w",
			err,
		)
	}

	*requests = model.AccountList{
		Type:     model.AccountListFollowRequests,
		Accounts: accounts,
	}

	return nil
}

func (g *GTSClient) AcceptFollowRequest(accountID string, _ *NoRPCResults) error {
	params := requestParameters{
		httpMethod:  http.MethodPost,
		url:         g.authentication.Instance + baseFollowRequestsPath + "/" + accountID + "/authorize",
		requestBody: nil,
		contentType: "",
		output:      nil,
	}

	if err := g.sendRequest(params); err != nil {
		return fmt.Errorf("received an error after sending the request to accept the follow request: %w", err)
	}

	return nil
}

func (g *GTSClient) RejectFollowRequest(accountID string, _ *NoRPCResults) error {
	params := requestParameters{
		httpMethod:  http.MethodPost,
		url:         g.authentication.Instance + baseFollowRequestsPath + "/" + accountID + "/reject",
		requestBody: nil,
		contentType: "",
		output:      nil,
	}

	if err := g.sendRequest(params); err != nil {
		return fmt.Errorf("received an error after sending the request to reject the follow request: %w", err)
	}

	return nil
}

func (g *GTSClient) GetMutedAccounts(limit int, muted *model.AccountList) error {
	var accounts []model.Account

	params := requestParameters{
		httpMethod:  http.MethodGet,
		url:         g.authentication.Instance + fmt.Sprintf("/api/v1/mutes?limit=%d", limit),
		requestBody: nil,
		contentType: "",
		output:      &accounts,
	}

	if err := g.sendRequest(params); err != nil {
		return fmt.Errorf(
			"received an error after sending the request to get the list of muted accounts: %w",
			err,
		)
	}

	*muted = model.AccountList{
		Type:     model.AccountListMuted,
		Accounts: accounts,
	}

	return nil
}

type MuteAccountArgs struct {
	AccountID     string
	Notifications bool
	Duration      int
}

func (g *GTSClient) MuteAccount(args MuteAccountArgs, _ *NoRPCResults) error {
	form := struct {
		Notifications bool `json:"notifications"`
		Duration      int  `json:"duration"`
	}{
		Notifications: args.Notifications,
		Duration:      args.Duration,
	}

	data, err := json.Marshal(form)
	if err != nil {
		return fmt.Errorf("unable to marshal the form: %w", err)
	}

	requestBody := bytes.NewBuffer(data)

	params := requestParameters{
		httpMethod:  http.MethodPost,
		url:         g.authentication.Instance + baseAccountsPath + "/" + args.AccountID + "/mute",
		requestBody: requestBody,
		contentType: applicationJSON,
		output:      nil,
	}

	if err := g.sendRequest(params); err != nil {
		return fmt.Errorf("received an error after sending the request to mute the account: %w", err)
	}

	return nil
}

func (g *GTSClient) UnmuteAccount(accountID string, _ *NoRPCResults) error {
	params := requestParameters{
		httpMethod:  http.MethodPost,
		url:         g.authentication.Instance + baseAccountsPath + "/" + accountID + "/unmute",
		requestBody: nil,
		contentType: "",
		output:      nil,
	}

	if err := g.sendRequest(params); err != nil {
		return fmt.Errorf("received an error after sending the request to unmute the account: %w", err)
	}

	return nil
}

type GetAccountStatusesArgs struct {
	AccountID      string
	Limit          int
	ExcludeReplies bool
	ExcludeReblogs bool
	Pinned         bool
	OnlyMedia      bool
	OnlyPublic     bool
}

func (g *GTSClient) GetAccountStatuses(args GetAccountStatusesArgs, statusList *model.StatusList) error {
	path := baseAccountsPath + "/" + args.AccountID + "/statuses"
	query := fmt.Sprintf(
		"?limit=%d&exclude_replies=%t&exclude_reblogs=%t&pinned=%t&only_media=%t&only_public=%t",
		args.Limit,
		args.ExcludeReplies,
		args.ExcludeReblogs,
		args.Pinned,
		args.OnlyMedia,
		args.OnlyPublic,
	)

	var statuses []model.Status

	params := requestParameters{
		httpMethod:  http.MethodGet,
		url:         g.authentication.Instance + path + query,
		requestBody: nil,
		contentType: "",
		output:      &statuses,
	}

	if err := g.sendRequest(params); err != nil {
		return fmt.Errorf(
			"received an error after sending the request to get the account's statuses: %w",
			err,
		)
	}

	*statusList = model.StatusList{
		Name:     "STATUSES:",
		Statuses: statuses,
	}

	return nil
}
