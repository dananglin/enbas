// SPDX-FileCopyrightText: 2024 Dan Anglin <d.n.i.anglin@gmail.com>
//
// SPDX-License-Identifier: GPL-3.0-or-later

package client

import (
	"fmt"
	"net/http"

	"codeflow.dananglin.me.uk/apollo/enbas/internal/model"
)

func (g *Client) GetHomeTimeline(limit int) (model.StatusList, error) {
	path := fmt.Sprintf("/api/v1/timelines/home?limit=%d", limit)

	timeline := model.StatusList{
		Name:     "Timeline: Home",
		Statuses: nil,
	}

	return g.getTimeline(path, timeline)
}

func (g *Client) GetPublicTimeline(limit int) (model.StatusList, error) {
	path := fmt.Sprintf("/api/v1/timelines/public?limit=%d", limit)

	timeline := model.StatusList{
		Name:     "Timeline: Public",
		Statuses: nil,
	}

	return g.getTimeline(path, timeline)
}

func (g *Client) GetListTimeline(listID, title string, limit int) (model.StatusList, error) {
	path := fmt.Sprintf("/api/v1/timelines/list/%s?limit=%d", listID, limit)

	timeline := model.StatusList{
		Name:     "Timeline: List (" + title + ")",
		Statuses: nil,
	}

	return g.getTimeline(path, timeline)
}

func (g *Client) GetTagTimeline(tag string, limit int) (model.StatusList, error) {
	path := fmt.Sprintf("/api/v1/timelines/tag/%s?limit=%d", tag, limit)

	timeline := model.StatusList{
		Name:     "Timeline: Tag (" + tag + ")",
		Statuses: nil,
	}

	return g.getTimeline(path, timeline)
}

func (g *Client) getTimeline(path string, timeline model.StatusList) (model.StatusList, error) {
	url := g.Authentication.Instance + path

	var statuses []model.Status

	if err := g.sendRequest(http.MethodGet, url, nil, &statuses); err != nil {
		return timeline, fmt.Errorf("received an error after sending the request to get the timeline: %w", err)
	}

	timeline.Statuses = statuses

	return timeline, nil
}
