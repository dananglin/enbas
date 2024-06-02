// SPDX-FileCopyrightText: 2024 Dan Anglin <d.n.i.anglin@gmail.com>
//
// SPDX-License-Identifier: GPL-3.0-or-later

package client

import (
	"fmt"
	"net/http"

	"codeflow.dananglin.me.uk/apollo/enbas/internal/model"
)

func (g *Client) GetInstance() (model.InstanceV2, error) {
	path := "/api/v2/instance"
	url := g.Authentication.Instance + path

	var instance model.InstanceV2

	if err := g.sendRequest(http.MethodGet, url, nil, &instance); err != nil {
		return model.InstanceV2{}, fmt.Errorf("received an error after sending the request to get the instance details; %w", err)
	}

	return instance, nil
}
