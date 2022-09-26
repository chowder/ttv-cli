package communitymomentschannel

const Topic = "community-moments-channel-v1"

type Response struct {
	Type string `json:"type"`
	Data struct {
		MomentId  string `json:"moment_id"`
		ChannelId string `json:"channel_id"`
		ClipSlug  string `json:"clip_slug"`
	} `json:"data"`
}
