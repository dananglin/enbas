package executor

import (
	"fmt"
	"path/filepath"

	"codeflow.dananglin.me.uk/apollo/enbas/internal/client"
	"codeflow.dananglin.me.uk/apollo/enbas/internal/media"
	"codeflow.dananglin.me.uk/apollo/enbas/internal/model"
	"codeflow.dananglin.me.uk/apollo/enbas/internal/utilities"
)

func (s *ShowExecutor) Execute() error {
	if s.resourceType == "" {
		return FlagNotSetError{flagText: flagType}
	}

	funcMap := map[string]func(*client.Client) error{
		resourceInstance:        s.showInstance,
		resourceAccount:         s.showAccount,
		resourceStatus:          s.showStatus,
		resourceTimeline:        s.showTimeline,
		resourceList:            s.showList,
		resourceFollowers:       s.showFollowers,
		resourceFollowing:       s.showFollowing,
		resourceBlocked:         s.showBlocked,
		resourceBookmarks:       s.showBookmarks,
		resourceLiked:           s.showLiked,
		resourceStarred:         s.showLiked,
		resourceFollowRequest:   s.showFollowRequests,
		resourcePoll:            s.showPoll,
		resourceMutedAccounts:   s.showMutedAccounts,
		resourceMedia:           s.showMedia,
		resourceMediaAttachment: s.showMediaAttachment,
	}

	doFunc, ok := funcMap[s.resourceType]
	if !ok {
		return UnsupportedTypeError{resourceType: s.resourceType}
	}

	gtsClient, err := client.NewClientFromFile(s.config.CredentialsFile)
	if err != nil {
		return fmt.Errorf("unable to create the GoToSocial client: %w", err)
	}

	return doFunc(gtsClient)
}

func (s *ShowExecutor) showInstance(gtsClient *client.Client) error {
	instance, err := gtsClient.GetInstance()
	if err != nil {
		return fmt.Errorf("unable to retrieve the instance details: %w", err)
	}

	s.printer.PrintInstance(instance)

	return nil
}

func (s *ShowExecutor) showAccount(gtsClient *client.Client) error {
	account, err := getAccount(gtsClient, s.myAccount, s.accountName)
	if err != nil {
		return fmt.Errorf("unable to get the account information: %w", err)
	}

	if s.showInBrowser {
		if err := utilities.OpenLink(s.config.Integrations.Browser, account.URL); err != nil {
			return fmt.Errorf("unable to open link: %w", err)
		}

		return nil
	}

	var (
		relationship *model.AccountRelationship
		preferences  *model.Preferences
		statuses     *model.StatusList
	)

	if !s.myAccount && !s.skipAccountRelationship {
		relationship, err = gtsClient.GetAccountRelationship(account.ID)
		if err != nil {
			return fmt.Errorf("unable to retrieve the relationship to this account: %w", err)
		}
	}

	if s.myAccount && s.showUserPreferences {
		preferences, err = gtsClient.GetUserPreferences()
		if err != nil {
			return fmt.Errorf("unable to retrieve the user preferences: %w", err)
		}
	}

	if s.showStatuses {
		form := client.GetAccountStatusesForm{
			AccountID:      account.ID,
			Limit:          s.limit,
			ExcludeReplies: s.excludeReplies,
			ExcludeReblogs: s.excludeBoosts,
			Pinned:         s.onlyPinned,
			OnlyMedia:      s.onlyMedia,
			OnlyPublic:     s.onlyPublic,
		}

		statuses, err = gtsClient.GetAccountStatuses(form)
		if err != nil {
			return fmt.Errorf("unable to retrieve the account's statuses: %w", err)
		}
	}

	s.printer.PrintAccount(account, relationship, preferences, statuses)

	return nil
}

func (s *ShowExecutor) showStatus(gtsClient *client.Client) error {
	if s.statusID == "" {
		return FlagNotSetError{flagText: flagStatusID}
	}

	status, err := gtsClient.GetStatus(s.statusID)
	if err != nil {
		return fmt.Errorf("unable to retrieve the status: %w", err)
	}

	if s.showInBrowser {
		if err := utilities.OpenLink(s.config.Integrations.Browser, status.URL); err != nil {
			return fmt.Errorf("unable to open link: %w", err)
		}

		return nil
	}

	s.printer.PrintStatus(status)

	return nil
}

func (s *ShowExecutor) showTimeline(gtsClient *client.Client) error {
	var (
		timeline model.StatusList
		err      error
	)

	switch s.timelineCategory {
	case model.TimelineCategoryHome:
		timeline, err = gtsClient.GetHomeTimeline(s.limit)
	case model.TimelineCategoryPublic:
		timeline, err = gtsClient.GetPublicTimeline(s.limit)
	case model.TimelineCategoryList:
		if s.listID == "" {
			return FlagNotSetError{flagText: flagListID}
		}

		var list model.List

		list, err = gtsClient.GetList(s.listID)
		if err != nil {
			return fmt.Errorf("unable to retrieve the list: %w", err)
		}

		timeline, err = gtsClient.GetListTimeline(list.ID, list.Title, s.limit)
	case model.TimelineCategoryTag:
		if s.tag == "" {
			return FlagNotSetError{flagText: flagTag}
		}

		timeline, err = gtsClient.GetTagTimeline(s.tag, s.limit)
	default:
		return model.InvalidTimelineCategoryError{Value: s.timelineCategory}
	}

	if err != nil {
		return fmt.Errorf("unable to retrieve the %s timeline: %w", s.timelineCategory, err)
	}

	if len(timeline.Statuses) == 0 {
		s.printer.PrintInfo("There are no statuses in this timeline.\n")

		return nil
	}

	s.printer.PrintStatusList(timeline)

	return nil
}

func (s *ShowExecutor) showList(gtsClient *client.Client) error {
	if s.listID == "" {
		return s.showLists(gtsClient)
	}

	list, err := gtsClient.GetList(s.listID)
	if err != nil {
		return fmt.Errorf("unable to retrieve the list: %w", err)
	}

	accounts, err := gtsClient.GetAccountsFromList(s.listID, 0)
	if err != nil {
		return fmt.Errorf("unable to retrieve the accounts from the list: %w", err)
	}

	if len(accounts) > 0 {
		accountMap := make(map[string]string)
		for i := range accounts {
			accountMap[accounts[i].Acct] = accounts[i].Username
		}

		list.Accounts = accountMap
	}

	s.printer.PrintList(list)

	return nil
}

func (s *ShowExecutor) showLists(gtsClient *client.Client) error {
	lists, err := gtsClient.GetAllLists()
	if err != nil {
		return fmt.Errorf("unable to retrieve the lists: %w", err)
	}

	if len(lists) == 0 {
		s.printer.PrintInfo("You have no lists.\n")

		return nil
	}

	s.printer.PrintLists(lists)

	return nil
}

func (s *ShowExecutor) showFollowers(gtsClient *client.Client) error {
	if s.fromResourceType == "" {
		return FlagNotSetError{flagText: flagFrom}
	}

	funcMap := map[string]func(*client.Client) error{
		resourceAccount: s.showFollowersFromAccount,
	}

	doFunc, ok := funcMap[s.fromResourceType]
	if !ok {
		return UnsupportedShowOperationError{
			ResourceType:         s.resourceType,
			ShowFromResourceType: s.fromResourceType,
		}
	}

	return doFunc(gtsClient)
}

func (s *ShowExecutor) showFollowersFromAccount(gtsClient *client.Client) error {
	accountID, err := getAccountID(gtsClient, s.myAccount, s.accountName, s.config.CredentialsFile)
	if err != nil {
		return fmt.Errorf("received an error while getting the account ID: %w", err)
	}

	followers, err := gtsClient.GetFollowers(accountID, s.limit)
	if err != nil {
		return fmt.Errorf("unable to retrieve the list of followers: %w", err)
	}

	if len(followers.Accounts) > 0 {
		s.printer.PrintAccountList(followers)
	} else {
		s.printer.PrintInfo("There are no followers for this account (or the list is hidden).\n")
	}

	return nil
}

func (s *ShowExecutor) showFollowing(gtsClient *client.Client) error {
	if s.fromResourceType == "" {
		return FlagNotSetError{flagText: flagFrom}
	}

	funcMap := map[string]func(*client.Client) error{
		resourceAccount: s.showFollowingFromAccount,
	}

	doFunc, ok := funcMap[s.fromResourceType]
	if !ok {
		return UnsupportedShowOperationError{
			ResourceType:         s.resourceType,
			ShowFromResourceType: s.fromResourceType,
		}
	}

	return doFunc(gtsClient)
}

func (s *ShowExecutor) showFollowingFromAccount(gtsClient *client.Client) error {
	accountID, err := getAccountID(gtsClient, s.myAccount, s.accountName, s.config.CredentialsFile)
	if err != nil {
		return fmt.Errorf("received an error while getting the account ID: %w", err)
	}

	following, err := gtsClient.GetFollowing(accountID, s.limit)
	if err != nil {
		return fmt.Errorf("unable to retrieve the list of followed accounts: %w", err)
	}

	if len(following.Accounts) > 0 {
		s.printer.PrintAccountList(following)
	} else {
		s.printer.PrintInfo("This account is not following anyone or the list is hidden.\n")
	}

	return nil
}

func (s *ShowExecutor) showBlocked(gtsClient *client.Client) error {
	blocked, err := gtsClient.GetBlockedAccounts(s.limit)
	if err != nil {
		return fmt.Errorf("unable to retrieve the list of blocked accounts: %w", err)
	}

	if len(blocked.Accounts) > 0 {
		s.printer.PrintAccountList(blocked)
	} else {
		s.printer.PrintInfo("You have no blocked accounts.\n")
	}

	return nil
}

func (s *ShowExecutor) showBookmarks(gtsClient *client.Client) error {
	bookmarks, err := gtsClient.GetBookmarks(s.limit)
	if err != nil {
		return fmt.Errorf("unable to retrieve the list of bookmarks: %w", err)
	}

	if len(bookmarks.Statuses) > 0 {
		s.printer.PrintStatusList(bookmarks)
	} else {
		s.printer.PrintInfo("You have no bookmarks.\n")
	}

	return nil
}

func (s *ShowExecutor) showLiked(gtsClient *client.Client) error {
	liked, err := gtsClient.GetLikedStatuses(s.limit, s.resourceType)
	if err != nil {
		return fmt.Errorf("unable to retrieve the list of your %s statuses: %w", s.resourceType, err)
	}

	if len(liked.Statuses) > 0 {
		s.printer.PrintStatusList(liked)
	} else {
		s.printer.PrintInfo("You have no " + s.resourceType + " statuses.\n")
	}

	return nil
}

func (s *ShowExecutor) showFollowRequests(gtsClient *client.Client) error {
	accounts, err := gtsClient.GetFollowRequests(s.limit)
	if err != nil {
		return fmt.Errorf("unable to retrieve the list of follow requests: %w", err)
	}

	if len(accounts.Accounts) > 0 {
		s.printer.PrintAccountList(accounts)
	} else {
		s.printer.PrintInfo("You have no follow requests.\n")
	}

	return nil
}

func (s *ShowExecutor) showPoll(gtsClient *client.Client) error {
	if s.pollID == "" {
		return FlagNotSetError{flagText: flagPollID}
	}

	poll, err := gtsClient.GetPoll(s.pollID)
	if err != nil {
		return fmt.Errorf("unable to retrieve the poll: %w", err)
	}

	s.printer.PrintPoll(poll)

	return nil
}

func (s *ShowExecutor) showMutedAccounts(gtsClient *client.Client) error {
	muted, err := gtsClient.GetMutedAccounts(s.limit)
	if err != nil {
		return fmt.Errorf("unable to retrieve the list of muted accounts: %w", err)
	}

	if len(muted.Accounts) > 0 {
		s.printer.PrintAccountList(muted)
	} else {
		s.printer.PrintInfo("You have not muted any accounts.\n")
	}

	return nil
}

func (s *ShowExecutor) showMediaAttachment(gtsClient *client.Client) error {
	if len(s.attachmentIDs) == 0 {
		return FlagNotSetError{flagText: flagAttachmentID}
	}

	if len(s.attachmentIDs) != 1 {
		return fmt.Errorf(
			"unexpected number of attachment IDs received: want 1, got %d",
			len(s.attachmentIDs),
		)
	}

	attachment, err := gtsClient.GetMediaAttachment(s.attachmentIDs[0])
	if err != nil {
		return fmt.Errorf("unable to retrieve the media attachment: %w", err)
	}

	s.printer.PrintMediaAttachment(attachment)

	return nil
}

func (s *ShowExecutor) showMedia(gtsClient *client.Client) error {
	if s.fromResourceType == "" {
		return FlagNotSetError{flagText: flagFrom}
	}

	funcMap := map[string]func(*client.Client) error{
		resourceStatus: s.showMediaFromStatus,
	}

	doFunc, ok := funcMap[s.fromResourceType]
	if !ok {
		return UnsupportedShowOperationError{
			ResourceType:         s.resourceType,
			ShowFromResourceType: s.fromResourceType,
		}
	}

	return doFunc(gtsClient)
}

func (s *ShowExecutor) showMediaFromStatus(gtsClient *client.Client) error {
	if s.statusID == "" {
		return FlagNotSetError{flagText: flagStatusID}
	}

	status, err := gtsClient.GetStatus(s.statusID)
	if err != nil {
		return fmt.Errorf("unable to retrieve the status: %w", err)
	}

	cacheDir := filepath.Join(
		utilities.CalculateCacheDir(s.config.CacheDirectory, utilities.GetFQDN(gtsClient.Authentication.Instance)),
		"media",
	)

	if err := utilities.EnsureDirectory(cacheDir); err != nil {
		return fmt.Errorf("unable to ensure the existence of the directory %q: %w", cacheDir, err)
	}

	mediaBundle := media.NewBundle(
		cacheDir,
		status.MediaAttachments,
		s.getAllImages,
		s.getAllVideos,
		s.attachmentIDs,
	)

	if err := mediaBundle.Download(gtsClient); err != nil {
		return fmt.Errorf("unable to download the media bundle: %w", err)
	}

	imageFiles := mediaBundle.ImageFiles()
	if len(imageFiles) > 0 {
		if err := utilities.OpenMedia(s.config.Integrations.ImageViewer, imageFiles); err != nil {
			return fmt.Errorf("unable to open the image viewer: %w", err)
		}
	}

	videoFiles := mediaBundle.VideoFiles()
	if len(videoFiles) > 0 {
		if err := utilities.OpenMedia(s.config.Integrations.VideoPlayer, videoFiles); err != nil {
			return fmt.Errorf("unable to open the video player: %w", err)
		}
	}

	return nil
}
