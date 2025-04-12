package resubscribe

type conf struct {
	conditions `json:"conditions"`
	fields     `json:"fields"`
}

type conditions struct {
	Tier  int `json:"tier,omitempty"`
	Month int `json:"month,omitempty"`
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

type resubscribe struct {
	UserID               string `json:"user_id,omitempty"`
	UserLogin            string `json:"user_login,omitempty"`
	UserName             string `json:"user_name,omitempty"`
	BroadcasterUserID    string `json:"broadcaster_user_id,omitempty"`
	BroadcasterUserLogin string `json:"broadcaster_user_login,omitempty"`
	BroadcasterUserName  string `json:"broadcaster_user_name,omitempty"`
	Tier                 int    `json:"tier,omitempty"`
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
