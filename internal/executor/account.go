package executor

import (
	"fmt"

	"codeflow.dananglin.me.uk/apollo/enbas/internal/client"
	"codeflow.dananglin.me.uk/apollo/enbas/internal/config"
	"codeflow.dananglin.me.uk/apollo/enbas/internal/model"
)

func getAccountID(gtsClient *client.Client, myAccount bool, accountName, path string) (string, error) {
	var (
		accountID string
		err       error
	)

	switch {
	case myAccount:
		accountID, err = getMyAccountID(gtsClient, path)
		if err != nil {
			return "", fmt.Errorf("unable to get your account ID: %w", err)
		}
	case accountName != "":
		accountID, err = getTheirAccountID(gtsClient, accountName)
		if err != nil {
			return "", fmt.Errorf("unable to get their account ID: %w", err)
		}
	default:
		return "", NoAccountSpecifiedError{}
	}

	return accountID, nil
}

func getTheirAccountID(gtsClient *client.Client, accountURI string) (string, error) {
	account, err := getAccount(gtsClient, accountURI)
	if err != nil {
		return "", fmt.Errorf("unable to retrieve your account: %w", err)
	}

	return account.ID, nil
}

func getMyAccountID(gtsClient *client.Client, path string) (string, error) {
	account, err := getMyAccount(gtsClient, path)
	if err != nil {
		return "", fmt.Errorf("received an error while getting your account details: %w", err)
	}

	return account.ID, nil
}

func getMyAccount(gtsClient *client.Client, path string) (model.Account, error) {
	authConfig, err := config.NewCredentialsConfigFromFile(path)
	if err != nil {
		return model.Account{}, fmt.Errorf("unable to retrieve the authentication configuration: %w", err)
	}

	accountURI := authConfig.CurrentAccount

	account, err := getAccount(gtsClient, accountURI)
	if err != nil {
		return model.Account{}, fmt.Errorf("unable to retrieve your account: %w", err)
	}

	return account, nil
}

func getAccount(gtsClient *client.Client, accountURI string) (model.Account, error) {
	account, err := gtsClient.GetAccount(accountURI)
	if err != nil {
		return model.Account{}, fmt.Errorf("unable to retrieve the account details: %w", err)
	}

	return account, nil
}
