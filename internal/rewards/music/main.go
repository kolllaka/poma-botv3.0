package music

import (
	"bytes"
	"encoding/json"
	"fmt"

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
	var music Music
	if err := json.Unmarshal(msg.Data, &music); err != nil {
		return r.rewardType, nil, fmt.Errorf("error on RunRoute: %w", err)
	}

	bSong := r.services.GetYoutubeMusicBy(music.Text, true)
	duration := bSong.Data.(model.Music).Duration

	var network bytes.Buffer
	json.NewEncoder(&network).Encode(bSong)

	if duration < 0 || r.conf.Duration > duration {
		return r.rewardType, network.Bytes(), nil
	}

	return r.rewardType, network.Bytes(), getErrorRequestToLong(duration, r.conf.Duration)
}
