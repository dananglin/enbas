package executor

import (
	"fmt"
	"net/rpc"
	"strings"

	"codeflow.dananglin.me.uk/apollo/enbas/internal/gtsclient"
	"codeflow.dananglin.me.uk/apollo/enbas/internal/server"
)

func (f *FollowExecutor) Execute() error {
	funcMap := map[string]func(*rpc.Client) error{
		resourceAccount: f.followAccount,
		resourceTag:     f.followTag,
	}

	doFunc, ok := funcMap[f.resourceType]
	if !ok {
		return UnsupportedTypeError{resourceType: f.resourceType}
	}

	client, err := server.Connect(f.config.Server, f.configDir)
	if err != nil {
		return fmt.Errorf("error creating the client for the daemon process: %w", err)
	}
	defer client.Close()

	return doFunc(client)
}

func (f *FollowExecutor) followAccount(client *rpc.Client) error {
	accountID, err := getAccountID(client, false, f.accountName)
	if err != nil {
		return fmt.Errorf("received an error while getting the account ID: %w", err)
	}

	if err := client.Call(
		"GTSClient.FollowAccount",
		gtsclient.FollowAccountArgs{
			AccountID:   accountID,
			ShowReposts: f.showReposts,
			Notify:      f.notify,
		},
		nil,
	); err != nil {
		return fmt.Errorf("error following the account: %w", err)
	}

	f.printer.PrintSuccess("Successfully sent the follow request.")

	return nil
}

func (f *FollowExecutor) followTag(client *rpc.Client) error {
	if f.tag == "" {
		return Error{"please provide the name of the tag"}
	}

	tag := strings.TrimLeft(f.tag, "#")

	if err := client.Call("GTSClient.FollowTag", tag, nil); err != nil {
		return fmt.Errorf("error following the tag: %w", err)
	}

	f.printer.PrintSuccess("You are now following '" + tag + "'.")

	return nil
}
