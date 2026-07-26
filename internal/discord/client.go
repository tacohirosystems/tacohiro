package discord

import (
	"fmt"
	"net/http"
)

type (
	Client struct {
		Config *DiscordConfig
		HTTPClient *http.Client
	}

	BotToken string
	ApplicationID string
)

var (
	ErrDiscordClientConfigNil = fmt.Errorf("discord: client's config cannot be nil")
	ErrDiscordClientHTTPClientNil = fmt.Errorf("discord: client's http client cannot be nil")
)

func NewClient(httpClient *http.Client, config *DiscordConfig) (*Client, error) {
	if config == nil {
		return nil, ErrDiscordClientConfigNil
	}

	if httpClient == nil {
		return nil, ErrDiscordClientHTTPClientNil
	}

	client := &Client{
		Config:     config,
		HTTPClient: httpClient,
	}
	return client, nil
}

// TODO: IMPLEMENT
func (cc Client) HealthCheck() error {
	return nil
}
