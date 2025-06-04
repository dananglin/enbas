package gtsclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"

	"codeflow.dananglin.me.uk/apollo/enbas/internal/model"
	"codeflow.dananglin.me.uk/apollo/enbas/internal/utilities"
)

const (
	baseMediaPath string = "/api/v1/media"
)

func (g *GTSClient) GetMediaAttachment(mediaAttachmentID string, attachment *model.MediaAttachment) error {
	params := requestParameters{
		httpMethod:  http.MethodGet,
		url:         g.auth.GetInstanceURL() + baseMediaPath + "/" + mediaAttachmentID,
		requestBody: nil,
		contentType: "",
		output:      attachment,
	}

	if err := g.sendRequest(params); err != nil {
		return fmt.Errorf("received an error after sending the request to get the media attachment: %w", err)
	}

	return nil
}

type CreateMediaAttachmentArgs struct {
	Path        string
	Description string
	Focus       string
}

func (g *GTSClient) CreateMediaAttachment(args CreateMediaAttachmentArgs, attachment *model.MediaAttachment) error {
	file, err := utilities.OpenFile(args.Path)
	if err != nil {
		return fmt.Errorf("unable to open the file: %w", err)
	}
	defer file.Close()

	// create the request body using a writer from the multipart package
	requestBody := bytes.Buffer{}
	requestBodyWriter := multipart.NewWriter(&requestBody)

	filename := filepath.Base(args.Path)

	part, err := requestBodyWriter.CreateFormFile("file", filename)
	if err != nil {
		return fmt.Errorf("unable to create the new part: %w", err)
	}

	if _, err := io.Copy(part, file); err != nil {
		return fmt.Errorf("unable to copy the file contents to the form: %w", err)
	}

	// add the description
	if args.Description != "" {
		descriptionFormFieldWriter, err := requestBodyWriter.CreateFormField("description")
		if err != nil {
			return fmt.Errorf(
				"unable to create the writer for the 'description' form field: %w",
				err,
			)
		}

		if _, err := io.WriteString(descriptionFormFieldWriter, args.Description); err != nil {
			return fmt.Errorf(
				"unable to write the description to the form: %w",
				err,
			)
		}
	}

	// add the focus values
	if args.Focus != "" {
		focusFormFieldWriter, err := requestBodyWriter.CreateFormField("focus")
		if err != nil {
			return fmt.Errorf(
				"unable to create the writer for the 'focus' form field: %w",
				err,
			)
		}

		if _, err := io.WriteString(focusFormFieldWriter, args.Focus); err != nil {
			return fmt.Errorf(
				"unable to write the focus values to the form: %w",
				err,
			)
		}
	}

	if err := requestBodyWriter.Close(); err != nil {
		return fmt.Errorf("unable to close the writer: %w", err)
	}

	url := g.auth.GetInstanceURL() + baseMediaPath

	params := requestParameters{
		httpMethod:  http.MethodPost,
		url:         url,
		requestBody: &requestBody,
		contentType: requestBodyWriter.FormDataContentType(),
		output:      attachment,
	}

	if err := g.sendRequest(params); err != nil {
		return fmt.Errorf(
			"received an error after sending the request to create the media attachment: %w",
			err,
		)
	}

	return nil
}

type UpdateMediaAttachmentArgs struct {
	MediaAttachmentID string
	Description       string
	Focus             string
}

func (g *GTSClient) UpdateMediaAttachment(args UpdateMediaAttachmentArgs, updated *model.MediaAttachment) error {
	form := struct {
		Description string `json:"description"`
		Focus       string `json:"focus"`
	}{
		Description: args.Description,
		Focus:       args.Focus,
	}

	data, err := json.Marshal(form)
	if err != nil {
		return fmt.Errorf("unable to marshal the form: %w", err)
	}

	requestBody := bytes.NewBuffer(data)

	params := requestParameters{
		httpMethod:  http.MethodPut,
		url:         g.auth.GetInstanceURL() + baseMediaPath + "/" + args.MediaAttachmentID,
		requestBody: requestBody,
		contentType: applicationJSON,
		output:      updated,
	}

	if err := g.sendRequest(params); err != nil {
		return fmt.Errorf(
			"received an error after sending the request to update the media attachment: %w",
			err,
		)
	}

	return nil
}

type DownloadMediaArgs struct {
	URL  string
	Path string
}

func (g *GTSClient) DownloadMedia(args DownloadMediaArgs, _ *NoRPCResults) error {
	ctx, cancel := context.WithTimeout(context.Background(), g.mediaTimeout)
	defer cancel()

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, args.URL, nil)
	if err != nil {
		return fmt.Errorf("unable to create the HTTP request: %w", err)
	}

	request.Header.Set("User-Agent", g.userAgent)

	response, err := g.httpClient.Do(request)
	if err != nil {
		return fmt.Errorf("received an error after attempting the download: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return BadStatusCodeError{
			statusCode: response.StatusCode,
		}
	}

	file, err := utilities.CreateFile(args.Path)
	if err != nil {
		return fmt.Errorf("unable to create %s: %w", args.Path, err)
	}
	defer file.Close()

	if _, err = io.Copy(file, response.Body); err != nil {
		return fmt.Errorf("unable to save the download to %s: %w", args.Path, err)
	}

	return nil
}

func (g *GTSClient) GetInstanceURL(_ NoRPCArgs, url *string) error {
	*url = g.auth.GetInstanceURL()

	return nil
}
