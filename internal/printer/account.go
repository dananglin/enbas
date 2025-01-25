package printer

import (
	"strconv"
	"strings"

	"codeflow.dananglin.me.uk/apollo/enbas/internal/model"
)

func (p Printer) PrintAccount(
	account model.Account,
	relationship model.AccountRelationship,
	preferences model.Preferences,
	statusList model.StatusList,
	userAccountID string,
) {
	var builder strings.Builder

	builder.WriteString("\n" + p.fullDisplayNameFormat(account.DisplayName, account.Acct))
	builder.WriteString("\n\n" + p.headerFormat("ACCOUNT ID:"))
	builder.WriteString("\n" + account.ID)
	builder.WriteString("\n\n" + p.headerFormat("JOINED ON:"))
	builder.WriteString("\n" + p.formatDate(account.CreatedAt))
	builder.WriteString("\n\n" + p.headerFormat("STATS:"))
	builder.WriteString("\n" + p.fieldFormat("Followers:"))
	builder.WriteString(" " + strconv.Itoa(account.FollowersCount))
	builder.WriteString("\n" + p.fieldFormat("Following:"))
	builder.WriteString(" " + strconv.Itoa(account.FollowingCount))
	builder.WriteString("\n" + p.fieldFormat("Statuses:"))
	builder.WriteString(" " + strconv.Itoa(account.StatusCount))
	builder.WriteString("\n\n" + p.headerFormat("BIOGRAPHY:"))
	builder.WriteString(p.convertHTMLToText(account.Note, true))
	builder.WriteString("\n\n" + p.headerFormat("METADATA:"))

	for _, field := range account.Fields {
		builder.WriteString("\n" + p.fieldFormat(field.Name) + ": " + p.convertHTMLToText(field.Value, false))
	}

	builder.WriteString("\n\n" + p.headerFormat("ACCOUNT URL:"))
	builder.WriteString("\n" + account.URL)

	if relationship.Print {
		builder.WriteString(p.accountRelationship(relationship))
	}

	if preferences.Print {
		builder.WriteString(p.userPreferences(preferences))
	}

	if statusList.Statuses != nil {
		builder.WriteString("\n\n" + p.statusList(statusList, userAccountID))
	}

	builder.WriteString("\n\n")

	p.print(builder.String())
}

func (p Printer) accountRelationship(relationship model.AccountRelationship) string {
	var builder strings.Builder

	builder.WriteString("\n\n" + p.headerFormat("YOUR RELATIONSHIP WITH THIS ACCOUNT:"))
	builder.WriteString("\n" + p.fieldFormat("Following:"))
	builder.WriteString(" " + strconv.FormatBool(relationship.Following))
	builder.WriteString("\n" + p.fieldFormat("Is following you:"))
	builder.WriteString(" " + strconv.FormatBool(relationship.FollowedBy))
	builder.WriteString("\n" + p.fieldFormat("A follow request was sent and is pending:"))
	builder.WriteString(" " + strconv.FormatBool(relationship.FollowRequested))
	builder.WriteString("\n" + p.fieldFormat("Received a pending follow request:"))
	builder.WriteString(" " + strconv.FormatBool(relationship.FollowRequestedBy))
	builder.WriteString("\n" + p.fieldFormat("Endorsed:"))
	builder.WriteString(" " + strconv.FormatBool(relationship.Endorsed))
	builder.WriteString("\n" + p.fieldFormat("Showing Reposts (boosts):"))
	builder.WriteString(" " + strconv.FormatBool(relationship.ShowingReblogs))
	builder.WriteString("\n" + p.fieldFormat("Muted:"))
	builder.WriteString(" " + strconv.FormatBool(relationship.Muting))
	builder.WriteString("\n" + p.fieldFormat("Notifications muted:"))
	builder.WriteString(" " + strconv.FormatBool(relationship.MutingNotifications))
	builder.WriteString("\n" + p.fieldFormat("Blocking:"))
	builder.WriteString(" " + strconv.FormatBool(relationship.Blocking))
	builder.WriteString("\n" + p.fieldFormat("Is blocking you:"))
	builder.WriteString(" " + strconv.FormatBool(relationship.BlockedBy))
	builder.WriteString("\n" + p.fieldFormat("Blocking account's domain:"))
	builder.WriteString(" " + strconv.FormatBool(relationship.DomainBlocking))

	if relationship.PrivateNote != "" {
		builder.WriteString("\n\n" + p.headerFormat("YOUR PRIVATE NOTE ABOUT THIS ACCOUNT:"))
		builder.WriteString("\n" + p.wrapLines(relationship.PrivateNote, 0))
	}

	return builder.String()
}

func (p Printer) userPreferences(preferences model.Preferences) string {
	var builder strings.Builder

	builder.WriteString("\n\n" + p.headerFormat("YOUR PREFERENCES:"))

	builder.WriteString("\n" + p.fieldFormat("Default post language:"))
	builder.WriteString(" " + preferences.PostingDefaultLanguage)

	builder.WriteString("\n" + p.fieldFormat("Default post visibility:"))
	builder.WriteString(" " + preferences.PostingDefaultVisibility)

	builder.WriteString("\n" + p.fieldFormat("Mark posts as sensitive by default:"))
	builder.WriteString(" " + strconv.FormatBool(preferences.PostingDefaultSensitive))

	return builder.String()
}

func (p Printer) PrintAccountList(list model.AccountList) {
	var builder strings.Builder

	builder.WriteString("\n")

	switch list.Type {
	case model.AccountListFollowers:
		builder.WriteString(p.headerFormat("Followed by:"))
	case model.AccountListFollowing:
		builder.WriteString(p.headerFormat("Following:"))
	case model.AccountListBlockedAccount:
		builder.WriteString(p.headerFormat("Blocked accounts:"))
	case model.AccountListFollowRequests:
		builder.WriteString(p.headerFormat("Accounts that have requested to follow you:"))
	case model.AccountListMuted:
		builder.WriteString(p.headerFormat("Muted accounts:"))
	default:
		builder.WriteString(p.headerFormat("Accounts:"))
	}

	if list.Type == model.AccountListBlockedAccount {
		for ind := range list.Accounts {
			builder.WriteString("\n" + symbolBullet + " " + list.Accounts[ind].Acct + " (" + list.Accounts[ind].ID + ")")
		}
	} else {
		for ind := range list.Accounts {
			builder.WriteString("\n" + symbolBullet + " " + p.fullDisplayNameFormat(list.Accounts[ind].DisplayName, list.Accounts[ind].Acct))
		}
	}

	builder.WriteString("\n")

	p.print(builder.String())
}
