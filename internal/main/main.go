package main

import (
	"UserInteractiveBot/internal/bot"
	"UserInteractiveBot/internal/config"
	"UserInteractiveBot/internal/storage/json"
	"log"
)

func main() {
	cfg, err := config.Load("internal/config/config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	storage := json.NewJSONStorage(
		cfg.Storage.UsersFile,
		cfg.Storage.AdminsFile,
		cfg.Storage.NicknamesFile,
		cfg.Storage.PostsFile,
	)

	bot, err := bot.New(cfg, storage)
	if err != nil {
		log.Fatalf("Failed to create bot: %v", err)
	}

	if err := bot.Start(); err != nil {
		log.Fatalf("Bot failed to start: %v", err)
	}
}
