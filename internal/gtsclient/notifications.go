package gtsclient

import (
	"fmt"
	"net/http"

	"codeflow.dananglin.me.uk/apollo/enbas/internal/model"
)

const baseNotificationsPath string = "/api/v1/notifications"

func (g *GTSClient) GetNotification(notificationID string, notification *model.Notification) error {
	params := requestParameters{
		httpMethod:  http.MethodGet,
		url:         g.authentication.Instance + baseNotificationsPath + "/" + notificationID,
		requestBody: nil,
		contentType: "",
		output:      notification,
	}

	if err := g.sendRequest(params); err != nil {
		return fmt.Errorf(
			"received an error after sending the request to get the notification: %w",
			err,
		)
	}

	return nil
}

type GetNotificationListArgs struct {
	Limit        int
	IncludeTypes []string
	ExcludeTypes []string
}

func (g *GTSClient) GetNotificationList(args GetNotificationListArgs, notifications *[]model.Notification) error {
	query := fmt.Sprintf("?limit=%d", args.Limit)

	for _, include := range args.IncludeTypes {
		query = query + "&types[]=" + include
	}

	for _, exclude := range args.ExcludeTypes {
		query = query + "&exclude_types[]=" + exclude
	}

	params := requestParameters{
		httpMethod:  http.MethodGet,
		url:         g.authentication.Instance + baseNotificationsPath + query,
		requestBody: nil,
		contentType: "",
		output:      &notifications,
	}

	if err := g.sendRequest(params); err != nil {
		return fmt.Errorf(
			"received an error after sending the request to get the list of notifications: %w",
			err,
		)
	}

	return nil
}

func (g *GTSClient) DeleteNotifications(_ NoRPCArgs, _ *NoRPCResults) error {
	params := requestParameters{
		httpMethod:  http.MethodPost,
		url:         g.authentication.Instance + baseNotificationsPath + "/clear",
		requestBody: nil,
		contentType: "",
		output:      nil,
	}

	if err := g.sendRequest(params); err != nil {
		return fmt.Errorf(
			"received an error after sending the request to delete the notifications: %w",
			err,
		)
	}

	return nil
}
