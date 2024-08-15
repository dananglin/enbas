package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"codeflow.dananglin.me.uk/apollo/enbas/internal/model"
)

const (
	baseMediaPath string = "/api/v1/media"
)

func (g *Client) GetMediaAttachment(mediaAttachmentID string) (model.Attachment, error) {
	url := g.Authentication.Instance + baseMediaPath + "/" + mediaAttachmentID

	var attachment model.Attachment

	params := requestParameters{
		httpMethod:  http.MethodGet,
		url:         url,
		requestBody: nil,
		contentType: "",
		output:      &attachment,
	}

	if err := g.sendRequest(params); err != nil {
		return model.Attachment{}, fmt.Errorf("received an error after sending the request to get the media attachment: %w", err)
	}

	return attachment, nil
}

func (g *Client) CreateMediaAttachment(path, description, focus string) (model.Attachment, error) {
	file, err := os.Open(path)
	if err != nil {
		return model.Attachment{}, fmt.Errorf("unable to open the file: %w", err)
	}
	defer file.Close()

	// create the request body using a writer from the multipart package
	requestBody := bytes.Buffer{}
	requestBodyWriter := multipart.NewWriter(&requestBody)

	filename := filepath.Base(path)

	part, err := requestBodyWriter.CreateFormFile("file", filename)
	if err != nil {
		return model.Attachment{}, fmt.Errorf("unable to create the new part: %w", err)
	}

	if _, err := io.Copy(part, file); err != nil {
		return model.Attachment{}, fmt.Errorf("unable to copy the file contents to the form: %w", err)
	}

	// add the description
	if description != "" {
		descriptionFormFieldWriter, err := requestBodyWriter.CreateFormField("description")
		if err != nil {
			return model.Attachment{}, fmt.Errorf(
				"unable to create the writer for the 'description' form field: %w",
				err,
			)
		}

		if _, err := io.WriteString(descriptionFormFieldWriter, description); err != nil {
			return model.Attachment{}, fmt.Errorf(
				"unable to write the description to the form: %w",
				err,
			)
		}
	}

	// add the focus values
	if focus != "" {
		focusFormFieldWriter, err := requestBodyWriter.CreateFormField("focus")
		if err != nil {
			return model.Attachment{}, fmt.Errorf(
				"unable to create the writer for the 'focus' form field: %w",
				err,
			)
		}

		if _, err := io.WriteString(focusFormFieldWriter, focus); err != nil {
			return model.Attachment{}, fmt.Errorf(
				"unable to write the focus values to the form: %w",
				err,
			)
		}
	}

	if err := requestBodyWriter.Close(); err != nil {
		return model.Attachment{}, fmt.Errorf("unable to close the writer: %w", err)
	}

	url := g.Authentication.Instance + baseMediaPath

	var attachment model.Attachment

	params := requestParameters{
		httpMethod:  http.MethodPost,
		url:         url,
		requestBody: &requestBody,
		contentType: requestBodyWriter.FormDataContentType(),
		output:      &attachment,
	}

	if err := g.sendRequest(params); err != nil {
		return model.Attachment{}, fmt.Errorf(
			"received an error after sending the request to create the media attachment: %w",
			err,
		)
	}

	return attachment, nil
}

func (g *Client) UpdateMediaAttachment(mediaAttachmentID, description, focus string) (model.Attachment, error) {
	form := struct {
		Description string `json:"description"`
		Focus       string `json:"focus"`
	}{
		Description: description,
		Focus:       focus,
	}

	data, err := json.Marshal(form)
	if err != nil {
		return model.Attachment{}, fmt.Errorf("unable to marshal the form: %w", err)
	}

	requestBody := bytes.NewBuffer(data)
	url := g.Authentication.Instance + baseMediaPath + "/" + mediaAttachmentID

	var updatedMediaAttachment model.Attachment

	params := requestParameters{
		httpMethod:  http.MethodPut,
		url:         url,
		requestBody: requestBody,
		contentType: applicationJSON,
		output:      &updatedMediaAttachment,
	}

	if err := g.sendRequest(params); err != nil {
		return model.Attachment{}, fmt.Errorf(
			"received an error after sending the request to update the media attachment: %w",
			err,
		)
	}

	return updatedMediaAttachment, nil
}
