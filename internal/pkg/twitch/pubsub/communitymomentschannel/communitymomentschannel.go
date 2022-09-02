package communitymomentschannel

type Response struct {
	Type string `json:"type"`
	Data struct {
		MomentId string `json:"moment_id"`
	} `json:"data"`
}
