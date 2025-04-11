package model

import "encoding/json"

const (
	NOTIFICATION_SUBSCRIBE = "subscribe"

	NOTIFICATION_NAME = "notification"
)

type NotificationMessage struct {
	RouteType string

	Data json.RawMessage
}
