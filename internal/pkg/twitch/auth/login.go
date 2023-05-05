package auth

import (
	"errors"
	"fmt"
	"github.com/Jeffail/gabs/v2"
	"github.com/go-resty/resty/v2"
	"math/rand"
	"time"
	"ttv-cli/internal/pkg/twitch"
)

const twitchAuthUrl = "https://id.twitch.tv/oauth2/device"
const twitchTokenUrl = "https://id.twitch.tv/oauth2/token"

type Verification struct {
	UserCode        string
	DeviceCode      string
	Interval        float64
	ExpiresIn       float64
	VerificationUrl string
}

var NotYetAuthenticatedError = errors.New("not yet authentication")

func GetAccessToken() (string, error) {
	verification, err := getVerificationCode()
	fmt.Println("Enter code:", verification.UserCode, "at", verification.VerificationUrl)

	if err != nil {
		return "", err
	}

	for {
		time.Sleep(5 * time.Second)
		accessToken, err := checkVerified(verification)
		if errors.Is(err, NotYetAuthenticatedError) {
			continue
		}

		if err != nil {
			return "", err
		}

		return accessToken, nil
	}
}

func checkVerified(verification Verification) (string, error) {
	client := resty.New()
	client = prepareHeaders(client)

	resp, err := client.R().
		SetFormData(map[string]string{
			"client_id":   twitch.DefaultClientId,
			"device_code": verification.DeviceCode,
			"grant_type":  "urn:ietf:params:oauth:grant-type:device_code",
		}).
		Post(twitchTokenUrl)

	if err != nil {
		return "", fmt.Errorf("error checking verification: %w", err)
	}

	if resp.StatusCode() != 200 {
		return "", NotYetAuthenticatedError
	}

	parsed, err := gabs.ParseJSON(resp.Body())
	return parsed.Path("access_token").Data().(string), err
}

func getVerificationCode() (Verification, error) {
	var verification Verification
	client := resty.New()
	client = prepareHeaders(client)

	resp, err := client.R().
		SetFormData(map[string]string{
			"client_id": twitch.DefaultClientId,
			"scopes":    "channel_read chat:read user_blocks_edit user_blocks_read user_follows_edit user_read",
		}).
		Post(twitchAuthUrl)

	if err != nil {
		return verification, fmt.Errorf("error getting verification codes: %w", err)
	}

	parsed, err := gabs.ParseJSON(resp.Body())
	if err != nil {
		return verification, fmt.Errorf("error parsing verification response: %w", err)
	}

	verification.DeviceCode = parsed.Path("device_code").Data().(string)
	verification.UserCode = parsed.Path("user_code").Data().(string)
	verification.ExpiresIn = parsed.Path("expires_in").Data().(float64)
	verification.Interval = parsed.Path("interval").Data().(float64)
	verification.VerificationUrl = parsed.Path("verification_uri").Data().(string)

	return verification, nil
}

func prepareHeaders(client *resty.Client) *resty.Client {
	return client.
		SetHeader("authority", "id.twitch.tv").
		SetHeader("accept", "application/json").
		SetHeader("accept-language", "en-US,en;q=0.9").
		SetHeader("content-type", "application/x-www-form-urlencoded").
		SetHeader("origin", "https://android.tv.twitch.tv").
		SetHeader("referer", "https://android.tv.twitch.tv/").
		SetHeader("user-agent", twitch.DefaultUserAgent)
}

func CreateRandomDeviceId() string {
	const chars = "abcdefghijklmnopqrstuvwxyz01234567890"
	const length = 26

	deviceId := make([]byte, length)
	for i := range deviceId {
		deviceId[i] = chars[rand.Intn(len(chars))]
	}

	return string(deviceId)
}
