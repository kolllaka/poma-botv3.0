package subscribe

type conf struct {
	conditions `json:"conditions"`
	fields     `json:"fields"`
}

type conditions struct {
	IsGift bool `json:"is_gift,omitempty"`
	Tier   int  `json:"tier,omitempty"`
}

type fields struct {
	Title string `json:"title"`
	Path  string `json:"path"`
}

type message struct {
	Title string `json:"title"`
	Link  string `json:"link"`
}

type subscribe struct {
	UserID               string `json:"user_id,omitempty"`
	UserLogin            string `json:"user_login,omitempty"`
	UserName             string `json:"user_name,omitempty"`
	BroadcasterUserID    string `json:"broadcaster_user_id,omitempty"`
	BroadcasterUserLogin string `json:"broadcaster_user_login,omitempty"`
	BroadcasterUserName  string `json:"broadcaster_user_name,omitempty"`
	Tier                 int    `json:"tier,omitempty"`
	IsGift               bool   `json:"is_gift,omitempty"`
}
