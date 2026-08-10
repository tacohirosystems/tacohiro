package discord

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"slices"
	"strconv"
)

type (
	GuildMember struct {
		User User   `json:"user"`
		Nick string `json:"nick"`
	}

	// https://docs.discord.com/developers/interactions/receiving-and-responding#interaction-object
	Interaction[InteractionData any] struct {
		// ID of the interaction
		ID InteractionID `json:"id"`
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
		GuildLocale *string                 `json:"guild_locale"`
		Context     *InteractionContextType `json:"context"`
		Token       InteractionToken        `json:"token"`
	}

	ApplicationCommandData struct {
		// ID of the invoked command
		ID   string                 `json:"id"`
		Name ApplicationCommandName `json:"name"`
		// type of the invoked command
		Type int64 `json:"type"`
		// Params + values from the user
		Options []ApplicationCommandDataOption `json:"options"`
	}

	ApplicationCommandDataOption struct {
		// Name of the parameter
		Name CommandOptionName `json:"name"`
		// Value of application command option type
		Type int64 `json:"type"`
		// Value of the option resulting from user input.
		// string, integer, double, or boolean
		Value ApplicationCommandDataOptionValue `json:"value"`
	}

	ApplicationCommandDataOptionValue struct {
		ValueString  *string
		ValueInt64   *int64
		ValueFloat64 *float64
		ValueBool    *bool
	}

	ApplicationCommandName  string
	InteractionType         int8
	InteractionContextType  int8
	InteractionCallbackType int8
	InteractionID           string
	InteractionToken        string

	// https://docs.discord.com/developers/interactions/receiving-and-responding#interaction-response-object
	InteractionResponseObject struct {
		Type InteractionCallbackType  `json:"type"`
		Data *InteractionCallbackData `json:"data"`
	}

	InteractionCreateCallbackResponse struct {
		ID    InteractionID
		Token InteractionToken
		// The payload that is sent to the request body when creating an interaction callback response.
		Body InteractionResponseObject
	}

	// https://docs.discord.com/developers/interactions/receiving-and-responding#interaction-response-object-interaction-callback-data-structure
	InteractionCallbackData struct {
		// Whether the response is TTS
		TTS *bool `json:"tts"`
		// Message content
		Content *string `json:"content"`
		// Message flags combined as a bitfield (only SUPPRESS_EMBEDS, EPHEMERAL, IS_COMPONENTS_V2, IS_VOICE_MESSAGE, and SUPPRESS_NOTIFICATIONS can be set)
		// https://docs.discord.com/developers/resources/message#message-object-message-flags
		Flags *int8 `json:"flags"`
	}

	MessageFlags int16
)

const (
	InteractionTypePing                           InteractionType = 1
	InteractionTypeApplicationCommand             InteractionType = 2
	InteractionTypeMessageComponent               InteractionType = 3
	InteractionTypeApplicationCommandAutocomplete InteractionType = 4
	InteractionTypeModalSubmit                    InteractionType = 5

	// Interaction can be used within servers
	InteractionContextTypeGuild InteractionContextType = 0
	// Interaction can be used within DMs with the app’s bot user
	InteractionContextTypeBotDM InteractionContextType = 1
	// Interaction can be used within Group DMs and DMs other than the app’s bot user
	InteractionContextTypePrivateChannel InteractionContextType = 2

	// Interaction Callback Response
	// https://docs.discord.com/developers/interactions/receiving-and-responding#interaction-response-object-interaction-callback-type

	// ACK a Ping
	InteractionCallbackTypePONG InteractionCallbackType = 1
	// Respond to an interaction with a message
	InteractionCallbackTypeCHANNEL_MESSAGE_WITH_SOURCE InteractionCallbackType = 4
	// ACK an interaction and edit a response later, the user sees a loading state
	InteractionCallbackTypeDEFERRED_CHANNEL_MESSAGE_WITH_SOURCE InteractionCallbackType = 5
	// For components, ACK an interaction and edit the original message later; the user does not see a loading state
	InteractionCallbackTypeDEFERRED_UPDATE_MESSAGE InteractionCallbackType = 6
	// For components, edit the message the component was attached to
	InteractionCallbackTypeUPDATE_MESSAGE InteractionCallbackType = 7
	// Respond to an autocomplete interaction with suggested choices
	InteractionCallbackTypeAPPLICATION_COMMAND_AUTOCOMPLETE_RESULT InteractionCallbackType = 8
	// Respond to an interaction with a popup modal
	InteractionCallbackTypeMODAL InteractionCallbackType = 9
	// Launch the Activity associated with the app. Only available for apps with Activities enabled
	InteractionCallbackTypeLAUNCH_ACTIVITY InteractionCallbackType = 12

	// this message has been published to subscribed channels (via Channel Following)
	MessageFlagBitCROSSPOSTED = 1 << 0
	// this message originated from a message in another channel (via Channel Following)
	MessageFlagBitIS_CROSSPOST = 1 << 1
	// do not include any embeds when serializing this message
	MessageFlagBitSUPPRESS_EMBEDS = 1 << 2
	// the source message for this crosspost has been deleted (via Channel Following)
	MessageFlagBitSOURCE_MESSAGE_DELETED MessageFlags = 1 << 3
	// this message came from the urgent message system
	MessageFlagBitURGENT MessageFlags = 1 << 4
	// this message has an associated thread, with the same id as the message
	MessageFlagBitHAS_THREAD MessageFlags = 1 << 5
	// this message is only visible to the user who invoked the Interaction
	MessageFlagBitEPHEMERAL MessageFlags = 1 << 6
	// this message is an Interaction Response and the bot is “thinking”
	MessageFlagBitLOADING MessageFlags = 1 << 7
	// this message failed to mention some roles and add their members to the thread
	MessageFlagBitFAILED_TO_MENTION_SOME_ROLES_IN_THREAD MessageFlags = 1 << 8
	// this message will not trigger push and desktop notifications
	MessageFlagBitSUPPRESS_NOTIFICATIONS MessageFlags = 1 << 12
	// this message is a voice message
	MessageFlagBitIS_VOICE_MESSAGE MessageFlags = 1 << 13
	// this message has a snapshot (via Message Forwarding)
	MessageFlagBitHAS_SNAPSHOT MessageFlags = 1 << 14
	// allows you to create fully component-driven messages
	MessageFlagBitIS_COMPONENTS_V2 = 1 << 15
)

func (value *ApplicationCommandDataOptionValue) UnmarshalJSON(b []byte) error {
	// string check
	if b[0] == '"' && b[len(b)-1] == '"' && len(b) > 1 {
		str := string(b[1 : len(b)-1])
		value.ValueString = &str
		return nil
	}

	if slices.Contains(b, '.') {
		f64, err := strconv.ParseFloat(string(b), 64)
		if err == nil {
			value.ValueFloat64 = &f64
			return nil
		}
		return nil
	}

	i64, err := strconv.ParseInt(string(b), 10, 64)
	if err == nil {
		value.ValueInt64 = &i64
		return nil
	}

	// bool check
	if slices.Equal(b, []byte("true")) || slices.Equal(b, []byte("false")) {
		v := slices.Equal(b, []byte("true"))
		value.ValueBool = &v
		return nil
	}

	return nil
}

func (c *Client) CreateInteractionResponse(params InteractionCreateCallbackResponse) error {
	bodyRaw, err := json.Marshal(params.Body)
	if err != nil {
		return fmt.Errorf("%s %w", err.Error(), ErrDiscordCommandFailedToSerializeCreatePayload)
	}

	// /interactions/{interaction.id}/{interaction.token}/callback
	url := fmt.Sprintf("%s/interactions/%s/%s/callback", DISCORD_BASE_URL, params.ID, params.Token)
	buf := bytes.NewBuffer(bodyRaw)
	res, err := c.Do(url, buf)
	if err != nil {
		return fmt.Errorf("%s %w", err.Error(), ErrDiscordCommandFailedToCreate)
	}

	if !(res.StatusCode == http.StatusOK || res.StatusCode == http.StatusNoContent) {
		rawRes, err := io.ReadAll(res.Body)
		if err != nil {
			return fmt.Errorf("%s %w", err.Error(), ErrDiscordCommandFailedToCreate)
		}
		fmt.Printf("%s", string(rawRes))
	}

	return nil
}
