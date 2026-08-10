package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"

	sqlite "github.com/tacohirosystems/tacohiro/internal/database"
	"github.com/tacohirosystems/tacohiro/internal/discord"
	"github.com/tacohirosystems/tacohiro/internal/interactions"
)

func main() {
	loggerOpts := slog.HandlerOptions{
		AddSource:   false,
		Level:       slog.LevelDebug,
		ReplaceAttr: nil,
	}
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &loggerOpts))
	if logger == nil {
		panic("Logger is nil")
	}
	logger.Debug("Logger: Initialized")

	// TODO: Make file path for user DBs configurable
	userDBPaths, err := filepath.Glob("./user-*.db")
	userDBs := make(map[string]*sqlite.DB, len(userDBPaths))
	for _, userDBPath := range userDBPaths {
		userID := userDBPath[5 : len(userDBPath)-3]
		logger.Debug("Initializing user DB", "step", "database", "path", userDBPath, "user_id", userID)

		userDB := &sqlite.DB{
			Path:      userDBPath,
			ReadPool:  make(chan struct{}, 10),
			WritePool: make(chan struct{}, 1),
			Logger:    logger,
		}

		userDB.SetMaxReadConnections(10)
		userDB.SetMaxWriteConnections(1)

		if err := userDB.SetPragmas(); err != nil {
			panic(err.Error())
		}
		logger.Debug("Initialized user DB", "step", "database", "path", userDBPath, "user_id", userID)
		userDBs[userID] = userDB
	}
	logger.Info("OK", "step", "database", "count", len(userDBPaths))

	// Initializes Discord configuration used for the SDK.
	config, err := discord.GetConfigFromEnv()
	if err != nil {
		logger.Error(err.Error(), "step", "discord")
		panic("Failed to initialize Discord")
	}

	client, err := discord.NewClient(http.DefaultClient, config)
	if err != nil {
		panic(fmt.Sprintf("%s", err.Error()))
	}

	logger.Debug("Initializing default slash commands...")
	// TODO: init commands
	err = interactions.InitCommands(client)
	if err != nil {
		panic(fmt.Sprintf("discord: Failed to initialize slash commands. %s", err.Error()))
	}
	logger.Info("OK", "step", "discord")

	logger.Debug("Registering routes...", "step", "server")
	interactionsHandler := interactions.Handler{
		DiscordBotConfig: config,
		DB:               userDBs,
		State: interactions.InMemoryCounter{
			SentLog:     make(map[string]int64),
			ReceivedLog: make(map[string]int64),
		},
		DiscordBotClient: client,
		Logger:           logger,
	}
	interactionsHandler.Routes()
	logger.Info("Registered routes", "step", "server")
	logger.Info("server: Running...", "step", "server")
	logger.Error(http.ListenAndServe(":8080", nil).Error())
}
