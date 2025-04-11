package resubscribe

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
	Msg   string `json:"msg"`
}

type reSubscribe struct {
	UserID               string `json:"user_id,omitempty"`
	UserLogin            string `json:"user_login,omitempty"`
	UserName             string `json:"user_name,omitempty"`
	BroadcasterUserID    string `json:"broadcaster_user_id,omitempty"`
	BroadcasterUserLogin string `json:"broadcaster_user_login,omitempty"`
	BroadcasterUserName  string `json:"broadcaster_user_name,omitempty"`
	Tier                 string `json:"tier,omitempty"`
	Message              struct {
		Text   string `json:"text,omitempty"`
		Emotes []struct {
			Begin int    `json:"begin,omitempty"`
			End   int    `json:"end,omitempty"`
			ID    string `json:"id,omitempty"`
		} `json:"emotes,omitempty"`
	} `json:"message,omitempty"`
	CumulativeMonths int `json:"cumulative_months,omitempty"`
	StreakMonths     int `json:"streak_months,omitempty"`
	DurationMonths   int `json:"duration_months,omitempty"`
}
