package music

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/kolllaka/poma-botv3.0/internal/model"
	"github.com/kolllaka/poma-botv3.0/internal/services"
)

type route struct {
	services   services.Service
	rewardType string
	conf       conf
}

func NewRoute(rewardType string, services services.Service, rawConf json.RawMessage) *route {
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
	var music music
	if err := json.Unmarshal(msg.Data, &music); err != nil {
		return r.rewardType, nil, fmt.Errorf("error on RunRoute: %w", err)
	}

	bSong, err := r.services.GetYoutubeMusicBy(model.Music{
		Link:   music.Text,
		Author: music.Author,
	}, music.IsReward)
	if err != nil {
		if errors.Is(err, model.ErrBadLink) {
			return r.rewardType, nil, err
		}

		return r.rewardType, nil, fmt.Errorf("error on GetYoutubeMusicBy: %w", err)
	}
	duration := bSong.Data.(model.Music).Duration

	var network bytes.Buffer
	json.NewEncoder(&network).Encode(bSong)

	if duration < 0 || r.conf.Duration > duration {
		return r.rewardType, network.Bytes(), nil
	}

	return r.rewardType, network.Bytes(), getErrorRequestToLong(duration, r.conf.Duration)
}
