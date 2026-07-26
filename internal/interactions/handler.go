package interactions

import (
	"crypto/ed25519"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/tacohirosystems/tacohiro/internal/discord"
)

type (
	Handler struct {
		Config *discord.DiscordConfig
	}
)

func (h Handler) Routes() {
	http.HandleFunc("POST /api/discord/interactions", h.ProcessInteractions)
}

func (h Handler) ProcessInteractions(w http.ResponseWriter, r *http.Request) {
	sigHex := r.Header.Get("X-Signature-Ed25519")
	sig, err := hex.DecodeString(sigHex)
	if err != nil {
		log.Print("Invalid signature")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	ts := r.Header.Get("X-Signature-Timestamp")

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Print("invalid body format")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	msg := fmt.Sprintf("%s%s", ts, body)
	log.Printf("\n%s\n", msg)
	if !ed25519.Verify(h.Config.PublicKey, []byte(msg), sig) {
		log.Print("Invalid message")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	var rawPayload map[string]any
	err = json.Unmarshal(body, &rawPayload)
	if err != nil {
		log.Fatalf("failed to deserialize %s", err.Error())
	}

	switch discord.InteractionType(int8(rawPayload["type"].(float64))) {
	case discord.InteractionTypePing:
		w.Write(body)
	case discord.InteractionTypeApplicationCommand:
		var payload discord.Interaction[discord.ApplicationCommandData]
		err = json.Unmarshal(body, &payload)
		if err != nil {
			log.Fatalf("%s", err.Error())
			return
		}

		log.Printf("%+v\n", payload)
		log.Printf("Received a slash command: %s\n", payload.Data.Name)
		w.WriteHeader(http.StatusOK)
	default:
		w.WriteHeader(http.StatusOK)
	}
}
