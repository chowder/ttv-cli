package gql

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

const getChannelQuery = `query Channel($name: String) {
	channel(name: $name) {
		id
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

	requestBody, err := json.Marshal(request)
	if err != nil {
		log.Fatalln(err)
	}

	// Make a POST request
	req, err := http.NewRequest("POST", gqlApiUrl, bytes.NewBuffer(requestBody))
	if err != nil {
		log.Fatalln(err)
	}
	req.Header.Add("Client-ID", defaultClientId)

	// Execute the POST request
	client := &http.Client{}
	httpResp, err := client.Do(req)

	if err != nil {
		log.Fatalln(err)
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Panic(err)
		}
	}(httpResp.Body)

	// Read the response body
	body, err := ioutil.ReadAll(httpResp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	// Unmarshal response
	var response GetChannelResponse
	if err := json.Unmarshal(body, &response); err != nil {
		log.Fatalln(err)
	}

	return response.Data.Channel
}
