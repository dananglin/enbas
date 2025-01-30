package executor

import (
	"fmt"
	"net/rpc"

	"codeflow.dananglin.me.uk/apollo/enbas/internal/gtsclient"
	"codeflow.dananglin.me.uk/apollo/enbas/internal/media"
	"codeflow.dananglin.me.uk/apollo/enbas/internal/model"
	"codeflow.dananglin.me.uk/apollo/enbas/internal/server"
	"codeflow.dananglin.me.uk/apollo/enbas/internal/utilities"
)

func (s *ShowExecutor) Execute() error {
	if s.resourceType == "" {
		return FlagNotSetError{flagText: flagType}
	}

	funcMap := map[string]func(*rpc.Client) error{
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
		resourceMutedAccounts:   s.showMutedAccounts,
		resourceMedia:           s.showMedia,
		resourceMediaAttachment: s.showMediaAttachment,
		resourceFollowedTags:    s.showFollowedTags,
		resourceTag:             s.showTag,
		resourceThread:          s.showThread,
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

func (s *ShowExecutor) showInstance(client *rpc.Client) error {
	var instance model.InstanceV2
	if err := client.Call("GTSClient.GetInstance", gtsclient.NoRPCArgs{}, &instance); err != nil {
		return fmt.Errorf("unable to retrieve the instance details: %w", err)
	}

	s.printer.PrintInstance(instance)

	return nil
}

func (s *ShowExecutor) showAccount(client *rpc.Client) error {
	account, err := getAccount(client, s.myAccount, s.accountName)
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
		relationship model.AccountRelationship
		preferences  model.Preferences
		statusList   model.StatusList
		myAccountID  string
	)

	relationship.Print = false
	preferences.Print = false
	statusList.Statuses = nil

	if !s.myAccount && !s.skipAccountRelationship {
		if err := client.Call("GTSClient.GetAccountRelationship", account.ID, &relationship); err != nil {
			return fmt.Errorf("unable to retrieve the relationship to this account: %w", err)
		}

		relationship.Print = true
	}

	if s.myAccount {
		myAccountID = account.ID
		if s.showUserPreferences {
			if err := client.Call("GTSClient.GetUserPreferences", gtsclient.NoRPCArgs{}, &preferences); err != nil {
				return fmt.Errorf("unable to retrieve the user preferences: %w", err)
			}

			preferences.Print = true
		}
	}

	if s.showStatuses {
		args := gtsclient.GetAccountStatusesArgs{
			AccountID:      account.ID,
			Limit:          s.limit,
			ExcludeReplies: s.excludeReplies,
			ExcludeReblogs: s.excludeBoosts,
			Pinned:         s.onlyPinned,
			OnlyMedia:      s.onlyMedia,
			OnlyPublic:     s.onlyPublic,
		}

		if err := client.Call("GTSClient.GetAccountStatuses", args, &statusList); err != nil {
			return fmt.Errorf("unable to retrieve the account's statuses: %w", err)
		}
	}

	s.printer.PrintAccount(account, relationship, preferences, statusList, myAccountID)

	return nil
}

func (s *ShowExecutor) showStatus(client *rpc.Client) error {
	if s.statusID == "" {
		return MissingIDError{
			resource: resourceStatus,
			action:   "view",
		}
	}

	var (
		status    model.Status
		boostedBy model.AccountList
		likedBy   model.AccountList
	)

	if err := client.Call("GTSClient.GetStatus", s.statusID, &status); err != nil {
		return fmt.Errorf("error retrieving the status: %w", err)
	}

	if s.showInBrowser {
		if err := utilities.OpenLink(s.config.Integrations.Browser, status.URL); err != nil {
			return fmt.Errorf("unable to open link: %w", err)
		}

		return nil
	}

	boostedBy.Accounts = nil
	likedBy.Accounts = nil

	if s.boostedBy {
		if err := client.Call("GTSClient.GetAccountsWhoRebloggedStatus", s.statusID, &boostedBy); err != nil {
			return fmt.Errorf(
				"error retrieving the list of accounts that boosted the status: %w",
				err,
			)
		}
	}

	if s.likedBy {
		if err := client.Call("GTSClient.GetAccountsWhoLikedStatus", s.statusID, &likedBy); err != nil {
			return fmt.Errorf(
				"error retrieving the list of accounts that liked the status: %w",
				err,
			)
		}
	}

	myAccountID, err := getAccountID(client, true, nil)
	if err != nil {
		return fmt.Errorf("unable to get your account ID: %w", err)
	}

	s.printer.PrintStatus(
		status,
		myAccountID,
		boostedBy,
		likedBy,
	)

	return nil
}

func (s *ShowExecutor) showTimeline(client *rpc.Client) error {
	var (
		timeline model.StatusList
		err      error
	)

	switch s.timelineCategory {
	case model.TimelineCategoryHome:
		err = client.Call("GTSClient.GetHomeTimeline", s.limit, &timeline)
	case model.TimelineCategoryPublic:
		err = client.Call("GTSClient.GetPublicTimeline", s.limit, &timeline)
	case model.TimelineCategoryList:
		if s.listID == "" {
			return MissingIDError{
				resource: resourceList,
				action:   "view the timeline in",
			}
		}

		var list model.List

		if err := client.Call("GTSClient.GetList", s.listID, &list); err != nil {
			return fmt.Errorf("unable to retrieve the list: %w", err)
		}

		err = client.Call(
			"GTSClient.GetListTimeline",
			gtsclient.GetListTimelineArgs{
				ListID: list.ID,
				Title:  list.Title,
				Limit:  s.limit,
			},
			&timeline,
		)
	case model.TimelineCategoryTag:
		if s.tag == "" {
			return Error{"please provide the name of the tag"}
		}

		err = client.Call(
			"GTSClient.GetTagTimeline",
			gtsclient.GetTagTimelineArgs{
				TagName: s.tag,
				Limit:   s.limit,
			},
			&timeline,
		)
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

	myAccountID, err := getAccountID(client, true, nil)
	if err != nil {
		return fmt.Errorf("unable to get your account ID: %w", err)
	}

	s.printer.PrintStatusList(timeline, myAccountID)

	return nil
}

func (s *ShowExecutor) showList(client *rpc.Client) error {
	if s.listID == "" {
		return s.showLists(client)
	}

	var list model.List

	if err := client.Call("GTSClient.GetList", s.listID, &list); err != nil {
		return fmt.Errorf("unable to retrieve the list: %w", err)
	}

	acctMap, err := getAccountsFromList(client, s.listID)
	if err != nil {
		return err
	}

	if len(acctMap) > 0 {
		list.Accounts = acctMap
	}

	s.printer.PrintList(list)

	return nil
}

func (s *ShowExecutor) showLists(client *rpc.Client) error {
	var lists []model.List
	if err := client.Call("GTSClient.GetAllLists", "", &lists); err != nil {
		return fmt.Errorf("unable to retrieve the lists: %w", err)
	}

	if len(lists) == 0 {
		s.printer.PrintInfo("You have no lists.\n")

		return nil
	}

	s.printer.PrintLists(lists)

	return nil
}

func (s *ShowExecutor) showFollowers(client *rpc.Client) error {
	if s.fromResourceType == "" {
		return FlagNotSetError{flagText: flagFrom}
	}

	funcMap := map[string]func(*rpc.Client) error{
		resourceAccount: s.showFollowersFromAccount,
	}

	doFunc, ok := funcMap[s.fromResourceType]
	if !ok {
		return UnsupportedShowOperationError{
			resourceType:         s.resourceType,
			showFromResourceType: s.fromResourceType,
		}
	}

	return doFunc(client)
}

func (s *ShowExecutor) showFollowersFromAccount(client *rpc.Client) error {
	accountID, err := getAccountID(client, s.myAccount, s.accountName)
	if err != nil {
		return fmt.Errorf("received an error while getting the account ID: %w", err)
	}

	var followers model.AccountList
	if err := client.Call(
		"GTSClient.GetFollowers",
		gtsclient.GetFollowersArgs{
			AccountID: accountID,
			Limit:     s.limit,
		},
		&followers,
	); err != nil {
		return fmt.Errorf("unable to retrieve the list of followers: %w", err)
	}

	if len(followers.Accounts) > 0 {
		s.printer.PrintAccountList(followers)
	} else {
		s.printer.PrintInfo("There are no followers for this account (or the list is hidden).\n")
	}

	return nil
}

func (s *ShowExecutor) showFollowing(client *rpc.Client) error {
	if s.fromResourceType == "" {
		return FlagNotSetError{flagText: flagFrom}
	}

	funcMap := map[string]func(*rpc.Client) error{
		resourceAccount: s.showFollowingFromAccount,
	}

	doFunc, ok := funcMap[s.fromResourceType]
	if !ok {
		return UnsupportedShowOperationError{
			resourceType:         s.resourceType,
			showFromResourceType: s.fromResourceType,
		}
	}

	return doFunc(client)
}

func (s *ShowExecutor) showFollowingFromAccount(client *rpc.Client) error {
	accountID, err := getAccountID(client, s.myAccount, s.accountName)
	if err != nil {
		return fmt.Errorf("received an error while getting the account ID: %w", err)
	}

	var followings model.AccountList
	if err := client.Call(
		"GTSClient.GetFollowing",
		gtsclient.GetFollowingsArgs{
			AccountID: accountID,
			Limit:     s.limit,
		},
		&followings,
	); err != nil {
		return fmt.Errorf("unable to retrieve the list of followed accounts: %w", err)
	}

	if len(followings.Accounts) > 0 {
		s.printer.PrintAccountList(followings)
	} else {
		s.printer.PrintInfo("This account is not following anyone or the list is hidden.\n")
	}

	return nil
}

func (s *ShowExecutor) showBlocked(client *rpc.Client) error {
	var blocked model.AccountList
	if err := client.Call("GTSClient.GetBlockedAccounts", s.limit, &blocked); err != nil {
		return fmt.Errorf("unable to retrieve the list of blocked accounts: %w", err)
	}

	if len(blocked.Accounts) > 0 {
		s.printer.PrintAccountList(blocked)
	} else {
		s.printer.PrintInfo("You have no blocked accounts.\n")
	}

	return nil
}

func (s *ShowExecutor) showBookmarks(client *rpc.Client) error {
	var bookmarks model.StatusList
	if err := client.Call("GTSClient.GetBookmarks", s.limit, &bookmarks); err != nil {
		return fmt.Errorf("unable to retrieve the list of bookmarks: %w", err)
	}

	if len(bookmarks.Statuses) > 0 {
		myAccountID, err := getAccountID(client, true, nil)
		if err != nil {
			return fmt.Errorf("unable to get your account ID: %w", err)
		}

		s.printer.PrintStatusList(bookmarks, myAccountID)
	} else {
		s.printer.PrintInfo("You have no bookmarks.\n")
	}

	return nil
}

func (s *ShowExecutor) showLiked(client *rpc.Client) error {
	var liked model.StatusList
	if err := client.Call(
		"GTSClient.GetLikedStatuses",
		gtsclient.GetLikedStatusesArgs{
			Limit:        s.limit,
			ResourceType: s.resourceType,
		},
		&liked,
	); err != nil {
		return fmt.Errorf("unable to retrieve the list of your %s statuses: %w", s.resourceType, err)
	}

	if len(liked.Statuses) > 0 {
		myAccountID, err := getAccountID(client, true, nil)
		if err != nil {
			return fmt.Errorf("unable to get your account ID: %w", err)
		}

		s.printer.PrintStatusList(liked, myAccountID)
	} else {
		s.printer.PrintInfo("You have no " + s.resourceType + " statuses.\n")
	}

	return nil
}

func (s *ShowExecutor) showFollowRequests(client *rpc.Client) error {
	var requests model.AccountList
	if err := client.Call("GTSClient.GetFollowRequests", s.limit, &requests); err != nil {
		return fmt.Errorf("unable to retrieve the list of follow requests: %w", err)
	}

	if len(requests.Accounts) > 0 {
		s.printer.PrintAccountList(requests)
	} else {
		s.printer.PrintInfo("You have no follow requests.\n")
	}

	return nil
}

func (s *ShowExecutor) showMutedAccounts(client *rpc.Client) error {
	var muted model.AccountList
	if err := client.Call("GTSClient.GetMutedAccounts", s.limit, &muted); err != nil {
		return fmt.Errorf("unable to retrieve the list of muted accounts: %w", err)
	}

	if len(muted.Accounts) > 0 {
		s.printer.PrintAccountList(muted)
	} else {
		s.printer.PrintInfo("You have not muted any accounts.\n")
	}

	return nil
}

func (s *ShowExecutor) showMediaAttachment(client *rpc.Client) error {
	if len(s.attachmentIDs) != 1 {
		return fmt.Errorf(
			"unexpected number of attachment IDs received: want 1, got %d",
			len(s.attachmentIDs),
		)
	}

	var attachment model.MediaAttachment
	if err := client.Call("GTSClient.GetMediaAttachment", s.attachmentIDs[0], &attachment); err != nil {
		return fmt.Errorf("unable to retrieve the media attachment: %w", err)
	}

	s.printer.PrintMediaAttachment(attachment)

	return nil
}

func (s *ShowExecutor) showMedia(client *rpc.Client) error {
	if s.fromResourceType == "" {
		return FlagNotSetError{flagText: flagFrom}
	}

	funcMap := map[string]func(*rpc.Client) error{
		resourceStatus: s.showMediaFromStatus,
	}

	doFunc, ok := funcMap[s.fromResourceType]
	if !ok {
		return UnsupportedShowOperationError{
			resourceType:         s.resourceType,
			showFromResourceType: s.fromResourceType,
		}
	}

	return doFunc(client)
}

func (s *ShowExecutor) showMediaFromStatus(client *rpc.Client) error {
	if s.statusID == "" {
		return MissingIDError{
			resource: resourceStatus,
			action:   "view the media from",
		}
	}

	var (
		status      model.Status
		instanceURL string
	)

	if err := client.Call("GTSClient.GetStatus", s.statusID, &status); err != nil {
		return fmt.Errorf("unable to retrieve the status: %w", err)
	}

	if err := client.Call("GTSClient.GetInstanceURL", gtsclient.NoRPCArgs{}, &instanceURL); err != nil {
		return fmt.Errorf("unable to retrieve the instance URL: %w", err)
	}

	cacheDir, err := utilities.CalculateMediaCacheDir(s.config.CacheDirectory, instanceURL)
	if err != nil {
		return fmt.Errorf("unable to calculate the media cache directory: %w", err)
	}

	if err := utilities.EnsureDirectory(cacheDir); err != nil {
		return fmt.Errorf("unable to ensure the existence of the directory %q: %w", cacheDir, err)
	}

	mediaBundle := media.NewBundle(
		cacheDir,
		status.MediaAttachments,
		s.getAllAudio,
		s.getAllImages,
		s.getAllVideos,
		s.attachmentIDs,
	)

	if err := mediaBundle.Download(client); err != nil {
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

	audioFiles := mediaBundle.AudioFiles()
	if len(audioFiles) > 0 {
		if err := utilities.OpenMedia(s.config.Integrations.AudioPlayer, audioFiles); err != nil {
			return fmt.Errorf("unable to open the audio player: %w", err)
		}
	}

	return nil
}

func (s *ShowExecutor) showTag(client *rpc.Client) error {
	if s.tag == "" {
		return Error{"please provide the name of the tag"}
	}

	var tag model.Tag
	if err := client.Call("GTSClient.GetTag", s.tag, &tag); err != nil {
		return fmt.Errorf("unable to get the details of the tag: %w", err)
	}

	s.printer.PrintTag(tag)

	return nil
}

func (s *ShowExecutor) showFollowedTags(client *rpc.Client) error {
	var list model.TagList
	if err := client.Call("GTSClient.GetFollowedTags", s.limit, &list); err != nil {
		return fmt.Errorf("unable to get the list of followed tags: %w", err)
	}

	if len(list.Tags) > 0 {
		s.printer.PrintTagList(list)
	} else {
		s.printer.PrintInfo("This account is not following any tags.\n")
	}

	return nil
}

func (s *ShowExecutor) showThread(client *rpc.Client) error {
	if s.fromResourceType == "" {
		return FlagNotSetError{flagText: flagFrom}
	}

	funcMap := map[string]func(*rpc.Client) error{
		resourceStatus: s.showThreadFromStatus,
	}

	doFunc, ok := funcMap[s.fromResourceType]
	if !ok {
		return UnsupportedShowOperationError{
			resourceType:         s.resourceType,
			showFromResourceType: s.fromResourceType,
		}
	}

	return doFunc(client)
}

func (s *ShowExecutor) showThreadFromStatus(client *rpc.Client) error {
	if s.statusID == "" {
		return MissingIDError{
			resource: resourceStatus,
			action:   "view the media from",
		}
	}

	myAccountID, err := getAccountID(client, true, nil)
	if err != nil {
		return fmt.Errorf("unable to get your account ID: %w", err)
	}

	var thread model.Thread
	if err := client.Call("GTSClient.GetThread", s.statusID, &thread); err != nil {
		return fmt.Errorf("error retrieving the thread: %w", err)
	}

	if err := client.Call("GTSClient.GetStatus", s.statusID, &thread.Context); err != nil {
		return fmt.Errorf("error retrieving the status in context: %w", err)
	}

	// Print the thread
	s.printer.PrintThread(thread, myAccountID)

	return nil
}
