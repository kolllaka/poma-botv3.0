package model

type Reward struct {
	Type string `json:"type"`
	Data Data   `json:"data"`
}

type Data struct {
	Timestamp  string     `json:"timestamp"`
	Redemption Redemption `json:"redemption"`
}

type Redemption struct {
	ID         string      `json:"id"`
	User       User        `json:"user"`
	ChannelID  string      `json:"channel_id"`
	RedeemedAt string      `json:"redeemed_at"`
	Reward     RewardClass `json:"reward"`
	UserInput  string      `json:"user_input"`
	Status     string      `json:"status"`
	Cursor     string      `json:"cursor"`
}

type RewardClass struct {
	ID                                string              `json:"id"`
	ChannelID                         string              `json:"channel_id"`
	Title                             string              `json:"title"`
	Prompt                            string              `json:"prompt"`
	Cost                              int64               `json:"cost"`
	IsUserInputRequired               bool                `json:"is_user_input_required"`
	IsSubOnly                         bool                `json:"is_sub_only"`
	Image                             any                 `json:"image"`
	DefaultImage                      DefaultImage        `json:"default_image"`
	BackgroundColor                   string              `json:"background_color"`
	IsEnabled                         bool                `json:"is_enabled"`
	IsPaused                          bool                `json:"is_paused"`
	IsInStock                         bool                `json:"is_in_stock"`
	MaxPerStream                      MaxPerStream        `json:"max_per_stream"`
	ShouldRedemptionsSkipRequestQueue bool                `json:"should_redemptions_skip_request_queue"`
	TemplateID                        any                 `json:"template_id"`
	UpdatedForIndicatorAt             string              `json:"updated_for_indicator_at"`
	MaxPerUserPerStream               MaxPerUserPerStream `json:"max_per_user_per_stream"`
	GlobalCooldown                    GlobalCooldown      `json:"global_cooldown"`
	RedemptionsRedeemedCurrentStream  any                 `json:"redemptions_redeemed_current_stream"`
	CooldownExpiresAt                 any                 `json:"cooldown_expires_at"`
}

type DefaultImage struct {
	URL1X string `json:"url_1x"`
	URL2X string `json:"url_2x"`
	URL4X string `json:"url_4x"`
}

type GlobalCooldown struct {
	IsEnabled             bool  `json:"is_enabled"`
	GlobalCooldownSeconds int64 `json:"global_cooldown_seconds"`
}

type MaxPerStream struct {
	IsEnabled    bool  `json:"is_enabled"`
	MaxPerStream int64 `json:"max_per_stream"`
}

type MaxPerUserPerStream struct {
	IsEnabled           bool  `json:"is_enabled"`
	MaxPerUserPerStream int64 `json:"max_per_user_per_stream"`
}

type User struct {
	ID          string `json:"id"`
	Login       string `json:"login"`
	DisplayName string `json:"display_name"`
}
