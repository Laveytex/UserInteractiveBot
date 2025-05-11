package main

import (
	"encoding/json"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"os"
)

const (
	AuthorizedUsersFile  = "authorized_users.json"
	AuthorizedAdminsFile = "authorized_admins.json"
	NicknamesFile        = "nicknames.json"
	SecretCode           = "your_secret_code"
	AdminCode            = "your_admin_code"
)

type UserNickname struct {
	ChatID   int64  `json:"chat_id"`
	Nickname string `json:"nickname"`
}

/***************************** РАБОТА С АВТОРИЗАЦИЕЙ *****************************/

// loadAuthorizedList загружает список авторизованных ID из JSON-файла
func loadAuthorizedList(filePath string) ([]int64, error) {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		emptyList := []int64{}
		data, err := json.Marshal(emptyList)
		if err != nil {
			return nil, err
		}
		if err := os.WriteFile(filePath, data, 0644); err != nil {
			return nil, err
		}
		return emptyList, nil
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var list []int64
	if err := json.Unmarshal(data, &list); err != nil {
		return nil, err
	}
	return list, nil
}

// saveAuthorizedList сохраняет список авторизованных ID в JSON-файл
func saveAuthorizedList(filePath string, list []int64) error {
	data, err := json.Marshal(list)
	if err != nil {
		return err
	}
	return os.WriteFile(filePath, data, 0644)
}

// isAuthorizedInList проверяет, есть ли chatID в указанном списке
func isAuthorizedInList(chatID int64, list []int64) bool {
	for _, id := range list {
		if id == chatID {
			return true
		}
	}
	return false
}

// isAuthorized проверяет, есть ли chatID в одном из списков
func isAuthorized(chatID int64, lists ...[]int64) bool {
	for _, list := range lists {
		for _, id := range list {
			if id == chatID {
				return true
			}
		}
	}
	return false
}

// removeFromList удаляет chatID из списка
func removeFromList(chatID int64, list []int64) []int64 {
	for i, id := range list {
		if id == chatID {
			return append(list[:i], list[i+1:]...)
		}
	}
	return list
}

/***************************** РАБОТА С НИКАМИ *****************************/

// saveNicknames сохраняет список никнеймов в JSON-файл
func saveNicknames(filePath string, list []UserNickname) error {
	data, err := json.Marshal(list)
	if err != nil {
		return err
	}
	return os.WriteFile(filePath, data, 0644)
}

// loadNicknames загружает список никнеймов из JSON-файла
func loadNicknames(filePath string) ([]UserNickname, error) {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		emptyList := []UserNickname{}
		data, err := json.Marshal(emptyList)
		if err != nil {
			return nil, err
		}
		if err := os.WriteFile(filePath, data, 0644); err != nil {
			return nil, err
		}
		return emptyList, nil
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var list []UserNickname
	if err := json.Unmarshal(data, &list); err != nil {
		return nil, err
	}
	return list, nil
}

// getNickname возвращает никнейм пользователя или пустую строку
func getNickname(chatID int64, nicknames []UserNickname) string {
	for _, user := range nicknames {
		if user.ChatID == chatID {
			return user.Nickname
		}
	}
	return ""
}

// setNickname устанавливает или обновляет никнейм пользователя
func setNickname(chatID int64, nickname string, nicknames []UserNickname) []UserNickname {
	for i, user := range nicknames {
		if user.ChatID == chatID {
			nicknames[i].Nickname = nickname
			return nicknames
		}
	}
	return append(nicknames, UserNickname{ChatID: chatID, Nickname: nickname})
}

func main() {
	// Получаем токен из переменной среды
	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	if token == "" {
		log.Fatal("TELEGRAM_BOT_TOKEN is not set in environment variables")
	}

	// Инициализация бота
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Fatalf("Failed to initialize bot: %v", err)
	}

	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	// Загружаем списки авторизованных пользователей, администраторов и никнеймов
	authorizedUsers, err := loadAuthorizedList(AuthorizedUsersFile)
	if err != nil {
		log.Fatalf("Failed to load authorized users: %v", err)
	}
	authorizedAdmins, err := loadAuthorizedList(AuthorizedAdminsFile)
	if err != nil {
		log.Fatalf("Failed to load authorized admins: %v", err)
	}
	nicknames, err := loadNicknames(NicknamesFile)
	if err != nil {
		log.Fatalf("Failed to load nicknames: %v", err)
	}

	// Настройка получения обновлений
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	// Обработка входящих сообщений
	for update := range updates {
		if update.Message == nil && update.CallbackQuery == nil {
			continue
		}

		chatID := int64(0)
		var msg tgbotapi.MessageConfig

		// Обработка callback-запросов от кнопок
		if update.CallbackQuery != nil {
			chatID = update.CallbackQuery.Message.Chat.ID
			msg = tgbotapi.NewMessage(chatID, "")
			data := update.CallbackQuery.Data

			switch data {
			case "create_nick":
				msg.Text = "Введите ваш новый никнейм:"
				bot.Send(tgbotapi.NewCallback(update.CallbackQuery.ID, ""))
			case "change_nick":
				msg.Text = "Введите ваш новый никнейм:"
				bot.Send(tgbotapi.NewCallback(update.CallbackQuery.ID, ""))
			case "my_nick":
				nick := getNickname(chatID, nicknames)
				if nick == "" {
					msg.Text = "У вас еще нет никнейма. Создайте его с помощью кнопки 'Создать ник'."
				} else {
					msg.Text = fmt.Sprintf("Ваш текущий никнейм: %s", nick)
				}
				bot.Send(tgbotapi.NewCallback(update.CallbackQuery.ID, ""))
			}
			if _, err := bot.Send(msg); err != nil {
				log.Printf("Failed to send message: %v", err)
			}
			continue
		}

		chatID = update.Message.Chat.ID
		msg = tgbotapi.NewMessage(chatID, "")

		// Обработка команды /start
		if update.Message.IsCommand() && update.Message.Command() == "start" {
			args := update.Message.CommandArguments()
			if args == SecretCode {
				if !isAuthorized(chatID, authorizedUsers, authorizedAdmins) {
					authorizedUsers = append(authorizedUsers, chatID)
					if err := saveAuthorizedList(AuthorizedUsersFile, authorizedUsers); err != nil {
						log.Printf("Failed to save authorized users: %v", err)
						msg.Text = "Ошибка при сохранении авторизации. Попробуйте позже."
					} else {
						msg.Text = "Вы подключились с секретным кодом! Теперь вы можете использовать бота как пользователь."
						// Добавляем кнопки для пользователей
						keyboard := tgbotapi.NewInlineKeyboardMarkup(
							tgbotapi.NewInlineKeyboardRow(
								tgbotapi.NewInlineKeyboardButtonData("Создать ник", "create_nick"),
								tgbotapi.NewInlineKeyboardButtonData("Поменять ник", "change_nick"),
								tgbotapi.NewInlineKeyboardButtonData("Мой ник", "my_nick"),
							),
						)
						msg.ReplyMarkup = keyboard
					}
				} else {
					msg.Text = "Вы уже авторизованы."
					keyboard := tgbotapi.NewInlineKeyboardMarkup(
						tgbotapi.NewInlineKeyboardRow(
							tgbotapi.NewInlineKeyboardButtonData("Создать ник", "create_nick"),
							tgbotapi.NewInlineKeyboardButtonData("Поменять ник", "change_nick"),
							tgbotapi.NewInlineKeyboardButtonData("Мой ник", "my_nick"),
						),
					)
					msg.ReplyMarkup = keyboard
				}
			} else if args == AdminCode {
				if !isAuthorizedInList(chatID, authorizedAdmins) {
					if isAuthorizedInList(chatID, authorizedUsers) {
						authorizedUsers = removeFromList(chatID, authorizedUsers)
						if err := saveAuthorizedList(AuthorizedUsersFile, authorizedUsers); err != nil {
							log.Printf("Failed to save authorized users: %v", err)
							msg.Text = "Ошибка при обновлении списка пользователей. Попробуйте позже."
						} else {
							authorizedAdmins = append(authorizedAdmins, chatID)
							if err := saveAuthorizedList(AuthorizedAdminsFile, authorizedAdmins); err != nil {
								log.Printf("Failed to save authorized admins: %v", err)
								msg.Text = "Ошибка при сохранении авторизации администратора. Попробуйте позже."
							} else {
								msg.Text = "Вы повышены до администратора! Теперь вы можете использовать бота как администратор."
							}
						}
					} else {
						authorizedAdmins = append(authorizedAdmins, chatID)
						if err := saveAuthorizedList(AuthorizedAdminsFile, authorizedAdmins); err != nil {
							log.Printf("Failed to save authorized admins: %v", err)
							msg.Text = "Ошибка при сохранении авторизации администратора. Попробуйте позже."
						} else {
							msg.Text = "Вы подключились с кодом администратора! Теперь вы можете использовать бота как администратор."
						}
					}
				}
			}
			if _, err := bot.Send(msg); err != nil {
				log.Printf("Failed to send message: %v", err)
			}
			continue
		}

		// Проверяем, авторизован ли пользователь
		if !isAuthorized(chatID, authorizedUsers, authorizedAdmins) {
			continue
		}

		// Обработка сообщений для установки никнейма
		if isAuthorizedInList(chatID, authorizedUsers) {
			// Проверяем, является ли сообщение ответом на запрос никнейма
			if update.Message.Text != "" {
				nickname := update.Message.Text
				nicknames = setNickname(chatID, nickname, nicknames)
				if err := saveNicknames(NicknamesFile, nicknames); err != nil {
					log.Printf("Failed to save nicknames: %v", err)
					msg.Text = "Ошибка при сохранении никнейма. Попробуйте позже."
				} else {
					msg.Text = fmt.Sprintf("Ваш никнейм установлен: %s", nickname)
					keyboard := tgbotapi.NewInlineKeyboardMarkup(
						tgbotapi.NewInlineKeyboardRow(
							tgbotapi.NewInlineKeyboardButtonData("Создать ник", "create_nick"),
							tgbotapi.NewInlineKeyboardButtonData("Поменять ник", "change_nick"),
							tgbotapi.NewInlineKeyboardButtonData("Мой ник", "my_nick"),
						),
					)
					msg.ReplyMarkup = keyboard
				}
			}
		} else if isAuthorizedInList(chatID, authorizedAdmins) {
			// Для администраторов отправляем их chatID
			msg.Text = fmt.Sprintf("Ваш ID: %d", chatID)
		}

		if _, err := bot.Send(msg); err != nil {
			log.Printf("Failed to send message: %v", err)
		}
	}
}
