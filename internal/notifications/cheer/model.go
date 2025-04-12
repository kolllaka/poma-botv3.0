package cheer

type conf struct {
	conditions `json:"conditions"`
	fields     `json:"fields"`
}

type conditions struct {
	IsAnonym bool `json:"is_anonym,omitempty"`
	Bits     int  `json:"tier,omitempty"`
}

type fields struct {
	Title string `json:"title"`
	Path  string `json:"path"`
}

type message struct {
	Title   string `json:"title"`
	Link    string `json:"link"`
	Message string `json:"message"`
}

type cheer struct {
	IsAnonymous          bool   `json:"is_anonymous,omitempty"`
	UserID               string `json:"user_id,omitempty"`
	UserLogin            string `json:"user_login,omitempty"`
	UserName             string `json:"user_name,omitempty"`
	BroadcasterUserID    string `json:"broadcaster_user_id,omitempty"`
	BroadcasterUserLogin string `json:"broadcaster_user_login,omitempty"`
	BroadcasterUserName  string `json:"broadcaster_user_name,omitempty"`
	Message              string `json:"message,omitempty"`
	Bits                 int    `json:"bits,omitempty"`
}
