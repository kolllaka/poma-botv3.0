package raid

type conf struct {
	conditions `json:"conditions"`
	fields     `json:"fields"`
}

type conditions struct {
	Viewers int `json:"viewers"`
}

type fields struct {
	Title string `json:"title"`
	Path  string `json:"path"`
}

type message struct {
	Title string `json:"title"`
	Link  string `json:"link"`
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
