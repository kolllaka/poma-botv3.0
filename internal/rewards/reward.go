package rewards

import "github.com/kolllaka/poma-botv3.0/internal/model"

type Route interface {
	RunRoute(msg model.RewardMessage) (string, []byte, error)
}
