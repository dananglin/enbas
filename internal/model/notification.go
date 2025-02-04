package model

import (
	"encoding/json"
	"fmt"
	"time"
)

type Notification struct {
	Account   *Account         `json:"account"`
	CreatedAt time.Time        `json:"created_at"`
	ID        string           `json:"id"`
	Status    *Status          `json:"status"`
	Type      NotificationType `json:"type"`
}

type NotificationType int

const (
	NotificationTypeFollow = iota
	NotificationTypeFollowRequest
	NotificationTypeMention
	NotificationTypeReblog
	NotificationTypeFavourite
	NotificationTypePoll
	NotificationTypeStatus
	NotificationTypeUnknown
)

const (
	notificationTypeFollowValue        = "follow"
	notificationTypeFollowRequestValue = "follow_request"
	notificationTypeMentionValue       = "mention"
	notificationTypeReblogValue        = "reblog"
	notificationTypeFavouriteValue     = "favourite"
	notificationTypePollValue          = "poll"
	notificationTypeStatusValue        = "status"
)

func (n NotificationType) String() string {
	mapped := map[NotificationType]string{
		NotificationTypeFollow:        notificationTypeFollowValue,
		NotificationTypeFollowRequest: notificationTypeFollowRequestValue,
		NotificationTypeMention:       notificationTypeMentionValue,
		NotificationTypeReblog:        notificationTypeReblogValue,
		NotificationTypeFavourite:     notificationTypeFavouriteValue,
		NotificationTypePoll:          notificationTypePollValue,
		NotificationTypeStatus:        notificationTypeStatusValue,
	}

	output, ok := mapped[n]
	if !ok {
		return unknownValue
	}

	return output
}

func ParseNotificationType(value string) (NotificationType, error) {
	mapped := map[string]NotificationType{
		notificationTypeFollowValue:        NotificationTypeFollow,
		notificationTypeFollowRequestValue: NotificationTypeFollowRequest,
		notificationTypeMentionValue:       NotificationTypeMention,
		notificationTypeReblogValue:        NotificationTypeReblog,
		notificationTypeFavouriteValue:     NotificationTypeFavourite,
		notificationTypePollValue:          NotificationTypePoll,
		notificationTypeStatusValue:        NotificationTypeStatus,
	}

	output, ok := mapped[value]
	if !ok {
		return NotificationTypeUnknown, InvalidNotificationTypeError{value: value}
	}

	return output, nil
}

func (n *NotificationType) UnmarshalJSON(data []byte) error {
	var (
		value string
		err   error
	)

	if err = json.Unmarshal(data, &value); err != nil {
		return fmt.Errorf("unable to decode the JSON data: %w", err)
	}

	*n, _ = ParseNotificationType(value)

	return nil
}

type InvalidNotificationTypeError struct {
	value string
}

func (e InvalidNotificationTypeError) Error() string {
	return "'" + e.value + "is not a valid notification type: valid values are " +
		notificationTypeFollowValue + ", " +
		notificationTypeFollowRequestValue + ", " +
		notificationTypeMentionValue + ", " +
		notificationTypeReblogValue + ", " +
		notificationTypeFavouriteValue + ", " +
		notificationTypePollValue + ", " +
		notificationTypeStatusValue
}
