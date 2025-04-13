package eventsub

import "time"

type EventRaid struct {
	FromBroadcasterUserID    string `json:"from_broadcaster_user_id,omitempty"`
	FromBroadcasterUserLogin string `json:"from_broadcaster_user_login,omitempty"`
	FromBroadcasterUserName  string `json:"from_broadcaster_user_name,omitempty"`
	ToBroadcasterUserID      string `json:"to_broadcaster_user_id,omitempty"`
	ToBroadcasterUserLogin   string `json:"to_broadcaster_user_login,omitempty"`
	ToBroadcasterUserName    string `json:"to_broadcaster_user_name,omitempty"`
	Viewers                  int    `json:"viewers,omitempty"`
}

type EventFollow struct {
	UserID               string    `json:"user_id,omitempty"`
	UserLogin            string    `json:"user_login,omitempty"`
	UserName             string    `json:"user_name,omitempty"`
	BroadcasterUserID    string    `json:"broadcaster_user_id,omitempty"`
	BroadcasterUserLogin string    `json:"broadcaster_user_login,omitempty"`
	BroadcasterUserName  string    `json:"broadcaster_user_name,omitempty"`
	FollowedAt           time.Time `json:"followed_at,omitempty"`
}

type EventSubscribe struct {
	UserID               string `json:"user_id,omitempty"`
	UserLogin            string `json:"user_login,omitempty"`
	UserName             string `json:"user_name,omitempty"`
	BroadcasterUserID    string `json:"broadcaster_user_id,omitempty"`
	BroadcasterUserLogin string `json:"broadcaster_user_login,omitempty"`
	BroadcasterUserName  string `json:"broadcaster_user_name,omitempty"`
	Tier                 int    `json:"tier,omitempty"`
	IsGift               bool   `json:"is_gift,omitempty"`
}

type EventSubscribeGift struct {
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

type EventReSubscribe struct {
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

type EventCheer struct {
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
