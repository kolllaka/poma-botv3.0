package notifications

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/kolllaka/poma-botv3.0/internal/model"
	"github.com/kolllaka/poma-botv3.0/internal/notifications/subscribe"
	"github.com/kolllaka/poma-botv3.0/internal/services"
	"github.com/kolllaka/poma-botv3.0/pkg/logging"
)

type Rewards interface {
	HandleReward()
	InitRewards(cfg *model.NotificationsConfig)

	GetRewardChannel(rewardType string) chan []byte
	GetPathToUrl() []string
}
type rewards struct {
	logger   *logging.Logger
	services services.Service

	routes map[string]Route
	reader chan model.NotificationMessage

	WritersChan map[string]chan []byte

	pathToUrl []string
}

func New(
	logger *logging.Logger,
	services services.Service,
	reader chan model.NotificationMessage,
) Rewards {
	return &rewards{
		logger:      logger,
		services:    services,
		routes:      make(map[string]Route),
		reader:      reader,
		WritersChan: make(map[string]chan []byte),
		pathToUrl:   []string{},
	}
}

func (r *rewards) HandleReward() {
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

			go func() {
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

				r.WritersChan[nType] <- rBody
			}()
		}
	}(r.reader)
}

func (r *rewards) InitRewards(cfg *model.NotificationsConfig) {
	for _, notification := range cfg.Notifications {
		notificationType := strings.ToLower(notification.Type)

		switch notificationType {
		case model.NOTIFICATION_SUBSCRIBE:
			r.routes[notificationType] = subscribe.NewRoute(notificationType, notification.Checks)

		default:
			r.logger.Warn("unknown notification type", logging.AnyAttr("notification", notification))

			continue
		}

		if _, ok := r.WritersChan[notification.Type]; !ok {
			r.WritersChan[notification.Type] = make(chan []byte)
		}
	}

	r.logger.Debug("rewards initialized", logging.AnyAttr("rewards", r.routes), logging.AnyAttr("channels", r.WritersChan))
}

// GetRewardChannel implements Rewards.
func (r *rewards) GetRewardChannel(rewardType string) chan []byte {
	return r.WritersChan[rewardType]
}

func (r *rewards) GetPathToUrl() []string {
	return r.pathToUrl
}

func (r *rewards) appendPathToUrl(checks json.RawMessage) string {
	type path struct {
		Path string `json:"path"`
	}
	f := path{}
	json.Unmarshal(checks, &f)

	r.pathToUrl = append(r.pathToUrl, f.Path)

	return fmt.Sprintf("/%s%d/", model.NOTIFICATION_NAME, len(r.pathToUrl)-1)
}
