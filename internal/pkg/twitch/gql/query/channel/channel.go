package channel

import (
	"encoding/json"
	"log"
	"ttv-cli/internal/pkg/twitch/gql"
)

const getChannelQuery = `query Channel($name: String) {
	channel(name: $name) {
		id
		displayName
		communityPointsSettings {
			customRewards {
				id
				backgroundColor
				title
				prompt
				cost
				cooldownExpiresAt
				globalCooldownSetting {
					globalCooldownSeconds
					isEnabled
				}
				isPaused
				isEnabled
			}
		}
	}
}`

type GetChannelRequest struct {
	Query     string `json:"query"`
	Variables struct {
		Name string `json:"name"`
	} `json:"variables"`
}

type GetChannelResponse struct {
	Data struct {
		Channel Channel `json:"channel"`
	} `json:"data"`
}

type Channel struct {
	Id                      string                         `json:"id"`
	DisplayName             string                         `json:"displayName"`
	CommunityPointsSettings CommunityPointsChannelSettings `json:"communityPointsSettings"`
}

type CommunityPointsChannelSettings struct {
	CustomRewards []CommunityPointsCustomReward `json:"customRewards"`
}

type CommunityPointsCustomReward struct {
	BackgroundColor       string                                           `json:"backgroundColor"`
	Title                 string                                           `json:"title"`
	Prompt                string                                           `json:"prompt"`
	Cost                  int                                              `json:"cost"`
	Id                    string                                           `json:"id"`
	CooldownExpiresAt     string                                           `json:"cooldownExpiresAt"`
	GlobalCooldownSetting CommunityPointsCustomRewardGlobalCooldownSetting `json:"globalCooldownSetting"`
	IsPaused              bool                                             `json:"isPaused"`
	IsEnabled             bool                                             `json:"isEnabled"`
}

type CommunityPointsCustomRewardGlobalCooldownSetting struct {
	GlobalCooldownSeconds int  `json:"globalCooldownSeconds"`
	IsEnabled             bool `json:"isEnabled"`
}

func GetChannel(name string) Channel {
	request := GetChannelRequest{Query: getChannelQuery}
	request.Variables.Name = name

	body, err := gql.Post(request)
	if err != nil {
		log.Fatalln(err)
	}

	var response GetChannelResponse
	if err := json.Unmarshal(body, &response); err != nil {
		log.Fatalln(err)
	}

	return response.Data.Channel
}
