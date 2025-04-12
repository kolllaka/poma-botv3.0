package model

import "encoding/json"

const (
	NOTIFICATION_SUBSCRIBE   = "subscribe"
	NOTIFICATION_RAID        = "raid"
	NOTIFICATION_FOLLOW      = "follow"
	NOTIFICATION_SUBGIFT     = "subgift"
	NOTIFICATION_CHEER       = "cheer"
	NOTIFICATION_RESUBSCRIBE = "resubscribe"

	NOTIFICATION_NAME = "notification"
)

type NotificationMessage struct {
	RouteType string

	Data json.RawMessage
}
