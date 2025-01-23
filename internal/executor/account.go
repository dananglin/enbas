package executor

import (
	"fmt"
	"net/rpc"

	internalFlag "codeflow.dananglin.me.uk/apollo/enbas/internal/flag"
	"codeflow.dananglin.me.uk/apollo/enbas/internal/gtsclient"
	"codeflow.dananglin.me.uk/apollo/enbas/internal/model"
)

func getAccountID(
	client *rpc.Client,
	myAccount bool,
	accountNames internalFlag.StringSliceValue,
) (string, error) {
	account, err := getAccount(client, myAccount, accountNames)
	if err != nil {
		return "", fmt.Errorf("unable to get the account information: %w", err)
	}

	return account.ID, nil
}

func getAccount(
	client *rpc.Client,
	myAccount bool,
	accountNames internalFlag.StringSliceValue,
) (model.Account, error) {
	var (
		account model.Account
		err     error
	)

	switch {
	case myAccount:
		account, err = getMyAccount(client)
		if err != nil {
			return account, fmt.Errorf("unable to get your account ID: %w", err)
		}
	case !accountNames.Empty():
		account, err = getOtherAccount(client, accountNames)
		if err != nil {
			return account, fmt.Errorf("unable to get the account ID: %w", err)
		}
	default:
		return account, NoAccountSpecifiedError{}
	}

	return account, nil
}

func getMyAccount(client *rpc.Client) (model.Account, error) {
	var account model.Account
	if err := client.Call("GTSClient.VerifyCredentials", gtsclient.NoRPCArgs{}, &account); err != nil {
		return model.Account{}, fmt.Errorf("unable to retrieve your account: %w", err)
	}

	return account, nil
}

func getOtherAccount(client *rpc.Client, accountNames internalFlag.StringSliceValue) (model.Account, error) {
	expectedNumAccountNames := 1
	if !accountNames.ExpectedLength(expectedNumAccountNames) {
		return model.Account{}, fmt.Errorf(
			"received an unexpected number of account names: want %d",
			expectedNumAccountNames,
		)
	}

	var account model.Account
	if err := client.Call("GTSClient.GetAccount", accountNames[0], &account); err != nil {
		return model.Account{}, fmt.Errorf("unable to retrieve the account details: %w", err)
	}

	return account, nil
}

func getOtherAccounts(client *rpc.Client, accountNames internalFlag.StringSliceValue) ([]model.Account, error) {
	numAccountNames := len(accountNames)
	accounts := make([]model.Account, numAccountNames)

	for ind := range numAccountNames {
		var account model.Account

		if err := client.Call("GTSClient.GetAccount", accountNames[ind], &account); err != nil {
			return nil, fmt.Errorf(
				"unable to retrieve the account information for %s: %w",
				accountNames[ind],
				err,
			)
		}

		accounts[ind] = account
	}

	return accounts, nil
}

func getAccountsFromList(client *rpc.Client, listID string) (map[string]string, error) {
	var acctList model.AccountList
	if err := client.Call(
		"GTSClient.GetAccountsFromList",
		gtsclient.GetAccountsFromListArgs{
			ListID: listID,
			Limit:  0,
		},
		&acctList,
	); err != nil {
		return map[string]string{}, fmt.Errorf("unable to retrieve the accounts from the list: %w", err)
	}

	if len(acctList.Accounts) == 0 {
		return map[string]string{}, nil
	}

	acctMap := make(map[string]string)

	for i := range acctList.Accounts {
		acctMap[acctList.Accounts[i].Acct] = acctList.Accounts[i].Username
	}

	return acctMap, nil
}
