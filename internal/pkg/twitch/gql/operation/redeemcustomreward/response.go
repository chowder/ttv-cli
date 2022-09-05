package redeemcustomreward

type Error struct {
	Code     string `json:"code"`
	Typename string `json:"__typename"`
}

type Response struct {
	Data struct {
		RedeemCommunityPointsCustomReward struct {
			Error    *Error `json:"error"`
			Typename string `json:"__typename"`
		} `json:"redeemCommunityPointsCustomReward"`
	} `json:"data"`
	Extensions struct {
		DurationMilliseconds int    `json:"durationMilliseconds"`
		OperationName        string `json:"operationName"`
		RequestID            string `json:"requestID"`
	} `json:"extensions"`
}

func ErrorToDisplay(code string) string {
	switch code {
	case "GLOBAL_COOLDOWN":
		return "Reward on cooldown"
	case "NOT_ENOUGH_POINTS":
		return "Not enough points"
	default:
		return code
	}
}
