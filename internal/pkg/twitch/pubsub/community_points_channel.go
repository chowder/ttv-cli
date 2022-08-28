package pubsub

type Image struct {
	Url1x string `json:"url_1x"`
	Url2x string `json:"url_2x"`
	Url4x string `json:"url_4x"`
}

type UpdatedReward struct {
	Id                  string `json:"id"`
	ChannelId           string `json:"channel_id"`
	Title               string `json:"title"`
	Prompt              string `json:"prompt"`
	Cost                int    `json:"cost"`
	IsUserInputRequired bool   `json:"is_user_input_required"`
	IsSubOnly           bool   `json:"is_sub_only"`
	Image               Image  `json:"image"`
	DefaultImage        Image  `json:"default_image"`
	BackgroundColor     string `json:"background_color"`
	IsEnabled           bool   `json:"is_enabled"`
	IsPaused            bool   `json:"is_paused"`
	IsInStock           bool   `json:"is_in_stock"`
	CooldownExpiresAt   string `json:"cooldown_expires_at"`
}

type Data struct {
	Timestamp     string        `json:"timestamp"`
	UpdatedReward UpdatedReward `json:"updated_reward"`
}

type CommunityPointsChannelResponse struct {
	Type string `json:"type"`
	Data Data   `json:"data"`
}
