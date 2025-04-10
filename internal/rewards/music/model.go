package music

type conf struct {
	Title    string `json:"title"`
	Duration int    `json:"duration"`
}

type Music struct {
	IsReward bool `json:"is_reward"`

	Author string `json:"username"`
	Text   string `json:"text"`
}
