package media

import (
	"fmt"
	"net/rpc"
	"path/filepath"
	"strings"

	"codeflow.dananglin.me.uk/apollo/enbas/internal/gtsclient"
	"codeflow.dananglin.me.uk/apollo/enbas/internal/model"
	"codeflow.dananglin.me.uk/apollo/enbas/internal/utilities"
)

const (
	mediaTypeImage string = "image"
	mediaTypeVideo string = "video"
)

type media struct {
	source      string
	destination string
	mediaType   string
}

func (m *media) download(client *rpc.Client) error {
	fileExists, err := utilities.FileExists(m.destination)
	if err != nil {
		return fmt.Errorf(
			"unable to determine if %s exists: %w",
			m.destination,
			err,
		)
	}

	if fileExists {
		return nil
	}

	if err := client.Call(
		"GTSClient.DownloadMedia",
		gtsclient.DownloadMediaArgs{
			URL:  m.source,
			Path: m.destination,
		},
		nil,
	); err != nil {
		return fmt.Errorf(
			"downloading %s -> %s failed: %w",
			m.source,
			m.destination,
			err,
		)
	}

	return nil
}

func newMediaHashmap(cacheDir string, attachments []model.Attachment) map[string]media {
	hashmap := make(map[string]media)

	for ind := range attachments {
		hashmap[attachments[ind].ID] = media{
			source:      attachments[ind].URL,
			destination: mediaFilepath(cacheDir, attachments[ind].URL),
			mediaType:   attachments[ind].Type,
		}
	}

	return hashmap
}

type Bundle struct {
	images []media
	videos []media
}

func NewBundle(
	cacheDir string,
	attachments []model.Attachment,
	getAllImages bool,
	getAllVideos bool,
	attachmentIDs []string,
) Bundle {
	mediaHashmap := newMediaHashmap(cacheDir, attachments)
	images := make([]media, 0)
	videos := make([]media, 0)

	if !getAllImages && !getAllVideos && len(attachmentIDs) == 0 {
		return Bundle{
			images: images,
			videos: videos,
		}
	}

	if getAllImages || getAllVideos {
		if getAllImages {
			for _, m := range mediaHashmap {
				if m.mediaType == mediaTypeImage {
					images = append(images, m)
				}
			}
		}

		if getAllVideos {
			for _, m := range mediaHashmap {
				if m.mediaType == mediaTypeVideo {
					videos = append(videos, m)
				}
			}
		}

		return Bundle{
			images: images,
			videos: videos,
		}
	}

	for _, attachmentID := range attachmentIDs {
		obj, ok := mediaHashmap[attachmentID]
		if !ok {
			continue
		}

		switch obj.mediaType {
		case mediaTypeImage:
			images = append(images, obj)
		case mediaTypeVideo:
			videos = append(videos, obj)
		}
	}

	return Bundle{
		images: images,
		videos: videos,
	}
}

func (m *Bundle) Download(client *rpc.Client) error {
	for ind := range m.images {
		if err := m.images[ind].download(client); err != nil {
			return fmt.Errorf("received an error trying to download the image files: %w", err)
		}
	}

	for ind := range m.videos {
		if err := m.videos[ind].download(client); err != nil {
			return fmt.Errorf("received an error trying to download the video files: %w", err)
		}
	}

	return nil
}

func (m *Bundle) ImageFiles() []string {
	filepaths := make([]string, len(m.images))

	for ind := range m.images {
		filepaths[ind] = m.images[ind].destination
	}

	return filepaths
}

func (m *Bundle) VideoFiles() []string {
	filepaths := make([]string, len(m.videos))

	for ind := range m.videos {
		filepaths[ind] = m.videos[ind].destination
	}

	return filepaths
}

func mediaFilepath(cacheDir, mediaURL string) string {
	split := strings.Split(mediaURL, "/")

	return filepath.Join(cacheDir, split[len(split)-1])
}
