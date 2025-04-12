package subgift

type conf struct {
	conditions `json:"conditions"`
	fields     `json:"fields"`
}

type conditions struct {
	IsAnonym bool `json:"is_anonym,omitempty"`
	Tier     int  `json:"tier,omitempty"`
}

type fields struct {
	Title string `json:"title"`
	Path  string `json:"path"`
}

type message struct {
	Title string `json:"title"`
	Link  string `json:"link"`
}

type subgift struct {
	UserID               string `json:"user_id,omitempty"`
	UserLogin            string `json:"user_login,omitempty"`
	UserName             string `json:"user_name,omitempty"`
	BroadcasterUserID    string `json:"broadcaster_user_id,omitempty"`
	BroadcasterUserLogin string `json:"broadcaster_user_login,omitempty"`
	BroadcasterUserName  string `json:"broadcaster_user_name,omitempty"`
	Total                int    `json:"total,omitempty"`
	Tier                 int    `json:"tier,omitempty"`
	CumulativeTotal      int    `json:"cumulative_total,omitempty"`
	IsAnonymous          bool   `json:"is_anonymous,omitempty"`
}
