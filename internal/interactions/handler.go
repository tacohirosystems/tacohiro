package interactions

import (
	"crypto/ed25519"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/tacohirosystems/tacohiro/internal/discord"
)

type (
	Handler struct {
		DiscordBotConfig *discord.DiscordConfig

		Service *Service

		Logger *slog.Logger
	}
)

const (
	SlashCommandGive discord.ApplicationCommandName = "give"
)

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

		if err := h.Service.ProcessInteraction(interactionPayload); err != nil {
			h.Logger.Error(fmt.Sprintf("failed to process interaction: %s", err.Error()))
		}

		w.WriteHeader(http.StatusOK)
	default:
		w.WriteHeader(http.StatusOK)
		return
	}
}
