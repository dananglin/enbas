package client

import (
	"fmt"
	"net/http"

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

//type CreateMediaAttachmentForm struct {
//	Description string
//	Focus       string
//	Filepath    string
//}
//
//func (g *Client) CreateMediaAttachment(form CreateMediaAttachmentForm) (model.Attachment, error) {
//	return model.Attachment{}, nil
//}
