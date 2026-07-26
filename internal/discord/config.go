package discord

import (
	"crypto/ed25519"
	"encoding/hex"
	"fmt"
	"os"
)

type (
	DiscordConfig struct {
		ApplicationID string
		PublicKey ed25519.PublicKey
		Token BotToken
	}
)

const (
	DISCORD_APPLICATION_ID_KEY string = "DISCORD_APPLICATION_ID"
	DISCORD_PUBLIC_KEY_KEY string = "DISCORD_PUBLIC_KEY"
	DISCORD_BOT_TOKEN_KEY string = "DISCORD_BOT_TOKEN"
)

var (
	ErrDiscordRequiredConfig = fmt.Errorf("discord: a required field is missing or empty")
	ErrDiscordInvalidPublicKey = fmt.Errorf("discord: public key is invalid. must be in hex")
)

func GetConfigFromEnv() (*DiscordConfig, error) {
	applicationID := os.Getenv(DISCORD_APPLICATION_ID_KEY)
	if applicationID == "" {
		return nil, fmt.Errorf("%w %s", DISCORD_APPLICATION_ID_KEY)
	}

	publicKey := os.Getenv(DISCORD_PUBLIC_KEY_KEY)
	if publicKey == "" {
		return nil, fmt.Errorf("%w %s", DISCORD_PUBLIC_KEY_KEY)
	}

	publicKeyBs, err := hex.DecodeString(publicKey)
	if err != nil {
		return nil, ErrDiscordInvalidPublicKey
	}

	botToken := os.Getenv(DISCORD_BOT_TOKEN_KEY)
	if botToken == "" {
		return nil, fmt.Errorf("%w %s", DISCORD_BOT_TOKEN_KEY)
	}

	config := DiscordConfig{
		ApplicationID: applicationID,
		PublicKey:     ed25519.PublicKey(publicKeyBs),
		Token:         BotToken(botToken),
	}
	return &config, nil
}
