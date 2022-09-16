package config

import (
	"fmt"
	"ttv-cli/internal/pkg/twitch/auth"
)

func (c *Config) GetTokenDetails() (TokenDetails, error) {
	if err := c.validateAuthToken(); err != nil {
		return TokenDetails{}, err
	}
	return *c.tokenDetails, nil
}

func (c *Config) validateAuthToken() error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if len(c.authToken) == 0 {
		fmt.Println("Auth token not found, generating a new one for you...")
		if err := c.refreshAuthToken(); err != nil {
			return fmt.Errorf("could not refresh auth token: %w", err)
		}
	}

	resp, err := auth.Validate(c.authToken)
	if err != nil {
		fmt.Println("Auth token is stale or invalid, generating a new one for you...")
		if err := c.refreshAuthToken(); err != nil {
			return fmt.Errorf("could not refresh auth token: %w", err)
		}
	}

	c.tokenDetails = &TokenDetails{
		ClientId:  resp.ClientId,
		Login:     resp.Login,
		Scopes:    resp.Scopes,
		UserId:    resp.UserId,
		ExpiresIn: resp.ExpiresIn,
	}

	return nil
}

func (c *Config) refreshAuthToken() error {
	authToken, err := auth.GetAccessToken("", "")
	if err != nil {
		return fmt.Errorf("validateAuthToken: Error getting Twitch access token: %w", err)
	}

	c.authToken = authToken
	if err := c.Save(); err != nil {
		return fmt.Errorf("validateAuthToken: Error saving config: %w", err)
	}

	return nil
}
