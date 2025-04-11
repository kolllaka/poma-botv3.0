package aug

type conf struct {
	Title string `json:"title"`
	Path  string `json:"path"`
}

type message struct {
	Title string `json:"title"`
	Link  string `json:"link"`
}

type RouteMessage struct {
	UserName string
}
