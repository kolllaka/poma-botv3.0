package notifications

import "github.com/kolllaka/poma-botv3.0/internal/model"

type Route interface {
	RunRoute(msg model.NotificationMessage) (string, []byte, error)
}
