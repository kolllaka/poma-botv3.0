package model

const (
	AUGURY = "augury"
	MUSIC  = "music"

	REWARD_NAME = "reward"
)

type RewardMessage struct {
	IsReward bool
	Name     string
	Username string
	Text     string
}
