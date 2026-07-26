package discord

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type (
	CommandsClient = *Client

	CommandOption struct {
		Type CommandOptionType `json:"type"`
		Name string `json:"name"`
		Description string `json:"description"`
		Required bool `json:"required"`
		Autocomplete bool `json:"autocomplete"`
	}

	CreateCommandParams struct {
		Name string `json:"name"`
		Type int64 `json:"type"`
		Description string `json:"description"`
		Options []CommandOption `json:"options"`
	}

	CommandOptionType int8
)

const (
	DISCORD_BASE_URL string = "https://discord.com/api/v10"

	// Command option types
	// https://docs.discord.com/developers/interactions/application-commands#application-command-object-application-command-option-type

	CommandOptionTypeSubCommand CommandOptionType = 1
	CommandOptionTypeSubCommandGroup CommandOptionType = 2
	CommandOptionTypeString CommandOptionType = 3
	// Any integer between -2^53+1 and 2^53-1
	CommandOptionTypeInteger CommandOptionType = 4
	CommandOptionTypeBoolean CommandOptionType = 5
	CommandOptionTypeUser CommandOptionType = 6
	// Includes all channel types + categories
	CommandOptionTypeChannel CommandOptionType = 7
	CommandOptionTypeRole CommandOptionType = 8
	// Includes users and roles
	CommandOptionTypeMentionable CommandOptionType = 9
	// Any double between -2^53 and 2^53
	CommandOptionTypeNumber CommandOptionType = 10
	CommandOptionTypeAttachment CommandOptionType = 11


	// Command Types
	// https://docs.discord.com/developers/interactions/application-commands

	CommandTypeChatInput = 1
	CommandTypeUser = 2
	CommandTypeMessage = 2
	CommandTypePrimaryEntryPoint = 2
)

var (
	ErrFailedToCreateCommand = fmt.Errorf("discord: failed to create command")
)

func (cc CommandsClient) CreateCommand(params CreateCommandParams) error {
	url := fmt.Sprintf("%s/applications/%s/commands", DISCORD_BASE_URL, cc.Config.ApplicationID)

	bodyRaw, err := json.Marshal(params)
	if err != nil {
		panic("bye")
	}

	buf := bytes.NewBuffer(bodyRaw)
	req, err := http.NewRequest(http.MethodPost, url, buf)
	if err != nil {
		panic(fmt.Sprintf("failed to register command: %s", err.Error()))
	}

	req.Header.Add("authorization", fmt.Sprintf("Bot %s", cc.Config.Token))
	req.Header.Add("content-type", "application/json")

	res, err := cc.HTTPClient.Do(req)
	if err != nil {
		panic(fmt.Sprintf("failed to send request: %s", err.Error()))
	}

	// log.Printf("\nSlash commands: Register command status code: %d\n", res.StatusCode)

	if !(res.StatusCode == 200 || res.StatusCode == 201) {
		rawRes, err := io.ReadAll(res.Body)
		if err != nil {
			panic("failed to read from body")
		}
		fmt.Printf("%s", string(rawRes))
		return fmt.Errorf("%s %w", err.Error(), ErrFailedToCreateCommand)
	}

	return nil
}
