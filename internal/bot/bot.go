package bot

import (
	"UserInteractiveBot/internal/config"
	"UserInteractiveBot/internal/storage"
	"errors"
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Bot struct {
	api      *tgbotapi.BotAPI
	cfg      *config.Config
	storage  storage.Storage
	handlers *Handlers
}

func New(cfg *config.Config, storage storage.Storage) (*Bot, error) {
	token := os.Getenv(cfg.Telegram.TokenEnv)
	if token == "" {
		return nil, errors.New("telegram token not set")
	}

	api, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}

	api.Debug = cfg.Telegram.Debug
	log.Printf("Authorized on account %s", api.Self.UserName)

	return &Bot{
		api:      api,
		cfg:      cfg,
		storage:  storage,
		handlers: NewHandlers(api, cfg, storage),
	}, nil
}

func (b *Bot) Start() error {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := b.api.GetUpdatesChan(u)

	for update := range updates {
		b.handlers.HandleUpdate(update)
	}
	return nil
}
