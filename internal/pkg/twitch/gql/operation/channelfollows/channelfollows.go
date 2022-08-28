package channelfollows

import (
	"encoding/json"
	"log"
	"time"
	"ttv-cli/internal/pkg/twitch/gql"
)

type variables struct {
	Limit int    `json:"limit"`
	Order string `json:"order"`
}

type extensions struct {
	PersistedQuery persistedQuery `json:"persistedQuery"`
}

type persistedQuery struct {
	Version    int    `json:"version"`
	Sha256Hash string `json:"sha256Hash"`
}

type channelFollowsRequest struct {
	OperationName string     `json:"operationName"`
	Variables     variables  `json:"variables"`
	Extensions    extensions `json:"extensions"`
}

func makeRequest() channelFollowsRequest {
	return channelFollowsRequest{
		OperationName: "ChannelFollows",
		Variables: variables{
			Limit: 100,
			Order: "ASC",
		},
		Extensions: extensions{
			PersistedQuery: persistedQuery{
				Version:    1,
				Sha256Hash: "4b9cb31b54b9213e5760f2f6e9e935ad09924cac2f78aac51f8a64d85f028ed0",
			},
		},
	}
}

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

type ChannelFollow struct {
	DisplayName string
	Login       string
}

// GetChannelFollows TODO: Implement cursor following to handle >100 follows
func GetChannelFollows(authToken string) []ChannelFollow {
	req := makeRequest()
	respBody, err := gql.PostWithAuth(req, authToken)
	if err != nil {
		log.Fatalln(err)
	}

	var resp response
	if err := json.Unmarshal(respBody, &resp); err != nil {
		log.Fatalln(err)
	}

	follows := make([]ChannelFollow, 0)
	for _, node := range resp.Data.User.Follows.Edges {
		if node.Typename == "FollowEdge" {
			follows = append(follows, ChannelFollow{
				DisplayName: node.Node.DisplayName,
				Login:       node.Node.Login,
			})
		}
	}

	return follows
}
