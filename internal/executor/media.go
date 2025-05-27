package executor

import (
	"fmt"
	"net/rpc"

	"codeflow.dananglin.me.uk/apollo/enbas/internal/cli"
	"codeflow.dananglin.me.uk/apollo/enbas/internal/command"
	"codeflow.dananglin.me.uk/apollo/enbas/internal/config"
	internalFlag "codeflow.dananglin.me.uk/apollo/enbas/internal/flag"
	"codeflow.dananglin.me.uk/apollo/enbas/internal/gtsclient"
	"codeflow.dananglin.me.uk/apollo/enbas/internal/media"
	"codeflow.dananglin.me.uk/apollo/enbas/internal/model"
	"codeflow.dananglin.me.uk/apollo/enbas/internal/printer"
	"codeflow.dananglin.me.uk/apollo/enbas/internal/server"
	"codeflow.dananglin.me.uk/apollo/enbas/internal/utilities"
)

func mediaFunc(
	cfg config.Config,
	_ printer.Settings,
	cmd command.Command,
) error {
	if cfg.IsZero() {
		return zeroConfigurationError{path: cfg.Path}
	}

	// Create the session to interact with the GoToSocial instance.
	session, err := server.StartSession(cfg.Server, cfg.Path)
	if err != nil {
		return fmt.Errorf("error creating the client for the daemon process: %w", err)
	}
	defer server.EndSession(session)

	switch cmd.Action {
	case cli.ActionShow:
		return mediaShow(
			session.Client(),
			cfg,
			cmd.RelatedTarget,
			cmd.RelatedTargetFlags,
		)
	default:
		return unsupportedActionError{action: cmd.Action, target: cli.TargetMedia}
	}
}

func mediaShow(
	client *rpc.Client,
	cfg config.Config,
	relatedTarget string,
	relatedTargetFlags []string,
) error {
	switch relatedTarget {
	case cli.TargetStatus:
		return mediaShowFromStatus(
			client,
			cfg.CacheDirectory,
			cfg.Integrations.AudioPlayer,
			cfg.Integrations.ImageViewer,
			cfg.Integrations.VideoPlayer,
			relatedTargetFlags,
		)
	default:
		return unsupportedTargetToTargetError{
			action:        cli.ActionShow,
			focusedTarget: cli.TargetMedia,
			preposition:   cli.TargetActionPreposition(cli.TargetMedia, cli.ActionShow),
			relatedTarget: relatedTarget,
		}
	}
}

func mediaShowFromStatus(
	client *rpc.Client,
	rootCacheDir string,
	audioPlayer string,
	imageViewer string,
	videoPlayer string,
	flags []string,
) error {
	var (
		statusID      string
		attachmentIDs internalFlag.MultiStringValue
		getAllAudio   bool
		getAllImages  bool
		getAllVideos  bool
	)

	// Parse the remaining flags
	if err := cli.ParseMediaShowFromStatusFlags(
		&statusID,
		&attachmentIDs,
		&getAllAudio,
		&getAllImages,
		&getAllVideos,
		flags,
	); err != nil {
		return err
	}

	if statusID == "" {
		return missingIDError{
			target: cli.TargetStatus,
			action: "view the media from",
		}
	}

	var (
		status      model.Status
		instanceURL string
	)

	if err := client.Call(
		"GTSClient.GetStatus",
		statusID,
		&status,
	); err != nil {
		return fmt.Errorf("error retrieving the status: %w", err)
	}

	if err := client.Call(
		"GTSClient.GetInstanceURL",
		gtsclient.NoRPCArgs{},
		&instanceURL,
	); err != nil {
		return fmt.Errorf("error retrieving the instance URL: %w", err)
	}

	cacheDir, err := utilities.CalculateMediaCacheDir(rootCacheDir, instanceURL)
	if err != nil {
		return fmt.Errorf("unable to calculate the media cache directory: %w", err)
	}

	if err := utilities.EnsureDirectory(cacheDir); err != nil {
		return fmt.Errorf("unable to ensure the existence of the directory %q: %w", cacheDir, err)
	}

	mediaBundle := media.NewBundle(
		cacheDir,
		status.MediaAttachments,
		getAllAudio,
		getAllImages,
		getAllVideos,
		attachmentIDs.Values(),
	)

	if err := mediaBundle.Download(client); err != nil {
		return fmt.Errorf("unable to download the media bundle: %w", err)
	}

	imageFiles := mediaBundle.ImageFiles()
	if len(imageFiles) > 0 {
		if err := utilities.OpenMedia(imageViewer, imageFiles); err != nil {
			return fmt.Errorf("unable to open the image viewer: %w", err)
		}
	}

	videoFiles := mediaBundle.VideoFiles()
	if len(videoFiles) > 0 {
		if err := utilities.OpenMedia(videoPlayer, videoFiles); err != nil {
			return fmt.Errorf("unable to open the video player: %w", err)
		}
	}

	audioFiles := mediaBundle.AudioFiles()
	if len(audioFiles) > 0 {
		if err := utilities.OpenMedia(audioPlayer, audioFiles); err != nil {
			return fmt.Errorf("unable to open the audio player: %w", err)
		}
	}

	return nil
}
