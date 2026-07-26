package discord

type (
	// https://docs.discord.com/developers/resources/user#user-object
	User struct {
		ID string `json:"id"`
		Username string `json:"username"`
		Discriminator string `json:"discriminator"`
		GlobalName string `json:"global_name"`
		Bot bool `json:"bot"`
		System bool `json:"system"`
	}

	GuildMember struct {
		User User `json:"user"`
		Nick string `json:"nick"`
	}

	// https://docs.discord.com/developers/interactions/receiving-and-responding#interaction-object
	Interaction[InteractionData any] struct {
		// ID of the interaction
		ID string `json:"id"`
		// ID of the application this interaction is for
		ApplicationID string `json:"application_id"`
		// Type of interaction
		Type InteractionType `json:"type"`
		// Interaction data payload
		Data *InteractionData `json:"data"`
		// Guild that the interaction was sent from
		GuildID string `json:"guild_id"`
		// Channel that the interaction was sent from
		ChannelID string `json:"channel_id"`
		// Guild member data for the invoking user, including permissions
		Member GuildMember `json:"member"`
		// Selected language of the invoking user
		Locale *string `json:"locale"`
		// Guild's preferred locale, if invoked in a guild
		GuildLocale *string `json:"guild_locale"`
		Context *InteractionContextType `json:"context"`
	}

	ApplicationCommandData struct {
		// ID of the invoked command
		ID string `json:"id"`
		Name string `json:"name"`
		// type of the invoked command
		Type int64 `json:"type"`
	}

	InteractionType int8
	InteractionContextType int8
)

const (
	InteractionTypePing InteractionType = 1
	InteractionTypeApplicationCommand InteractionType = 2
	InteractionTypeMessageComponent InteractionType = 3
	InteractionTypeApplicationCommandAutocomplete InteractionType = 4
	InteractionTypeModalSubmit InteractionType = 5

	// Interaction can be used within servers
	InteractionContextTypeGuild InteractionContextType = 0
	// Interaction can be used within DMs with the app’s bot user
	InteractionContextTypeBotDM InteractionContextType = 1
	// Interaction can be used within Group DMs and DMs other than the app’s bot user
	InteractionContextTypePrivateChannel InteractionContextType = 2
)
