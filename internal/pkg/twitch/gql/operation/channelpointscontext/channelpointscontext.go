package channelpointscontext

import (
	"encoding/json"
	"fmt"
	"time"
	"ttv-cli/internal/pkg/twitch/gql"
)

type persistedQuery struct {
	Version    int    `json:"version"`
	Sha256Hash string `json:"sha256Hash"`
}

type extensions struct {
	PersistedQuery persistedQuery `json:"persistedQuery"`
}

type variables struct {
	ChannelLogin string `json:"channelLogin"`
}

type request struct {
	OperationName string     `json:"operationName"`
	Variables     variables  `json:"variables"`
	Extensions    extensions `json:"extensions"`
}

func makeRequest(channelLogin string) request {
	return request{
		OperationName: "ChannelPointsContext",
		Variables: variables{
			ChannelLogin: channelLogin,
		},
		Extensions: extensions{
			PersistedQuery: persistedQuery{
				Version:    1,
				Sha256Hash: "9988086babc615a918a1e9a722ff41d98847acac822645209ac7379eecb27152",
			},
		},
	}
}

type ChannelPointsContext struct {
	Data struct {
		Community struct {
			Id          string `json:"id"`
			DisplayName string `json:"displayName"`
			Channel     struct {
				Id   string `json:"id"`
				Self struct {
					CommunityPoints struct {
						AvailableClaim          interface{}   `json:"availableClaim"`
						Typename                string        `json:"__typename"`
						Balance                 int           `json:"balance"`
						ActiveMultipliers       []interface{} `json:"activeMultipliers"`
						CanRedeemRewardsForFree bool          `json:"canRedeemRewardsForFree"`
						LastViewedContent       []struct {
							ContentType  string    `json:"contentType"`
							LastViewedAt time.Time `json:"lastViewedAt"`
							Typename     string    `json:"__typename"`
						} `json:"lastViewedContent"`
						UserRedemptions []interface{} `json:"userRedemptions"`
					} `json:"communityPoints"`
					Typename string `json:"__typename"`
				} `json:"self"`
				Typename                string `json:"__typename"`
				CommunityPointsSettings struct {
					Name  string `json:"name"`
					Image struct {
						Url      string `json:"url"`
						Url2X    string `json:"url2x"`
						Url4X    string `json:"url4x"`
						Typename string `json:"__typename"`
					} `json:"image"`
					Typename         string `json:"__typename"`
					AutomaticRewards []struct {
						Id                     string      `json:"id"`
						BackgroundColor        interface{} `json:"backgroundColor"`
						Cost                   interface{} `json:"cost"`
						DefaultBackgroundColor string      `json:"defaultBackgroundColor"`
						DefaultCost            int         `json:"defaultCost"`
						DefaultImage           struct {
							Url      string `json:"url"`
							Url2X    string `json:"url2x"`
							Url4X    string `json:"url4x"`
							Typename string `json:"__typename"`
						} `json:"defaultImage"`
						Image                         interface{} `json:"image"`
						IsEnabled                     bool        `json:"isEnabled"`
						IsHiddenForSubs               bool        `json:"isHiddenForSubs"`
						MinimumCost                   int         `json:"minimumCost"`
						Type                          string      `json:"type"`
						UpdatedForIndicatorAt         interface{} `json:"updatedForIndicatorAt"`
						GloballyUpdatedForIndicatorAt time.Time   `json:"globallyUpdatedForIndicatorAt"`
						Typename                      string      `json:"__typename"`
					} `json:"automaticRewards"`
					CustomRewards   []interface{} `json:"customRewards"`
					Goals           []interface{} `json:"goals"`
					IsEnabled       bool          `json:"isEnabled"`
					RaidPointAmount int           `json:"raidPointAmount"`
					EmoteVariants   []struct {
						Id           string `json:"id"`
						IsUnlockable bool   `json:"isUnlockable"`
						Emote        struct {
							Id       string `json:"id"`
							Token    string `json:"token"`
							Typename string `json:"__typename"`
						} `json:"emote"`
						Modifications []struct {
							Id    string `json:"id"`
							Emote struct {
								Id       string `json:"id"`
								Token    string `json:"token"`
								Typename string `json:"__typename"`
							} `json:"emote"`
							ModifierIconDark struct {
								Url      string `json:"url"`
								Url2X    string `json:"url2x"`
								Url4X    string `json:"url4x"`
								Typename string `json:"__typename"`
							} `json:"modifierIconDark"`
							ModifierIconLight struct {
								Url      string `json:"url"`
								Url2X    string `json:"url2x"`
								Url4X    string `json:"url4x"`
								Typename string `json:"__typename"`
							} `json:"modifierIconLight"`
							Title                         string    `json:"title"`
							GloballyUpdatedForIndicatorAt time.Time `json:"globallyUpdatedForIndicatorAt"`
							Typename                      string    `json:"__typename"`
						} `json:"modifications"`
						Typename string `json:"__typename"`
					} `json:"emoteVariants"`
					Earning struct {
						Id                     string `json:"id"`
						AveragePointsPerHour   int    `json:"averagePointsPerHour"`
						CheerPoints            int    `json:"cheerPoints"`
						ClaimPoints            int    `json:"claimPoints"`
						FollowPoints           int    `json:"followPoints"`
						PassiveWatchPoints     int    `json:"passiveWatchPoints"`
						RaidPoints             int    `json:"raidPoints"`
						SubscriptionGiftPoints int    `json:"subscriptionGiftPoints"`
						WatchStreakPoints      []struct {
							Points   int    `json:"points"`
							Typename string `json:"__typename"`
						} `json:"watchStreakPoints"`
						Multipliers []struct {
							ReasonCode string  `json:"reasonCode"`
							Factor     float64 `json:"factor"`
							Typename   string  `json:"__typename"`
						} `json:"multipliers"`
						Typename string `json:"__typename"`
					} `json:"earning"`
				} `json:"communityPointsSettings"`
			} `json:"channel"`
			Typename string `json:"__typename"`
			Self     struct {
				IsModerator bool   `json:"isModerator"`
				Typename    string `json:"__typename"`
			} `json:"self"`
		} `json:"community"`
		CurrentUser struct {
			Id              string `json:"id"`
			CommunityPoints struct {
				LastViewedContent []struct {
					ContentID    string    `json:"contentID"`
					ContentType  string    `json:"contentType"`
					LastViewedAt time.Time `json:"lastViewedAt"`
					Typename     string    `json:"__typename"`
				} `json:"lastViewedContent"`
				Typename string `json:"__typename"`
			} `json:"communityPoints"`
			Typename string `json:"__typename"`
		} `json:"currentUser"`
	} `json:"data"`
	Extensions struct {
		DurationMilliseconds int    `json:"durationMilliseconds"`
		OperationName        string `json:"operationName"`
		RequestID            string `json:"requestID"`
	} `json:"extensions"`
}

func Get(channelLogin string, authToken string) (ChannelPointsContext, error) {
	req := makeRequest(channelLogin)
	resp, err := gql.PostWithAuth(req, authToken)
	if err != nil {
		return ChannelPointsContext{}, fmt.Errorf("error with GQL request: %w", err)
	}

	var c ChannelPointsContext
	err = json.Unmarshal(resp, &c)

	return c, err
}
