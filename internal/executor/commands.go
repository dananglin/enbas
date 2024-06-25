// SPDX-FileCopyrightText: 2024 Dan Anglin <d.n.i.anglin@gmail.com>
//
// SPDX-License-Identifier: GPL-3.0-or-later

package executor

const (
	CommandAccept   string = "accept"
	CommandAdd      string = "add"
	CommandBlock    string = "block"
	CommandCreate   string = "create"
	CommandDelete   string = "delete"
	CommandEdit     string = "edit"
	CommandFollow   string = "follow"
	CommandInit     string = "init"
	CommandLogin    string = "login"
	CommandMute     string = "mute"
	CommandReject   string = "reject"
	CommandRemove   string = "remove"
	CommandShow     string = "show"
	CommandSwitch   string = "switch"
	CommandUnblock  string = "unblock"
	CommandUnfollow string = "unfollow"
	CommandUnmute   string = "unmute"
	CommandVersion  string = "version"
	CommandWhoami   string = "whoami"

	commandAcceptSummary   string = "Accept a request (e.g. a follow request)"
	commandAddSummary      string = "Add a resource to another resource"
	commandBlockSummary    string = "Block a resource (e.g. an account)"
	commandCreateSummary   string = "Create a specific resource"
	commandDeleteSummary   string = "Delete a specific resource"
	commandEditSummary     string = "Edit a specific resource"
	commandFollowSummary   string = "Follow a resource (e.g. an account)"
	commandInitSummary     string = "Create a new configuration file in the specified configuration directory"
	commandLoginSummary    string = "Login to an account on GoToSocial"
	commandMuteSummary     string = "Mute a resource (e.g. an account)"
	commandRejectSummary   string = "Reject a request (e.g. a follow request)"
	commandRemoveSummary   string = "Remove a resource from another resource"
	commandShowSummary     string = "Print details about a specified resource"
	commandSwitchSummary   string = "Perform a switch operation (e.g. switch logged in accounts)"
	commandUnblockSummary  string = "Unblock a resource (e.g. an account)"
	commandUnfollowSummary string = "Unfollow a resource (e.g. an account)"
	commandUnmuteSummary   string = "Unmute a resource (e.g. an account)"
	commandVersionSummary  string = "Print the application's version and build information"
	commandWhoamiSummary   string = "Print the account that you are currently logged in to"
)

func CommandSummaryMap() map[string]string {
	return map[string]string{
		CommandAccept:   commandAcceptSummary,
		CommandAdd:      commandAddSummary,
		CommandBlock:    commandBlockSummary,
		CommandCreate:   commandCreateSummary,
		CommandDelete:   commandDeleteSummary,
		CommandEdit:     commandEditSummary,
		CommandFollow:   commandFollowSummary,
		CommandInit:     commandInitSummary,
		CommandLogin:    commandLoginSummary,
		CommandMute:     commandMuteSummary,
		CommandReject:   commandRejectSummary,
		CommandRemove:   commandRemoveSummary,
		CommandShow:     commandShowSummary,
		CommandSwitch:   commandSwitchSummary,
		CommandUnblock:  commandUnblockSummary,
		CommandUnfollow: commandUnfollowSummary,
		CommandUnmute:   commandUnmuteSummary,
		CommandVersion:  commandVersionSummary,
		CommandWhoami:   commandWhoamiSummary,
	}
}

func CommandSummaryLookup(command string) string {
	commandMap := CommandSummaryMap()

	summary, ok := commandMap[command]
	if !ok {
		return "This command does not have a summary"
	}

	return summary
}
