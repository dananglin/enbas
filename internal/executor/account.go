package executor

import (
	"fmt"
	"net/rpc"

	internalFlag "codeflow.dananglin.me.uk/apollo/enbas/internal/flag"
	"codeflow.dananglin.me.uk/apollo/enbas/internal/gtsclient"
	"codeflow.dananglin.me.uk/apollo/enbas/internal/model"
)

func getAccountID(client *rpc.Client, accountNames internalFlag.StringSliceValue) (string, error) {
	account, err := getAccount(client, accountNames)
	if err != nil {
		return "", fmt.Errorf("unable to get the account information: %w", err)
	}

	return account.ID, nil
}

func getAccount(client *rpc.Client, accountNames internalFlag.StringSliceValue) (model.Account, error) {
	if accountNames.Empty() {
		return model.Account{}, NoAccountSpecifiedError{}
	}

	if !accountNames.ExpectedLength(1) {
		return model.Account{}, fmt.Errorf(
			"received an unexpected number of account names: want 1, got %d",
			len(accountNames),
		)
	}

	var account model.Account
	if err := client.Call("GTSClient.GetAccount", accountNames[0], &account); err != nil {
		return model.Account{}, fmt.Errorf("unable to retrieve the account details: %w", err)
	}

	return account, nil
}

func getMultipleAccounts(client *rpc.Client, accountNames internalFlag.StringSliceValue) ([]model.Account, error) {
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
