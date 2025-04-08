package model

const (
	AUGURY       = "augury"
	MUSIC        = "music"
	NOTIFICATION = "notification"

	REWARD_NAME = "reward"
)

type RewardMessage struct {
	Reward struct {
		IsReward bool
		Name     string
		Username string
		Text     string
	}
	Event struct {
	}
}
