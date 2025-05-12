package bot

import (
	"UserInteractiveBot/internal/config"
	"UserInteractiveBot/internal/models"
	"UserInteractiveBot/internal/storage"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Handlers struct {
	api     *tgbotapi.BotAPI
	cfg     *config.Config
	storage storage.Storage
	state   *State
}

func NewHandlers(api *tgbotapi.BotAPI, cfg *config.Config, storage storage.Storage) *Handlers {
	return &Handlers{
		api:     api,
		cfg:     cfg,
		storage: storage,
		state:   NewState(),
	}
}

func (h *Handlers) HandleUpdate(update tgbotapi.Update) {
	if update.Message == nil && update.CallbackQuery == nil {
		return
	}

	if update.Message != nil && update.Message.IsCommand() && update.Message.Command() == "start" {
		h.handleStartCommand(update.Message)
	}

	// Добавим обработку callback-запросов
	if update.CallbackQuery != nil {
		h.handleCallbackQuery(update.CallbackQuery)
	}

	// Добавим обработку сообщений в зависимости от состояния
	if update.Message != nil && !update.Message.IsCommand() {
		h.handleMessage(update.Message)
	}
}

func (h *Handlers) handleStartCommand(msg *tgbotapi.Message) {
	chatID := msg.Chat.ID
	response := tgbotapi.NewMessage(chatID, "")

	args := msg.CommandArguments()
	if args == h.cfg.Auth.SecretCode {
		users, err := h.storage.Users().Load()
		if err != nil {
			log.Printf("Failed to load users: %v", err)
			return
		}
		users = append(users, chatID)
		if err := h.storage.Users().Save(users); err != nil {
			log.Printf("Failed to save users: %v", err)
			return
		}
		response.Text = "Вы авторизованы как пользователь!"
		response.ReplyMarkup = UserKeyboard()
	} else if args == h.cfg.Auth.AdminCode {
		admins, err := h.storage.Admins().Load()
		if err != nil {
			log.Printf("Failed to load admins: %v", err)
			return
		}
		admins = append(admins, chatID)
		if err := h.storage.Admins().Save(admins); err != nil {
			log.Printf("Failed to save admins: %v", err)
			return
		}
		response.Text = "Вы авторизованы как админ!"
		response.ReplyMarkup = AdminKeyboard()
	} else {
		response.Text = "Неверный код авторизации."
	}

	if _, err := h.api.Send(response); err != nil {
		log.Printf("Failed to send message: %v", err)
	}
}

func (h *Handlers) handleCallbackQuery(query *tgbotapi.CallbackQuery) {
	chatID := query.Message.Chat.ID
	data := query.Data

	switch data {
	case "create_nick", "change_nick":
		h.state.SetAdminState(chatID, "setting_nick")
		msg := tgbotapi.NewMessage(chatID, "Введите ваш новый ник:")
		h.api.Send(msg)
	case "my_nick":
		nicks, err := h.storage.Nicknames().Load()
		if err != nil {
			log.Printf("Failed to load nicknames: %v", err)
			return
		}
		for _, nick := range nicks {
			if nick.ChatID == chatID {
				msg := tgbotapi.NewMessage(chatID, "Ваш ник: "+nick.Nickname)
				h.api.Send(msg)
				return
			}
		}
		msg := tgbotapi.NewMessage(chatID, "У вас нет ника.")
		h.api.Send(msg)
	case "add_text":
		h.state.SetAdminState(chatID, "adding_text")
		msg := tgbotapi.NewMessage(chatID, "Введите текст поста:")
		h.api.Send(msg)
	case "add_photo":
		h.state.SetAdminState(chatID, "adding_photo")
		msg := tgbotapi.NewMessage(chatID, "Отправьте фото:")
		h.api.Send(msg)
	case "preview_post":
		post := h.state.GetCreatingPost(chatID)
		if post == nil {
			msg := tgbotapi.NewMessage(chatID, "Пост не создан.")
			h.api.Send(msg)
			return
		}
		msg := tgbotapi.NewMessage(chatID, "Превью поста:\n"+post.Text)
		h.api.Send(msg)
	}
}

func (h *Handlers) handleMessage(msg *tgbotapi.Message) {
	chatID := msg.Chat.ID
	state := h.state.GetAdminState(chatID)

	switch state {
	case "setting_nick":
		nick := msg.Text
		nicks, err := h.storage.Nicknames().Load()
		if err != nil {
			log.Printf("Failed to load nicknames: %v", err)
			return
		}
		for i, n := range nicks {
			if n.ChatID == chatID {
				nicks[i].Nickname = nick
			} else {
				nicks = append(nicks, models.UserNickname{ChatID: chatID, Nickname: nick})
			}
		}
		if err := h.storage.Nicknames().Save(nicks); err != nil {
			log.Printf("Failed to save nicknames: %v", err)
			return
		}
		h.state.ClearState(chatID)
		response := tgbotapi.NewMessage(chatID, "Ник установлен: "+nick)
		response.ReplyMarkup = UserKeyboard()
		h.api.Send(response)
	case "adding_text":
		post := h.state.GetCreatingPost(chatID)
		if post == nil {
			post = &models.CreatingPost{AdminChatID: chatID}
		}
		post.Text = msg.Text
		h.state.SetCreatingPost(chatID, post)
		h.state.SetAdminState(chatID, "idle")
		response := tgbotapi.NewMessage(chatID, "Текст добавлен.")
		response.ReplyMarkup = AdminKeyboard()
		h.api.Send(response)
	case "adding_photo":
		if msg.Photo == nil {
			response := tgbotapi.NewMessage(chatID, "Пожалуйста, отправьте фото.")
			h.api.Send(response)
			return
		}
		photo := msg.Photo[len(msg.Photo)-1].FileID
		post := h.state.GetCreatingPost(chatID)
		if post == nil {
			post = &models.CreatingPost{AdminChatID: chatID}
		}
		post.PhotoIDs = append(post.PhotoIDs, photo)
		h.state.SetCreatingPost(chatID, post)
		h.state.SetAdminState(chatID, "idle")
		response := tgbotapi.NewMessage(chatID, "Фото добавлено.")
		response.ReplyMarkup = AdminKeyboard()
		h.api.Send(response)
	}
}
