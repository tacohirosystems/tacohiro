package main

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"

	sqlite "github.com/tacohirosystems/tacohiro/internal/database"
	"github.com/tacohirosystems/tacohiro/internal/discord"
	"github.com/tacohirosystems/tacohiro/internal/interactions"
)

func main() {
	userDBPaths, err := filepath.Glob("./user-*.db")
	userDBs := make(map[string]*sqlite.DB, len(userDBPaths))
	for _, userDBPath := range userDBPaths {
		userID := userDBPath[5 : len(userDBPath)-3]
		fmt.Println(userDBPath)
		fmt.Println(userID)

		userDB := &sqlite.DB{
			Path:      userDBPath,
			ReadPool:  make(chan struct{}, 10),
			WritePool: make(chan struct{}, 1),
		}

		userDB.SetMaxReadConnections(10)
		userDB.SetMaxWriteConnections(1)

		if err := userDB.SetPragmas(); err != nil {
			panic(err.Error())
		}
		userDBs[userID] = userDB
	}

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
		DB:               userDBs,
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
