package executor

import (
	"fmt"
	"net/rpc"

	"codeflow.dananglin.me.uk/apollo/enbas/internal/gtsclient"
	"codeflow.dananglin.me.uk/apollo/enbas/internal/model"
)

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
