// SPDX-FileCopyrightText: 2024 Dan Anglin <d.n.i.anglin@gmail.com>
//
// SPDX-License-Identifier: GPL-3.0-or-later

package client

import (
	"fmt"
	"net/http"

	"codeflow.dananglin.me.uk/apollo/enbas/internal/model"
)

func (g *Client) GetHomeTimeline(limit int) (model.Timeline, error) {
	path := fmt.Sprintf("/api/v1/timelines/home?limit=%d", limit)

	timeline := model.Timeline{
		Name:     "HOME TIMELINE",
		Statuses: nil,
	}

	return g.getTimeline(path, timeline)
}

func (g *Client) GetPublicTimeline(limit int) (model.Timeline, error) {
	path := fmt.Sprintf("/api/v1/timelines/public?limit=%d", limit)

	timeline := model.Timeline{
		Name:     "PUBLIC TIMELINE",
		Statuses: nil,
	}

	return g.getTimeline(path, timeline)
}

func (g *Client) GetListTimeline(listID string, limit int) (model.Timeline, error) {
	path := fmt.Sprintf("/api/v1/timelines/list/%s?limit=%d", listID, limit)

	timeline := model.Timeline{
		Name:     "LIST: " + listID,
		Statuses: nil,
	}

	return g.getTimeline(path, timeline)
}

func (g *Client) GetTagTimeline(tag string, limit int) (model.Timeline, error) {
	path := fmt.Sprintf("/api/v1/timelines/tag/%s?limit=%d", tag, limit)

	timeline := model.Timeline{
		Name:     "TAG: " + tag,
		Statuses: nil,
	}

	return g.getTimeline(path, timeline)
}

func (g *Client) getTimeline(path string, timeline model.Timeline) (model.Timeline, error) {
	url := g.Authentication.Instance + path

	var statuses []model.Status

	if err := g.sendRequest(http.MethodGet, url, nil, &statuses); err != nil {
		return timeline, fmt.Errorf("received an error after sending the request to get the timeline; %w", err)
	}

	timeline.Statuses = statuses

	return timeline, nil
}
