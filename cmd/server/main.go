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

	// Registers the default TacoHiro commands.
	// TODO: Register commands
	// - Leaderboards?
	log.Println("Slash commands: Initializing...")
	commandsClient := discord.CommandsClient(client)
	err = commandsClient.CreateCommand(discord.CreateCommandParams{
		Name:        "give",
		Type:        discord.CommandTypeChatInput,
		Description: "Give tacos to others!",
		Options:     []discord.CommandOption{
			{
				Type:         discord.CommandOptionTypeMentionable,
				Name:         "users",
				Description:  "Who would you like to give tacos to?",
				Required:     true,
				Autocomplete: true,
			},
			{
				Type:        discord.CommandOptionTypeInteger,
				Name:        "quantity",
				Description: "How many?",
				Required:    true,
				Autocomplete: true,
			},
	},
	})
	if err != nil {
		panic(err)
	}
	log.Println("Slash commands: OK")

	interactionsHandler := interactions.Handler{Config: config}
	interactionsHandler.Routes()
	log.Fatal(http.ListenAndServe(":8080", nil))
}
