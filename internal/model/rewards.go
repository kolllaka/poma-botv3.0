package model

import "encoding/json"

const (
	AUGURY       = "augury"
	MUSIC        = "music"
	NOTIFICATION = "notification"

	REWARD_NAME = "reward"
)

type RewardMessage struct {
	RouteName string

	Data json.RawMessage
}
