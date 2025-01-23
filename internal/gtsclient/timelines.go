package gtsclient

import (
	"fmt"
	"net/http"

	"codeflow.dananglin.me.uk/apollo/enbas/internal/model"
)

func (g *GTSClient) GetHomeTimeline(limit int, timeline *model.StatusList) error {
	path := fmt.Sprintf("/api/v1/timelines/home?limit=%d", limit)

	*timeline = model.StatusList{
		Name:     "Timeline: Home",
		Statuses: nil,
	}

	return g.getTimeline(path, timeline)
}

func (g *GTSClient) GetPublicTimeline(limit int, timeline *model.StatusList) error {
	path := fmt.Sprintf("/api/v1/timelines/public?limit=%d", limit)

	*timeline = model.StatusList{
		Name:     "Timeline: Public",
		Statuses: nil,
	}

	return g.getTimeline(path, timeline)
}

type GetListTimelineArgs struct {
	ListID string
	Title  string
	Limit  int
}

func (g *GTSClient) GetListTimeline(args GetListTimelineArgs, timeline *model.StatusList) error {
	path := fmt.Sprintf("/api/v1/timelines/list/%s?limit=%d", args.ListID, args.Limit)

	*timeline = model.StatusList{
		Name:     "Timeline: List (" + args.Title + ")",
		Statuses: nil,
	}

	return g.getTimeline(path, timeline)
}

type GetTagTimelineArgs struct {
	TagName string
	Limit   int
}

func (g *GTSClient) GetTagTimeline(args GetTagTimelineArgs, timeline *model.StatusList) error {
	path := fmt.Sprintf("/api/v1/timelines/tag/%s?limit=%d", args.TagName, args.Limit)

	*timeline = model.StatusList{
		Name:     "Timeline: Tag (" + args.TagName + ")",
		Statuses: nil,
	}

	return g.getTimeline(path, timeline)
}

func (g *GTSClient) getTimeline(path string, timeline *model.StatusList) error {
	var statuses []model.Status

	params := requestParameters{
		httpMethod:  http.MethodGet,
		url:         g.Authentication.Instance + path,
		requestBody: nil,
		contentType: "",
		output:      &statuses,
	}

	if err := g.sendRequest(params); err != nil {
		return fmt.Errorf("received an error after sending the request to get the timeline: %w", err)
	}

	timeline.Statuses = statuses

	return nil
}
