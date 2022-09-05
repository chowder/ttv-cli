package channelpointscontext

import "time"

type Response struct {
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
