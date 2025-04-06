package model

import (
	"encoding/json"
)

type EnvConfig struct {
	Host           string `env:"HOST" env-default:"localhost"`
	Port           int    `env:"PORT" env-default:"8080"`
	UserId         int    `env:"USERID" env-required:"true"`
	YoutubeKey     string `env:"YOUTUBEKEY" env-required:"true"`
	MyPlaylistPath string `env:"MYPLAYLISTPATH" env-required:"true"`
}

func NewEnvConfig() *EnvConfig {
	return &EnvConfig{}
}

type RewardsConfig struct {
	Rewards []RewardConfig `json:"rewards"`
}

type RewardConfig struct {
	Type   string          `json:"type"`
	Name   string          `json:"name"`
	Fields json.RawMessage `json:"fields"`
}

func NewRewardConfig() *RewardsConfig {
	return &RewardsConfig{}
}
