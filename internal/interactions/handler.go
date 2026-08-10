package interactions

import (
	"crypto/ed25519"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"math/rand/v2"
	"net/http"
	"sync"

	sqlite "github.com/tacohirosystems/tacohiro/internal/database"
	"github.com/tacohirosystems/tacohiro/internal/discord"
)

type (
	Handler struct {
		DiscordBotConfig *discord.DiscordConfig
		// FIXME: Placeholder
		State            InMemoryCounter
		DiscordBotClient *discord.Client
		DB               map[string]*sqlite.DB
		Logger           *slog.Logger
	}

	InMemoryCounter struct {
		sync.Mutex
		SentLog     map[string]int64
		ReceivedLog map[string]int64
	}
)

const (
	SlashCommandGive discord.ApplicationCommandName = "give"
)

func (s *InMemoryCounter) GetSentLog() map[string]int64 {
	s.Lock()
	defer s.Unlock()
	return s.SentLog
}

func (s *InMemoryCounter) GetReceivedLog() map[string]int64 {
	s.Lock()
	defer s.Unlock()
	return s.ReceivedLog
}

func (h *Handler) Routes() {
	http.HandleFunc("POST /api/discord/interactions", h.ProcessInteractions)
}

func (h *Handler) ProcessInteractions(w http.ResponseWriter, r *http.Request) {
	sigHex := r.Header.Get("X-Signature-Ed25519")
	sig, err := hex.DecodeString(sigHex)
	if err != nil {
		h.Logger.Info("Invalid signature")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	ts := r.Header.Get("X-Signature-Timestamp")

	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.Logger.Info("Invalid body format")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	msg := fmt.Sprintf("%s%s", ts, body)
	if !ed25519.Verify(h.DiscordBotConfig.PublicKey, []byte(msg), sig) {
		h.Logger.Info("Invalid message")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	var rawPayload map[string]any
	err = json.Unmarshal(body, &rawPayload)
	if err != nil {
		h.Logger.Error(fmt.Sprintf("failed to deserialize %s\n", err.Error()))
	}

	switch discord.InteractionType(int8(rawPayload["type"].(float64))) {
	case discord.InteractionTypePing:
		w.Write(body)
	case discord.InteractionTypeApplicationCommand:
		var interactionPayload discord.Interaction[discord.ApplicationCommandData]
		err = json.Unmarshal(body, &interactionPayload)
		if err != nil {
			h.Logger.Error(fmt.Sprintf("failed to deserialize interaction: %s", err.Error()))
			return
		}

		h.Logger.Debug(fmt.Sprintf("Interaction payload: %+v\n", interactionPayload))
		h.Logger.Debug(fmt.Sprintf("Received a slash command: %s\n", interactionPayload.Data.Name))
		if interactionPayload.Data == nil {
			return
		}

		h.processSlashCommand(interactionPayload)
		h.Logger.Debug(fmt.Sprintf("Counters: %+v %+v", h.State.GetSentLog(), h.State.GetReceivedLog()))
		w.WriteHeader(http.StatusOK)
	default:
		w.WriteHeader(http.StatusOK)
		return
	}
}

func (h *Handler) processSlashCommand(interaction discord.Interaction[discord.ApplicationCommandData]) {
	senderID := interaction.Member.User.ID
	var recipientIDs []string
	var quantity int64
	switch interaction.Data.Name {
	case SlashCommandGive:
		// TODO: Check if recipient is a role, then we send to all members in that role.
		// TODO: Check if recipient is a bot. Maybe consider allowing bots to send and receive.
		for _, o := range interaction.Data.Options {
			switch o.Name {
			case CommandOptionNameRecipient:
				if o.Value.ValueString == nil {
					return
				}
				h.Logger.Debug(fmt.Sprintf("Recipient ID: %s\n", *o.Value.ValueString))
				recipientIDs = append(recipientIDs, *o.Value.ValueString)
			case CommandOptionNameRecipientExtra1:
				if o.Value.ValueString == nil {
					return
				}
				h.Logger.Debug(fmt.Sprintf("Recipient ID: %s\n", *o.Value.ValueString))
				recipientIDs = append(recipientIDs, *o.Value.ValueString)
			case CommandOptionNameRecipientExtra2:
				if o.Value.ValueString == nil {
					return
				}
				h.Logger.Debug(fmt.Sprintf("Recipient ID: %s\n", *o.Value.ValueString))
				recipientIDs = append(recipientIDs, *o.Value.ValueString)
			case "quantity":
				if o.Value.ValueInt64 != nil {
					h.Logger.Debug(fmt.Sprintf("Quantity of tacos: %d\n", *o.Value.ValueInt64))
					quantity = *o.Value.ValueInt64
				} else {
					h.Logger.Debug(fmt.Sprintf("Quantity of tacos: %+v", o.Value))
					return
				}
			default:
				h.Logger.Error(fmt.Sprintf("Unknown option %s\n", o.Name))
			}
		}
	default:
		return
	}

	var recipientsStr string
	for i := 0; i < len(recipientIDs); i++ {
		if i == 0 {
			recipientsStr = fmt.Sprintf("<@%s>", recipientIDs[i])
			continue
		}

		recipientsStr = fmt.Sprintf("%s <@%s>", recipientsStr, recipientIDs[i])
	}

	var msg string
	if quantity == 0 {
		msg = fmt.Sprintf("<@%s> gave %s %dx :taco:. _Hm... someone is stingy_", senderID, recipientsStr, quantity)
	}

	if quantity < 0 {
		msg = fmt.Sprintf("<@%s> gave %s a :taco: debt of %d.", senderID, recipientsStr, quantity*-1)
	}

	var loseOneTaco bool
	if rand.Float64() < 0.05 {
		loseOneTaco = true
	}

	if loseOneTaco && quantity >= 1 {
		quantity -= 1
		msg = fmt.Sprintf("<@%s> tried to give %s %dx :taco:! But Hiro ate one so it's one less...", senderID, recipientsStr, quantity+1)
	} else if quantity >= 1 {
		msg = fmt.Sprintf("<@%s> gave %s %dx :taco:!", senderID, recipientsStr, quantity)
	}

	h.State.Lock()
	h.State.SentLog[senderID] += quantity
	for _, recipientID := range recipientIDs {
		h.State.ReceivedLog[recipientID] += quantity
	}
	h.State.Unlock()

	err := h.DiscordBotClient.CreateInteractionResponse(discord.InteractionCreateCallbackResponse{
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
	})
	if err != nil {
		h.Logger.Error(fmt.Sprintf("%s", err.Error()))
	}
}
