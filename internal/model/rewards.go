package model

import "encoding/json"

const (
	REWARD_AUGURY = "augury"
	REWARD_MUSIC  = "music"

	REWARD_NAME = "reward"
)

type RewardMessage struct {
	RouteName string

	Data json.RawMessage
}
