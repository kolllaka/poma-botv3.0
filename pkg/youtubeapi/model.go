package youtubeapi

type Response struct {
	StatusCode int
	Data       []byte
}

type SongResp struct {
	Items []struct {
		Id      string `json:"id"`
		Snippet struct {
			Title string `json:"title,omitempty"`
		} `json:"snippet,omitempty"`
		ContentDetails struct {
			Duration string `json:"duration,omitempty"`
		} `json:"contentDetails,omitempty"`
	} `json:"items,omitempty"`
}
