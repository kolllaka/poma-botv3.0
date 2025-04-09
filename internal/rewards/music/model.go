package music

type conf struct {
	Title    string `json:"title"`
	Duration int    `json:"duration"`
}

type Music struct {
	IsReward bool
	Source   string

	Author string
	Text   string
}
