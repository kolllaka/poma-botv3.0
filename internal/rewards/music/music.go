package music

import (
	"encoding/json"

	"github.com/kolllaka/poma-botv3.0/internal/model"
	"github.com/kolllaka/poma-botv3.0/internal/services"
)

type route struct {
	services   services.Service
	rewardType string
	conf       conf
}

func NewRoute(services services.Service, rewardType string, rawConf json.RawMessage) *route {
	var conf conf
	json.Unmarshal(rawConf, &conf)

	return &route{
		services:   services,
		rewardType: rewardType,
		conf:       conf,
	}
}

// RunRoute implements Route.
func (r *route) RunRoute(msg model.RewardMessage) (string, []byte, error) {
	bSong := r.services.GetYoutubeMusicBy(msg)

	var resp Resp

	json.Unmarshal(bSong, &resp)

	if resp.Data.Duration < 0 || r.conf.Duration > resp.Data.Duration {
		return r.rewardType, bSong, nil
	}

	return r.rewardType, bSong, getErrorRequestToLong(resp.Data.Duration, r.conf.Duration)
}
