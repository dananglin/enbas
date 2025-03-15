package executor

import (
	"fmt"
	"net/rpc"

	"codeflow.dananglin.me.uk/apollo/enbas/internal/gtsclient"
	"codeflow.dananglin.me.uk/apollo/enbas/internal/model"
	"codeflow.dananglin.me.uk/apollo/enbas/internal/printer"
	"codeflow.dananglin.me.uk/apollo/enbas/internal/server"
)

func (s *SearchExecutor) Execute() error {
	if s.resourceType == "" {
		return FlagNotSetError{flagText: flagTo}
	}

	if s.query == "" {
		return Error{"please enter a search query"}
	}

	funcMap := map[string]func(*rpc.Client) error{
		resourceTag:     s.searchTags,
		resourceAccount: s.searchAccounts,
		resourceStatus:  s.searchStatuses,
	}

	doFunc, ok := funcMap[s.resourceType]
	if !ok {
		return UnsupportedTypeError{resourceType: s.resourceType}
	}

	client, err := server.Connect(s.config.Server, s.configDir)
	if err != nil {
		return fmt.Errorf("error creating the client for the daemon process: %w", err)
	}
	defer client.Close()

	return doFunc(client)
}

func (s *SearchExecutor) searchTags(client *rpc.Client) error {
	var results model.TagList

	if err := client.Call(
		"GTSClient.SearchTags",
		gtsclient.SearchTagsArgs{
			Limit: s.limit,
			Query: s.query,
		},
		&results,
	); err != nil {
		return fmt.Errorf("error searching for tags: %w", err)
	}

	if err := printer.PrintTagList(s.printSettings, results); err != nil {
		return fmt.Errorf("error printing the search result: %w", err)
	}

	return nil
}

func (s *SearchExecutor) searchAccounts(client *rpc.Client) error {
	var results model.AccountList

	if err := client.Call(
		"GTSClient.SearchAccounts",
		gtsclient.SearchAccountsArgs{
			Limit:     s.limit,
			Query:     s.query,
			Resolve:   s.resolve,
			Following: s.following,
		},
		&results,
	); err != nil {
		return fmt.Errorf("error searching for accounts: %w", err)
	}

	if err := printer.PrintAccountList(s.printSettings, results); err != nil {
		return fmt.Errorf("error printing the search result: %w", err)
	}

	return nil
}

func (s *SearchExecutor) searchStatuses(client *rpc.Client) error {
	var (
		results model.StatusList
		err     error
	)

	accountID := ""
	if !s.accountName.Empty() {
		accountID, err = getAccountID(client, s.accountName)
		if err != nil {
			return fmt.Errorf("unable to get the account ID: %w", err)
		}
	}

	if err := client.Call(
		"GTSClient.SearchStatuses",
		gtsclient.SearchStatusesArgs{
			Limit:     s.limit,
			Query:     s.query,
			AccountID: accountID,
		},
		&results,
	); err != nil {
		return fmt.Errorf("error searching for statuses: %w", err)
	}

	var myAccountID string
	if err := client.Call("GTSClient.GetMyAccountID", gtsclient.NoRPCArgs{}, &myAccountID); err != nil {
		return fmt.Errorf("unable to get your account ID: %w", err)
	}

	if err := printer.PrintStatusList(s.printSettings, results, myAccountID); err != nil {
		return fmt.Errorf("error printing the search result: %w", err)
	}

	return nil
}
