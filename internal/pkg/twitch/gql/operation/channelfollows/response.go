package channelfollows

import "time"

type response struct {
	Data struct {
		User struct {
			Id      string `json:"id"`
			Follows struct {
				Edges []struct {
					Cursor string `json:"cursor"`
					Node   struct {
						Id              string  `json:"id"`
						BannerImageURL  *string `json:"bannerImageURL"`
						DisplayName     string  `json:"displayName"`
						Login           string  `json:"login"`
						ProfileImageURL string  `json:"profileImageURL"`
						Self            struct {
							CanFollow bool `json:"canFollow"`
							Follower  struct {
								DisableNotifications bool      `json:"disableNotifications"`
								FollowedAt           time.Time `json:"followedAt"`
								Typename             string    `json:"__typename"`
							} `json:"follower"`
							Typename   string      `json:"__typename"`
							Friendship interface{} `json:"friendship"`
						} `json:"self"`
						Typename     string      `json:"__typename"`
						Activity     interface{} `json:"activity"`
						Availability interface{} `json:"availability"`
					} `json:"node"`
					Typename string `json:"__typename"`
				} `json:"edges"`
				PageInfo struct {
					HasNextPage bool   `json:"hasNextPage"`
					Typename    string `json:"__typename"`
				} `json:"pageInfo"`
				Typename string `json:"__typename"`
			} `json:"follows"`
			Typename string `json:"__typename"`
		} `json:"user"`
	} `json:"data"`
	Extensions struct {
		DurationMilliseconds int    `json:"durationMilliseconds"`
		OperationName        string `json:"operationName"`
		RequestID            string `json:"requestID"`
	} `json:"extensions"`
}
