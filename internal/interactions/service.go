package interactions

import (
	"github.com/tacohirosystems/tacohiro/internal/discord"
)

const (
	CommandOptionNameTacoQuantity discord.CommandOptionName = "quantity"
	CommandOptionNameRecipient discord.CommandOptionName = "recipient"
	CommandOptionNameRecipientExtra1 discord.CommandOptionName = "recipient_2"
	CommandOptionNameRecipientExtra2 discord.CommandOptionName = "recipient_3"
)

// Registers the default TacoHiro commands
func InitCommands(client *discord.Client) error {
	// TODO: Register commands
	// - Leaderboards?
	return client.CreateCommand(discord.CreateCommandParams{
		Name:        "give",
		Type:        discord.CommandTypeChatInput,
		Description: "Give tacos to others!",
		Options:     []discord.CommandOption{
			{
				Type:         discord.CommandOptionTypeMentionable,
				Name:         CommandOptionNameRecipient,
				Description:  "Who would you like to give tacos to?",
				Required:     true,
				Autocomplete: true,
			},
			{
				Type:        discord.CommandOptionTypeInteger,
				Name:        CommandOptionNameTacoQuantity,
				Description: "How many?",
				Required:    true,
				Autocomplete: true,
			},
			{
				Type:         discord.CommandOptionTypeMentionable,
				Name:         CommandOptionNameRecipientExtra1,
				Description:  "Who else would you like to give tacos to?",
				Required:     false,
				Autocomplete: true,
			},
			{
				Type:         discord.CommandOptionTypeMentionable,
				Name:         CommandOptionNameRecipientExtra2,
				Description:  "Who else would you like to give tacos to?",
				Required:     false,
				Autocomplete: true,
			},
		},
	})
}
