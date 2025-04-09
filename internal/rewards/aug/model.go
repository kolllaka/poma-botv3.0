package aug

type conf struct {
	fields `json:"fields"`
	Url    string `json:"url"`
}

type fields struct {
	Title string `json:"title"`
	Path  string `json:"path"`
}

type Message struct {
	Title string `json:"title"`
	Link  string `json:"link"`
}

type RouteMessage struct {
	UserName string
}
