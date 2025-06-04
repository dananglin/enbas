package gtsclient

import (
	"fmt"
	"net/http"

	"codeflow.dananglin.me.uk/apollo/enbas/internal/model"
)

const (
	baseTagsPath     = "/api/v1/tags"
	followedTagsPath = "/api/v1/followed_tags"
)

func (g *GTSClient) GetFollowedTags(limit int, list *model.TagList) error {
	var tags []model.Tag

	params := requestParameters{
		httpMethod:  http.MethodGet,
		url:         g.auth.GetInstanceURL() + fmt.Sprintf("%s?limit=%d", followedTagsPath, limit),
		requestBody: nil,
		contentType: "",
		output:      &tags,
	}

	if err := g.sendRequest(params); err != nil {
		return fmt.Errorf(
			"received an error after sending the request to get the list of followed tags: %w",
			err,
		)
	}

	*list = model.TagList{
		Name: "Followed tags:",
		Tags: tags,
	}

	return nil
}

func (g *GTSClient) GetTag(name string, tag *model.Tag) error {
	params := requestParameters{
		httpMethod:  http.MethodGet,
		url:         g.auth.GetInstanceURL() + baseTagsPath + "/" + name,
		contentType: "",
		requestBody: nil,
		output:      tag,
	}

	if err := g.sendRequest(params); err != nil {
		return fmt.Errorf(
			"received an error after sending the request to get the details of the tag: %w",
			err,
		)
	}

	return nil
}

func (g *GTSClient) FollowTag(name string, _ *NoRPCResults) error {
	params := requestParameters{
		httpMethod:  http.MethodPost,
		url:         g.auth.GetInstanceURL() + baseTagsPath + "/" + name + "/follow",
		contentType: "",
		requestBody: nil,
		output:      nil,
	}

	if err := g.sendRequest(params); err != nil {
		return fmt.Errorf(
			"received an error after sending the request to follow the tag %q: %w",
			name,
			err,
		)
	}

	return nil
}

func (g *GTSClient) UnfollowTag(name string, _ *NoRPCResults) error {
	params := requestParameters{
		httpMethod:  http.MethodPost,
		url:         g.auth.GetInstanceURL() + baseTagsPath + "/" + name + "/unfollow",
		contentType: "",
		requestBody: nil,
		output:      nil,
	}

	if err := g.sendRequest(params); err != nil {
		return fmt.Errorf(
			"received an error after sending the request to unfollow the tag %q: %w",
			name,
			err,
		)
	}

	return nil
}
