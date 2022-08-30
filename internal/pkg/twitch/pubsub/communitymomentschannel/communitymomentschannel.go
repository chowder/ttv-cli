package communitymomentschannel

type Data struct {
	MomentId string `json:"moment_id"`
}

type Response struct {
	Type string `json:"type"`
	Data Data   `json:"data"`
}
