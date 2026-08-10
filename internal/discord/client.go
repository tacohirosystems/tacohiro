package discord

import (
	"fmt"
	"io"
	"net/http"
)

type (
	Client struct {
		Config     *DiscordConfig
		HTTPClient *http.Client
	}

	BotToken      string
	ApplicationID string
)

var (
	ErrDiscordClientConfigNil     = fmt.Errorf("discord: client's config cannot be nil")
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

func (c *Client) Do(url string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		panic(fmt.Sprintf("failed to register command: %s", err.Error()))
	}

	req.Header.Add("authorization", fmt.Sprintf("Bot %s", c.Config.Token))
	req.Header.Add("content-type", "application/json")
	return c.HTTPClient.Do(req)
}
