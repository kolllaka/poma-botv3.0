package notification

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
	Link  string `json:"path"`
}

type raid struct {
	FromBroadcasterUserID    string `json:"from_broadcaster_user_id,omitempty"`
	FromBroadcasterUserLogin string `json:"from_broadcaster_user_login,omitempty"`
	FromBroadcasterUserName  string `json:"from_broadcaster_user_name,omitempty"`
	ToBroadcasterUserID      string `json:"to_broadcaster_user_id,omitempty"`
	ToBroadcasterUserLogin   string `json:"to_broadcaster_user_login,omitempty"`
	ToBroadcasterUserName    string `json:"to_broadcaster_user_name,omitempty"`
	Viewers                  int    `json:"viewers,omitempty"`
}
