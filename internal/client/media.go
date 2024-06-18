package client

import (
	"fmt"
	"net/http"

	"codeflow.dananglin.me.uk/apollo/enbas/internal/model"
)

func (g *Client) GetMediaAttachment(mediaAttachmentID string) (model.Attachment, error) {
	url := g.Authentication.Instance + "/api/v1/media/" + mediaAttachmentID

	var attachment model.Attachment

	if err := g.sendRequest(http.MethodGet, url, nil, &attachment); err != nil {
		return model.Attachment{}, fmt.Errorf("received an error after sending the request to get the media attachment: %w", err)
	}

	return attachment, nil
}
