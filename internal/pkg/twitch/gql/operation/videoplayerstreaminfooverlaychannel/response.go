package videoplayerstreaminfooverlaychannel

type Stream struct {
	Id           string `json:"id"`
	ViewersCount int    `json:"viewersCount"`
	Tags         []struct {
		Id            string `json:"id"`
		LocalizedName string `json:"localizedName"`
		Typename      string `json:"__typename"`
	} `json:"tags"`
	Typename string `json:"__typename"`
}

type Response struct {
	Data struct {
		User struct {
			Id                string `json:"id"`
			ProfileURL        string `json:"profileURL"`
			DisplayName       string `json:"displayName"`
			Login             string `json:"login"`
			ProfileImageURL   string `json:"profileImageURL"`
			BroadcastSettings struct {
				Id    string `json:"id"`
				Title string `json:"title"`
				Game  struct {
					Id          string `json:"id"`
					DisplayName string `json:"displayName"`
					Name        string `json:"name"`
					Typename    string `json:"__typename"`
				} `json:"game"`
				Typename string `json:"__typename"`
			} `json:"broadcastSettings"`
			Stream   *Stream `json:"stream"`
			Typename string  `json:"__typename"`
		} `json:"user"`
	} `json:"data"`
	Extensions struct {
		DurationMilliseconds int    `json:"durationMilliseconds"`
		OperationName        string `json:"operationName"`
		RequestID            string `json:"requestID"`
	} `json:"extensions"`
}
