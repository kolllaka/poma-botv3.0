package music

type conf struct {
	Title    string `json:"title"`
	Duration int    `json:"duration"`
}

type Resp struct {
	Data struct {
		Duration int `json:"duration"`
	} `json:"data"`
}
