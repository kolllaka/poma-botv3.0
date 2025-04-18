package notifications

import (
	"strings"

	"github.com/kolllaka/poma-botv3.0/internal/model"
	"github.com/kolllaka/poma-botv3.0/internal/notifications/cheer"
	"github.com/kolllaka/poma-botv3.0/internal/notifications/follow"
	"github.com/kolllaka/poma-botv3.0/internal/notifications/raid"
	"github.com/kolllaka/poma-botv3.0/internal/notifications/resubscribe"
	"github.com/kolllaka/poma-botv3.0/internal/notifications/subgift"
	"github.com/kolllaka/poma-botv3.0/internal/notifications/subscribe"
	"github.com/kolllaka/poma-botv3.0/internal/services"
	"github.com/kolllaka/poma-botv3.0/pkg/logging"
)

type Notifications interface {
	HandleNotification()
	InitNotifications(cfg *model.NotificationsConfig)

	GetNotificationChannel(notificationType string) chan []byte
}
type notifications struct {
	logger   *logging.Logger
	services services.Service

	routes map[string]Route
	reader chan model.NotificationMessage

	WritersChan map[string]chan []byte
	ErrorChan   chan error
}

func New(
	logger *logging.Logger,
	services services.Service,
	reader chan model.NotificationMessage,
) Notifications {
	return &notifications{
		logger:      logger,
		services:    services,
		routes:      make(map[string]Route),
		reader:      reader,
		WritersChan: make(map[string]chan []byte),
		ErrorChan:   make(chan error),
	}
}

func (r *notifications) HandleNotification() {
	go func(reader chan model.NotificationMessage) {
		for {
			notification := <-reader

			notificationsType := strings.ToLower(notification.RouteType)

			r.logger.Debug("notification received", logging.AnyAttr("notification", notification))

			route, ok := r.routes[notificationsType]
			if !ok {
				r.logger.Warn("unknown notification route", logging.AnyAttr("notification", notification))

				continue
			}

			nType, rBody, err := route.RunRoute(notification)
			if err != nil {
				r.logger.Warn(
					"RunRoute",
					logging.StringAttr("notification type", nType),
					logging.AnyAttr("body", string(rBody)),
					logging.ErrAttr(err),
				)

				return
			}

			r.logger.Debug("rBody received",
				logging.AnyAttr("rBody", rBody),
				logging.StringAttr("nType", nType),
				logging.AnyAttr("chan", r.WritersChan[nType]),
			)

			r.WritersChan[nType] <- rBody
		}
	}(r.reader)
}

func (r *notifications) InitNotifications(cfg *model.NotificationsConfig) {
	for _, notification := range cfg.Notifications {
		notificationType := strings.ToLower(notification.Type)

		switch notificationType {
		case model.NOTIFICATION_SUBSCRIBE:
			r.routes[notificationType] = subscribe.NewRoute(notificationType, notification.Checks)
		case model.NOTIFICATION_RAID:
			r.routes[notificationType] = raid.NewRoute(notificationType, notification.Checks)
		case model.NOTIFICATION_FOLLOW:
			r.routes[notificationType] = follow.NewRoute(notificationType, notification.Checks)
		case model.NOTIFICATION_SUBGIFT:
			r.routes[notificationType] = subgift.NewRoute(notificationType, notification.Checks)
		case model.NOTIFICATION_CHEER:
			r.routes[notificationType] = cheer.NewRoute(notificationType, notification.Checks)
		case model.NOTIFICATION_RESUBSCRIBE:
			r.routes[notificationType] = resubscribe.NewRoute(notificationType, notification.Checks)

		default:
			r.logger.Warn("unknown notification type", logging.AnyAttr("notification", notification))

			continue
		}

		if _, ok := r.WritersChan[notificationType]; !ok {
			r.WritersChan[notificationType] = make(chan []byte)
		}
	}

	r.logger.Debug("notifications initialized", logging.AnyAttr("notifications", r.routes), logging.AnyAttr("channels", r.WritersChan))
}

// GetNotificationChannel implements Notifications.
func (r *notifications) GetNotificationChannel(notificationType string) chan []byte {
	return r.WritersChan[notificationType]
}
