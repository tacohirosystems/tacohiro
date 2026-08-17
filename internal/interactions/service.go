package interactions

import (
	"fmt"
	"log/slog"
	"math/rand/v2"

	"github.com/tacohirosystems/tacohiro/internal/discord"
	"github.com/tacohirosystems/tacohiro/internal/events"
)

type (
	Service struct {
		EventsService    *events.Service
		DiscordBotClient *discord.Client
		Logger           *slog.Logger
	}
)

const (
	CommandOptionNameTacoQuantity    discord.CommandOptionName = "quantity"
	CommandOptionNameRecipient       discord.CommandOptionName = "recipient"
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
		Options: []discord.CommandOption{
			{
				Type:         discord.CommandOptionTypeMentionable,
				Name:         CommandOptionNameRecipient,
				Description:  "Who would you like to give tacos to?",
				Required:     true,
				Autocomplete: true,
			},
			{
				Type:         discord.CommandOptionTypeInteger,
				Name:         CommandOptionNameTacoQuantity,
				Description:  "How many?",
				Required:     true,
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

func (s *Service) ProcessInteraction(interaction discord.Interaction[discord.ApplicationCommandData]) error {
	var event events.Event

	event.SenderID = interaction.Member.User.ID
	// TODO: Check if recipient is a role, then we send to all members in that role.
	// TODO: Check if recipient is a bot. Maybe consider allowing bots to send and receive.
	switch interaction.Data.Name {
	case SlashCommandGive:
		for _, o := range interaction.Data.Options {
			switch o.Name {
			case CommandOptionNameRecipient:
				if o.Value.ValueString == nil {
					// TODO: Return error
					return nil
				}
				s.Logger.Debug(fmt.Sprintf("Recipient ID: %s\n", *o.Value.ValueString))
				event.RecipientIDs = append(event.RecipientIDs, discord.UserID(*o.Value.ValueString))
			case CommandOptionNameRecipientExtra1:
				if o.Value.ValueString == nil {
					// TODO: Return error
					return nil
				}
				s.Logger.Debug(fmt.Sprintf("Recipient ID: %s\n", *o.Value.ValueString))
				event.RecipientIDs = append(event.RecipientIDs, discord.UserID(*o.Value.ValueString))
			case CommandOptionNameRecipientExtra2:
				if o.Value.ValueString == nil {
					// TODO: Return error
					return nil
				}
				s.Logger.Debug(fmt.Sprintf("Recipient ID: %s\n", *o.Value.ValueString))
				event.RecipientIDs = append(event.RecipientIDs, discord.UserID(*o.Value.ValueString))
			case CommandOptionNameTacoQuantity:
				if o.Value.ValueInt64 != nil {
					s.Logger.Debug(fmt.Sprintf("Quantity of tacos: %d\n", *o.Value.ValueInt64))
					event.Quantity = *o.Value.ValueInt64
				} else {
					s.Logger.Debug(fmt.Sprintf("Quantity of tacos: %+v", o.Value))
					// TODO: Return error
					return nil
				}
			default:
				// TODO: Return error
				s.Logger.Error(fmt.Sprintf("Unknown option %s\n", o.Name))
			}
		}
	default:
		// TODO: Return error
		return nil
	}

	s.Logger.Debug(fmt.Sprintf("%+v", event))

	if err := s.EventsService.CreateEvent(event); err != nil {
		return err
	}

	// FIXME: Message builder refactor
	var recipientsStr string
	for i := 0; i < len(event.RecipientIDs); i++ {
		if i == 0 {
			recipientsStr = fmt.Sprintf("<@%s>", event.RecipientIDs[i])
			continue
		}

		recipientsStr = fmt.Sprintf("%s <@%s>", recipientsStr, event.RecipientIDs[i])
	}

	var msg string
	if event.Quantity == 0 {
		msg = fmt.Sprintf("<@%s> gave %s %dx :taco:. _Hm... someone is stingy_", event.SenderID, recipientsStr, event.Quantity)
	}

	if event.Quantity < 0 {
		msg = fmt.Sprintf("<@%s> gave %s a :taco: debt of %d.", event.SenderID, recipientsStr, event.Quantity*-1)
	}

	var loseOneTaco bool
	if rand.Float64() < 0.05 {
		loseOneTaco = true
	}

	if loseOneTaco && event.Quantity >= 1 {
		event.Quantity -= 1
		msg = fmt.Sprintf("<@%s> tried to give %s %dx :taco:! But Hiro ate one so it's one less...", event.SenderID, recipientsStr, event.Quantity+1)
	} else if event.Quantity >= 1 {
		msg = fmt.Sprintf("<@%s> gave %s %dx :taco:!", event.SenderID, recipientsStr, event.Quantity)
	}

	if err := s.DiscordBotClient.CreateInteractionResponse(discord.InteractionCreateCallbackResponse{
		ID:    interaction.ID,
		Token: interaction.Token,
		Body: discord.InteractionResponseObject{
			Type: discord.InteractionCallbackTypeCHANNEL_MESSAGE_WITH_SOURCE,
			Data: &discord.InteractionCallbackData{
				TTS:     new(false),
				Content: new(msg),
				Flags:   nil,
			},
		},
	}); err != nil {
		s.Logger.Error(fmt.Sprintf("%s", err.Error()))
	}

	return nil
}
