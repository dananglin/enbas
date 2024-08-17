package executor

import (
	"fmt"

	"codeflow.dananglin.me.uk/apollo/enbas/internal/client"
	internalFlag "codeflow.dananglin.me.uk/apollo/enbas/internal/flag"
	"codeflow.dananglin.me.uk/apollo/enbas/internal/model"
)

func getAccountID(
	gtsClient *client.Client,
	myAccount bool,
	accountNames internalFlag.StringSliceValue,
) (string, error) {
	account, err := getAccount(gtsClient, myAccount, accountNames)
	if err != nil {
		return "", fmt.Errorf("unable to get the account information: %w", err)
	}

	return account.ID, nil
}

func getAccount(
	gtsClient *client.Client,
	myAccount bool,
	accountNames internalFlag.StringSliceValue,
) (model.Account, error) {
	var (
		account model.Account
		err     error
	)

	switch {
	case myAccount:
		account, err = getMyAccount(gtsClient)
		if err != nil {
			return account, fmt.Errorf("unable to get your account ID: %w", err)
		}
	case !accountNames.Empty():
		account, err = getOtherAccount(gtsClient, accountNames)
		if err != nil {
			return account, fmt.Errorf("unable to get the account ID: %w", err)
		}
	default:
		return account, NoAccountSpecifiedError{}
	}

	return account, nil
}

func getMyAccount(gtsClient *client.Client) (model.Account, error) {
	account, err := gtsClient.VerifyCredentials()
	if err != nil {
		return model.Account{}, fmt.Errorf("unable to retrieve your account: %w", err)
	}

	return account, nil
}

func getOtherAccount(gtsClient *client.Client, accountNames internalFlag.StringSliceValue) (model.Account, error) {
	expectedNumAccountNames := 1
	if !accountNames.ExpectedLength(expectedNumAccountNames) {
		return model.Account{}, fmt.Errorf(
			"received an unexpected number of account names: want %d",
			expectedNumAccountNames,
		)
	}

	account, err := gtsClient.GetAccount(accountNames[0])
	if err != nil {
		return model.Account{}, fmt.Errorf("unable to retrieve the account details: %w", err)
	}

	return account, nil
}

func getOtherAccounts(gtsClient *client.Client, accountNames internalFlag.StringSliceValue) ([]model.Account, error) {
	numAccountNames := len(accountNames)
	accounts := make([]model.Account, numAccountNames)

	for ind := range numAccountNames {
		var err error

		accounts[ind], err = gtsClient.GetAccount(accountNames[ind])
		if err != nil {
			return nil, fmt.Errorf("unable to retrieve the account information for %s: %w", accountNames[ind], err)
		}
	}

	return accounts, nil
}
