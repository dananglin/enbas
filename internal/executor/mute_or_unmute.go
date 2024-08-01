package executor

import (
	"flag"
	"fmt"
	"time"

	"codeflow.dananglin.me.uk/apollo/enbas/internal/client"
	"codeflow.dananglin.me.uk/apollo/enbas/internal/config"
	"codeflow.dananglin.me.uk/apollo/enbas/internal/printer"
)

type MuteOrUnmuteExecutor struct {
	*flag.FlagSet

	printer           *printer.Printer
	config            *config.Config
	accountName       string
	command           string
	resourceType      string
	muteDuration      TimeDurationFlagValue
	muteNotifications bool
}

func NewMuteOrUnmuteExecutor(printer *printer.Printer, config *config.Config, name, summary string) *MuteOrUnmuteExecutor {
	muteDuration := TimeDurationFlagValue{time.Duration(0 * time.Second)}

	exe := MuteOrUnmuteExecutor{
		FlagSet: flag.NewFlagSet(name, flag.ExitOnError),

		printer:      printer,
		config:       config,
		command:      name,
		muteDuration: muteDuration,
	}

	exe.StringVar(&exe.accountName, flagAccountName, "", "Specify the account name in full (username@domain)")
	exe.StringVar(&exe.resourceType, flagType, "", "Specify the type of resource to mute or unmute")
	exe.BoolVar(&exe.muteNotifications, flagMuteNotifications, false, "Mute notifications as well as posts")
	exe.Var(&exe.muteDuration, flagMuteDuration, "Specify how long the mute should last for. To mute indefinitely, set this to 0s")

	exe.Usage = commandUsageFunc(name, summary, exe.FlagSet)

	return &exe
}

func (m *MuteOrUnmuteExecutor) Execute() error {
	funcMap := map[string]func(*client.Client) error{
		resourceAccount: m.muteOrUnmuteAccount,
	}

	doFunc, ok := funcMap[m.resourceType]
	if !ok {
		return UnsupportedTypeError{resourceType: m.resourceType}
	}

	gtsClient, err := client.NewClientFromFile(m.config.CredentialsFile)
	if err != nil {
		return fmt.Errorf("unable to create the GoToSocial client: %w", err)
	}

	return doFunc(gtsClient)
}

func (m *MuteOrUnmuteExecutor) muteOrUnmuteAccount(gtsClient *client.Client) error {
	if m.accountName == "" {
		return FlagNotSetError{flagText: flagAccountName}
	}

	accountID, err := getAccountID(gtsClient, false, m.accountName, m.config.CredentialsFile)
	if err != nil {
		return fmt.Errorf("received an error while getting the account ID: %w", err)
	}

	switch m.command {
	case CommandMute:
		return m.muteAccount(gtsClient, accountID)
	case CommandUnmute:
		return m.unmuteAccount(gtsClient, accountID)
	default:
		return nil
	}
}

func (m *MuteOrUnmuteExecutor) muteAccount(gtsClient *client.Client, accountID string) error {
	form := client.MuteAccountForm{
		Notifications: m.muteNotifications,
		Duration:      int(m.muteDuration.Duration.Seconds()),
	}

	if err := gtsClient.MuteAccount(accountID, form); err != nil {
		return fmt.Errorf("unable to mute the account: %w", err)
	}

	m.printer.PrintSuccess("Successfully muted the account.")

	return nil
}

func (m *MuteOrUnmuteExecutor) unmuteAccount(gtsClient *client.Client, accountID string) error {
	if err := gtsClient.UnmuteAccount(accountID); err != nil {
		return fmt.Errorf("unable to unmute the account: %w", err)
	}

	m.printer.PrintSuccess("Successfully unmuted the account.")

	return nil
}
