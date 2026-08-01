package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/tacohirosystems/tacohiro/internal/discord"
	"github.com/tacohirosystems/tacohiro/internal/interactions"
)

const (
)

func main() {
	log.Println("Discord: Initializing...")
	// Initializes Discord configuration used for the SDK.
	config, err := discord.GetConfigFromEnv()
	if err != nil {
		panic(fmt.Sprintf("%s", err.Error()))
	}

	client, err := discord.NewClient(http.DefaultClient, config)
	if err != nil {
		panic(fmt.Sprintf("%s", err.Error()))
	}
	log.Println("Discord: OK")

	log.Println("Discord: Initializing default slash commands...")
	// TODO: init commands
	err = interactions.InitCommands(client)
	if err != nil {
		panic(fmt.Sprintf("discord: Failed to initialize slash commands. %s", err.Error()))
	}

	log.Println("server: Registering routes...")
	interactionsHandler := interactions.Handler{
		DiscordBotConfig: config,
		State: interactions.InMemoryCounter{
			SentLog:     make(map[string]int64),
			ReceivedLog: make(map[string]int64),
		},
		DiscordBotClient: client,
	}
	interactionsHandler.Routes()
	log.Println("server: Registered routes")
	log.Println("server: Running...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
