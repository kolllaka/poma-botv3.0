package rewards

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/kolllaka/poma-botv3.0/internal/model"
	"github.com/kolllaka/poma-botv3.0/internal/rewards/aug"
	"github.com/kolllaka/poma-botv3.0/internal/rewards/music"
	"github.com/kolllaka/poma-botv3.0/internal/services"
	"github.com/kolllaka/poma-botv3.0/pkg/logging"
)

type Rewards interface {
	HandleReward()
	InitRewards(cfg *model.RewardsConfig)

	GetRewardChannel(rewardType string) chan []byte
	GetPathToUrl() []string
}
type rewards struct {
	logger   *logging.Logger
	services services.Service

	routes map[string]Route
	reader chan model.RewardMessage

	WritersChan map[string]chan []byte

	pathToUrl []string
}

func New(
	logger *logging.Logger,
	services services.Service,
	reader chan model.RewardMessage,
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
	go func(reader chan model.RewardMessage) {
		for {
			reward := <-reader

			rewardName := strings.ToLower(reward.Name)

			r.logger.Debug("reward received", logging.AnyAttr("reward", reward))

			route, ok := r.routes[rewardName]
			if !ok {
				r.logger.Warn("unknown reward route", logging.AnyAttr("reward", reward))

				continue
			}

			go func() {
				rType, rBody, err := route.RunRoute(reward)
				if err != nil {
					r.logger.Warn(
						"RunRoute",
						logging.StringAttr("reward type", rType),
						logging.AnyAttr("body", string(rBody)),
						logging.ErrAttr(err),
					)

					return
				}

				r.WritersChan[rType] <- rBody
			}()
		}
	}(r.reader)
}

func (r *rewards) InitRewards(cfg *model.RewardsConfig) {
	for _, reward := range cfg.Rewards {
		rewardName := strings.ToLower(reward.Name)

		switch reward.Type {
		case model.AUGURY:
			fields := r.appendPathToUrl(reward.Fields)
			r.routes[rewardName] = aug.NewRoute(model.AUGURY, fields)
		case model.MUSIC:
			r.routes[rewardName] = music.NewRoute(r.services, model.MUSIC, reward.Fields)
		default:
			r.logger.Error("unknown reward type", logging.AnyAttr("reward", reward))

			continue
		}

		if _, ok := r.WritersChan[reward.Type]; !ok {
			r.WritersChan[reward.Type] = make(chan []byte)
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

func (r *rewards) appendPathToUrl(fields json.RawMessage) json.RawMessage {
	type field struct {
		Path string `json:"path"`
	}
	type newField struct {
		F   json.RawMessage `json:"fields"`
		Url string          `json:"url"`
	}
	f := field{}
	json.Unmarshal(fields, &f)

	r.pathToUrl = append(r.pathToUrl, f.Path)

	newFields, _ := json.Marshal(newField{
		F:   fields,
		Url: fmt.Sprintf("/%s%d/", model.REWARD_NAME, len(r.pathToUrl)-1),
	})

	return newFields
}
