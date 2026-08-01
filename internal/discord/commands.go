package discord

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
)

type (
	CommandOption struct {
		Type CommandOptionType `json:"type"`
		Name CommandOptionName `json:"name"`
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
	CommandOptionName string
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
	CommandTypeMessage = 3
	CommandTypePrimaryEntryPoint = 4
)

var (
	ErrDiscordCommandFailedToCreate = fmt.Errorf("discord: failed to create command")
	ErrDiscordCommandFailedToSerializeCreatePayload = fmt.Errorf("discord: failed to serialize create command payload")
)

func (c *Client) CreateCommand(body CreateCommandParams) error {
	bodyRaw, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("%s %w", err.Error(), ErrDiscordCommandFailedToSerializeCreatePayload)
	}

	url := fmt.Sprintf("%s/applications/%s/commands", DISCORD_BASE_URL, c.Config.ApplicationID)
	buf := bytes.NewBuffer(bodyRaw)
	res, err := c.Do(url, buf)
	if err != nil {
		return fmt.Errorf("%s %w", err.Error(), ErrDiscordCommandFailedToCreate)
	}

	if !(res.StatusCode == 200 || res.StatusCode == 201) {
		rawRes, err := io.ReadAll(res.Body)
		if err != nil {
			return fmt.Errorf("%s %w", err.Error(), ErrDiscordCommandFailedToCreate)
		}
		fmt.Printf("%s", string(rawRes))
	}

	return nil
}
