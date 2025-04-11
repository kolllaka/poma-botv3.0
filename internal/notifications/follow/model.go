package follow

import "time"

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

type follow struct {
	UserID               string    `json:"user_id,omitempty"`
	UserLogin            string    `json:"user_login,omitempty"`
	UserName             string    `json:"user_name,omitempty"`
	BroadcasterUserID    string    `json:"broadcaster_user_id,omitempty"`
	BroadcasterUserLogin string    `json:"broadcaster_user_login,omitempty"`
	BroadcasterUserName  string    `json:"broadcaster_user_name,omitempty"`
	FollowedAt           time.Time `json:"followed_at,omitempty"`
}
