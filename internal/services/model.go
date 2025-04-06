package services

const (
	FILE = iota
	YOUTUBE
)

type Responce struct {
	Source   int  `json:"source"`
	IsReward bool `json:"is_reward"`
	Data     any  `json:"data"`
}
