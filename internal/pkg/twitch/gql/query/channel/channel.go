package channel

import (
	"encoding/json"
	"fmt"
	"ttv-cli/internal/pkg/config"
	"ttv-cli/internal/pkg/twitch/gql"
)

const getChannelQuery = `query Channel($name: String) {
	channel(name: $name) {
		id
		name
		displayName
		communityPointsSettings {
			customRewards {
				id
				backgroundColor
				title
				prompt
				cost
				defaultImage {
					url
				}
				image {
					url
				}
				cooldownExpiresAt
				globalCooldownSetting {
					globalCooldownSeconds
					isEnabled
				}
				isPaused
				isEnabled
				isUserInputRequired
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

type Image struct {
	Url string `json:"url"`
}

type Channel struct {
	Id                      string                         `json:"id"`
	Name                    string                         `json:"name"`
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
	DefaultImage          Image                                            `json:"defaultImage"`
	Image                 Image                                            `json:"image"`
	Id                    string                                           `json:"id"`
	CooldownExpiresAt     string                                           `json:"cooldownExpiresAt"`
	GlobalCooldownSetting CommunityPointsCustomRewardGlobalCooldownSetting `json:"globalCooldownSetting"`
	IsPaused              bool                                             `json:"isPaused"`
	IsEnabled             bool                                             `json:"isEnabled"`
	IsUserInputRequired   bool                                             `json:"isUserInputRequired"`
}

type CommunityPointsCustomRewardGlobalCooldownSetting struct {
	GlobalCooldownSeconds int  `json:"globalCooldownSeconds"`
	IsEnabled             bool `json:"isEnabled"`
}

// GetChannel Note that this function will not throw if a channel was not found for the provided name
func GetChannel(config config.Config, name string) (Channel, error) {
	request := GetChannelRequest{Query: getChannelQuery}
	request.Variables.Name = name

	body, err := gql.Post(config, request)
	if err != nil {
		return Channel{}, fmt.Errorf("GetChannel: error with GQL request: %w", err)
	}

	var response GetChannelResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return Channel{}, fmt.Errorf("GetChannel: error unmarshalling GQL response: %w", err)
	}

	return response.Data.Channel, nil
}
