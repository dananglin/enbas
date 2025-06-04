package gtsclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"slices"

	"codeflow.dananglin.me.uk/apollo/enbas/internal/model"
)

const (
	baseAccountsPath       = "/api/v1/accounts"
	baseFollowRequestsPath = "/api/v1/follow_requests"
)

func (g *GTSClient) GetMyAccount(_ NoRPCArgs, account *model.Account) error {
	var err error

	*account, err = g.verifyCredentials()
	if err != nil {
		return fmt.Errorf("error getting the account information: %w", err)
	}

	g.auth.UpdateCurrentAccountID(account.ID)

	return nil
}

func (g *GTSClient) GetMyAccountID(_ NoRPCArgs, accountID *string) error {
	currentAccountID := g.auth.GetCurrentAccountID()
	if currentAccountID != "" {
		*accountID = currentAccountID

		return nil
	}

	account, err := g.verifyCredentials()
	if err != nil {
		return fmt.Errorf("error getting the account information: %w", err)
	}

	*accountID = account.ID

	// Store the account ID in the GTSClient value before returning.
	g.auth.UpdateCurrentAccountID(account.ID)

	return nil
}

func (g *GTSClient) verifyCredentials() (model.Account, error) {
	var account model.Account

	params := requestParameters{
		httpMethod:  http.MethodGet,
		url:         g.auth.GetInstanceURL() + baseAccountsPath + "/verify_credentials",
		requestBody: nil,
		contentType: "",
		output:      &account,
	}

	if err := g.sendRequest(params); err != nil {
		return model.Account{}, fmt.Errorf(
			"received an error after sending the request to verify the credentials: %w",
			err,
		)
	}

	return account, nil
}

func (g *GTSClient) GetAccount(accountURI string, account *model.Account) error {
	var err error

	*account, err = g.getAccount(accountURI)
	if err != nil {
		return fmt.Errorf("error getting the account information: %w", err)
	}

	return nil
}

func (g *GTSClient) GetMultipleAccounts(accountURIs []string, accounts *[]model.Account) error {
	output := make([]model.Account, len(accountURIs))

	for idx, uri := range slices.All(accountURIs) {
		acct, err := g.getAccount(uri)
		if err != nil {
			return fmt.Errorf("error retrieving the account details for %q: %w", uri, err)
		}

		output[idx] = acct
	}

	*accounts = output

	return nil
}

func (g *GTSClient) GetAccountID(accountURI string, accountID *string) error {
	account, err := g.getAccount(accountURI)
	if err != nil {
		return fmt.Errorf("error getting the account information: %w", err)
	}

	*accountID = account.ID

	return nil
}

func (g *GTSClient) GetMultipleAccountIDs(accountURIs []string, accountIDs *[]string) error {
	ids := make([]string, len(accountURIs))

	for idx, uri := range slices.All(accountURIs) {
		acct, err := g.getAccount(uri)
		if err != nil {
			return fmt.Errorf("error retrieving the account ID of %q: %w", uri, err)
		}

		ids[idx] = acct.ID
	}

	*accountIDs = ids

	return nil
}

func (g *GTSClient) getAccount(accountURI string) (model.Account, error) {
	var account model.Account

	params := requestParameters{
		httpMethod:  http.MethodGet,
		url:         g.auth.GetInstanceURL() + baseAccountsPath + "/lookup?acct=" + accountURI,
		requestBody: nil,
		contentType: "",
		output:      &account,
	}

	if err := g.sendRequest(params); err != nil {
		return model.Account{}, fmt.Errorf(
			"received an error after sending the request to get the account information: %w",
			err,
		)
	}

	return account, nil
}

func (g *GTSClient) GetAccountRelationship(accountID string, relationship *model.AccountRelationship) error {
	var relationships []model.AccountRelationship

	params := requestParameters{
		httpMethod:  http.MethodGet,
		url:         g.auth.GetInstanceURL() + baseAccountsPath + "/relationships?id=" + accountID,
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
		url:         g.auth.GetInstanceURL() + baseAccountsPath + "/" + form.AccountID + "/follow",
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
		url:         g.auth.GetInstanceURL() + baseAccountsPath + "/" + accountID + "/unfollow",
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
		url:         g.auth.GetInstanceURL() + fmt.Sprintf("%s/%s/followers?limit=%d", baseAccountsPath, args.AccountID, args.Limit),
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
		Label:           "Followed by",
		Accounts:        accounts,
		BlockedAccounts: false,
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
		url:         g.auth.GetInstanceURL() + fmt.Sprintf("%s/%s/following?limit=%d", baseAccountsPath, args.AccountID, args.Limit),
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
		Label:           "Following",
		Accounts:        accounts,
		BlockedAccounts: false,
	}

	return nil
}

func (g *GTSClient) BlockAccount(accountID string, _ *NoRPCResults) error {
	params := requestParameters{
		httpMethod:  http.MethodPost,
		url:         g.auth.GetInstanceURL() + baseAccountsPath + "/" + accountID + "/block",
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
		url:         g.auth.GetInstanceURL() + baseAccountsPath + "/" + accountID + "/unblock",
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
		url:         g.auth.GetInstanceURL() + fmt.Sprintf("/api/v1/blocks?limit=%d", limit),
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
		Label:           "Blocked accounts",
		Accounts:        accounts,
		BlockedAccounts: true,
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
		url:         g.auth.GetInstanceURL() + baseAccountsPath + "/" + args.AccountID + "/note",
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
		url:         g.auth.GetInstanceURL() + fmt.Sprintf("%s?limit=%d", baseFollowRequestsPath, limit),
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
		Label:           "Accounts that have requested to follow you",
		Accounts:        accounts,
		BlockedAccounts: false,
	}

	return nil
}

func (g *GTSClient) AcceptFollowRequest(accountID string, _ *NoRPCResults) error {
	params := requestParameters{
		httpMethod:  http.MethodPost,
		url:         g.auth.GetInstanceURL() + baseFollowRequestsPath + "/" + accountID + "/authorize",
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
		url:         g.auth.GetInstanceURL() + baseFollowRequestsPath + "/" + accountID + "/reject",
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
		url:         g.auth.GetInstanceURL() + fmt.Sprintf("/api/v1/mutes?limit=%d", limit),
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
		Label:           "Muted accounts",
		Accounts:        accounts,
		BlockedAccounts: false,
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
		url:         g.auth.GetInstanceURL() + baseAccountsPath + "/" + args.AccountID + "/mute",
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
		url:         g.auth.GetInstanceURL() + baseAccountsPath + "/" + accountID + "/unmute",
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
		url:         g.auth.GetInstanceURL() + path + query,
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
